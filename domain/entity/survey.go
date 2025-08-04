package entity

// Survey represents a survey entity
type Survey struct {
	ID       string                 `json:"id"`
	Title    string                 `json:"title"`
	Status   string                 `json:"status"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// SurveyResponse represents a response to a survey
type SurveyResponse struct {
	ID        string                 `json:"id"`
	SurveyID  string                 `json:"survey_id"`
	Answers   map[string]interface{} `json:"answers"`
	CreatedAt int64                  `json:"created_at"`
}

// ReportJob represents a job to generate a report
type ReportJob struct {
	SurveyID string `json:"survey_id"`
}
