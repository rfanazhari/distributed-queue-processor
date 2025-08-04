package usecase

import (
	"context"
	"github.com/rfanazhari/distributed-queue-processor/domain/entity"
)

// ReportUseCase defines the interface for report generation use cases
type ReportUseCase interface {
	// SubmitResponse handles a new survey response submission
	// It checks for a lock and publishes a report job if needed
	SubmitResponse(ctx context.Context, response entity.SurveyResponse) error

	// GenerateReport generates a report for the given survey ID
	// This is the actual report generation logic that will be executed by the worker
	GenerateReport(ctx context.Context, surveyID string) error
}

// ReportWorkerUseCase defines the interface for the report worker
type ReportWorkerUseCase interface {
	// StartWorker starts the worker that consumes report jobs
	StartWorker(ctx context.Context) error

	// StopWorker stops the worker
	StopWorker() error
}
