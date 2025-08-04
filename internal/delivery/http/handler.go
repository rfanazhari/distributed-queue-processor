package http

import (
	"encoding/json"
	"github.com/rfanazhari/distributed-queue-processor/internal/usecase"
	"net/http"
	"time"

	"github.com/rfanazhari/distributed-queue-processor/domain/entity"
)

// Handler handles HTTP requests
type Handler struct {
	reportUseCase usecase.ReportUseCase
}

// NewHandler creates a new HTTP handler
func NewHandler(reportUseCase usecase.ReportUseCase) *Handler {
	return &Handler{
		reportUseCase: reportUseCase,
	}
}

// SubmitResponse handles the submission of a survey response
func (h *Handler) SubmitResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		SurveyID string                 `json:"survey_id"`
		Answers  map[string]interface{} `json:"answers"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.SurveyID == "" {
		http.Error(w, "Survey ID is required", http.StatusBadRequest)
		return
	}

	// Create a survey response
	response := entity.SurveyResponse{
		ID:        generateID(),
		SurveyID:  request.SurveyID,
		Answers:   request.Answers,
		CreatedAt: time.Now().Unix(),
	}

	// Submit the response
	if err := h.reportUseCase.SubmitResponse(r.Context(), response); err != nil {
		http.Error(w, "Failed to submit response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Response submitted successfully",
		"id":      response.ID,
	})
}

// SetupRoutes sets up the HTTP routes
func (h *Handler) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/api/survey/submit", h.SubmitResponse)

	return mux
}

// generateID generates a simple ID for the response
// In a real application, you would use a more robust ID generation method
func generateID() string {
	return time.Now().Format("20060102150405.000000")
}
