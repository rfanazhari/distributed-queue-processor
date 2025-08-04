# PowerShell script to test multiple instances of the survey application

# Wait for services to be ready
Write-Host "Waiting for services to be ready..."
Start-Sleep -Seconds 10

# Array of ports for the different instances
$ports = @(8080, 8081, 8082)

# Array of survey IDs to test
$surveyIds = @("survey123", "survey456", "survey789")

# Function to submit a survey response
function Submit-SurveyResponse {
    param (
        [string]$Port,
        [string]$SurveyId,
        [string]$InstanceName
    )
    
    Write-Host "`n[$InstanceName] Submitting response for survey ID: $SurveyId"
    
    $body = @{
        survey_id = $SurveyId
        answers = @{
            question1 = "answer1 for $SurveyId on $InstanceName"
            question2 = "answer2 for $SurveyId on $InstanceName"
            question3 = "answer3 for $SurveyId on $InstanceName"
        }
    } | ConvertTo-Json
    
    try {
        $response = Invoke-RestMethod -Uri "http://localhost:$Port/api/survey/submit" `
            -Method Post `
            -ContentType "application/json" `
            -Body $body
        
        Write-Host "[$InstanceName] Response: $($response | ConvertTo-Json)"
    }
    catch {
        Write-Host "[$InstanceName] Error: $_" -ForegroundColor Red
    }
}

# Test 1: Submit different surveys to different instances
Write-Host "`n=== Test 1: Submit different surveys to different instances ==="
for ($i = 0; $i -lt $ports.Count; $i++) {
    Submit-SurveyResponse -Port $ports[$i] -SurveyId $surveyIds[$i] -InstanceName "Instance $($i+1)"
}

# Test 2: Submit the same survey to all instances (testing Redis lock)
Write-Host "`n=== Test 2: Submit the same survey to all instances (testing Redis lock) ==="
$sharedSurveyId = "shared-survey-001"
foreach ($port in $ports) {
    $instanceNum = $ports.IndexOf($port) + 1
    Submit-SurveyResponse -Port $port -SurveyId $sharedSurveyId -InstanceName "Instance $instanceNum"
    # Small delay between submissions
    Start-Sleep -Seconds 1
}

# Test 3: Submit to the same instance multiple times (testing debounce)
Write-Host "`n=== Test 3: Submit to the same instance multiple times (testing debounce) ==="
$port = $ports[0]
$debounceSurveyId = "debounce-survey-001"

for ($i = 1; $i -le 3; $i++) {
    Submit-SurveyResponse -Port $port -SurveyId $debounceSurveyId -InstanceName "Instance 1 (Submission $i)"
    # Small delay between submissions
    Start-Sleep -Milliseconds 500
}

# Test 4: Wait for lock to expire, then submit again
Write-Host "`n=== Test 4: Wait for lock to expire, then submit again ==="
Write-Host "Waiting 31 seconds for Redis lock to expire..."
Start-Sleep -Seconds 31

Submit-SurveyResponse -Port $ports[0] -SurveyId $debounceSurveyId -InstanceName "Instance 1 (After lock expiry)"

Write-Host "`nTest completed."