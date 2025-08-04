#!/bin/bash

# This script tests the survey response submission API

# Wait for the server to start
echo "Waiting for server to start..."
sleep 5

# Submit a survey response
echo "Submitting survey response..."
curl -X POST http://localhost:8080/api/survey/submit \
  -H "Content-Type: application/json" \
  -d '{
    "survey_id": "survey123",
    "answers": {
      "question1": "answer1",
      "question2": "answer2",
      "question3": "answer3"
    }
  }'

echo -e "\n\nSubmitting another response for the same survey (should be debounced)..."
curl -X POST http://localhost:8080/api/survey/submit \
  -H "Content-Type: application/json" \
  -d '{
    "survey_id": "survey123",
    "answers": {
      "question1": "different answer",
      "question2": "different answer",
      "question3": "different answer"
    }
  }'

echo -e "\n\nSubmitting response for a different survey..."
curl -X POST http://localhost:8080/api/survey/submit \
  -H "Content-Type: application/json" \
  -d '{
    "survey_id": "survey456",
    "answers": {
      "question1": "answer1",
      "question2": "answer2",
      "question3": "answer3"
    }
  }'

echo -e "\n\nTest completed."