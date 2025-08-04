# PowerShell script to test the survey response submission API

# Wait for the server to start
Write-Host "Waiting for server to start..."
Start-Sleep -Seconds 5

# Submit a survey response
Write-Host "Submitting survey response..."
$response1 = Invoke-RestMethod -Uri "http://localhost:8080/api/survey/submit" `
    -Method Post `
    -ContentType "application/json" `
    -Body '{
        "survey_id": "survey123",
        "answers": {
            "question1": "answer1",
            "question2": "answer2",
            "question3": "answer3"
        }
    }'

Write-Host "Response: $($response1 | ConvertTo-Json)"

Write-Host "`nSubmitting another response for the same survey (should be debounced)..."
$response2 = Invoke-RestMethod -Uri "http://localhost:8080/api/survey/submit" `
    -Method Post `
    -ContentType "application/json" `
    -Body '{
        "survey_id": "survey123",
        "answers": {
            "question1": "different answer",
            "question2": "different answer",
            "question3": "different answer"
        }
    }'

Write-Host "Response: $($response2 | ConvertTo-Json)"

Write-Host "`nSubmitting response for a different survey..."
$response3 = Invoke-RestMethod -Uri "http://localhost:8080/api/survey/submit" `
    -Method Post `
    -ContentType "application/json" `
    -Body '{
        "survey_id": "survey456",
        "answers": {
            "question1": "answer1",
            "question2": "answer2",
            "question3": "answer3"
        }
    }'

Write-Host "Response: $($response3 | ConvertTo-Json)"

Write-Host "`nTest completed."