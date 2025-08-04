package usecase

import (
	"context"
	"fmt"

	"github.com/rfanazhari/distributed-queue-processor/domain/entity"
	"github.com/rfanazhari/distributed-queue-processor/domain/repository"
)

// reportWorkerUseCase implements the ReportWorkerUseCase interface
type reportWorkerUseCase struct {
	queueRepo     repository.QueueRepository
	reportUseCase ReportUseCase
	ctx           context.Context
	cancelFunc    context.CancelFunc
}

// NewReportWorkerUseCase creates a new report worker use case
func NewReportWorkerUseCase(
	queueRepo repository.QueueRepository,
	reportUseCase ReportUseCase,
) ReportWorkerUseCase {
	return &reportWorkerUseCase{
		queueRepo:     queueRepo,
		reportUseCase: reportUseCase,
	}
}

// StartWorker starts the worker that consumes report jobs
func (uc *reportWorkerUseCase) StartWorker(ctx context.Context) error {
	// Create a new context with cancel function
	uc.ctx, uc.cancelFunc = context.WithCancel(ctx)

	// Start consuming jobs
	err := uc.queueRepo.ConsumeReportJobs(uc.ctx, func(job entity.ReportJob) error {
		// Process the job by calling the report use case
		return uc.processJob(uc.ctx, job)
	})

	if err != nil {
		return fmt.Errorf("failed to start worker: %w", err)
	}

	fmt.Println("Report worker started successfully")
	return nil
}

// StopWorker stops the worker
func (uc *reportWorkerUseCase) StopWorker() error {
	if uc.cancelFunc != nil {
		uc.cancelFunc()
	}

	fmt.Println("Report worker stopped")
	return nil
}

// processJob processes a report job
func (uc *reportWorkerUseCase) processJob(ctx context.Context, job entity.ReportJob) error {
	fmt.Printf("Processing report job for survey ID: %s\n", job.SurveyID)

	// Call the report use case to generate the report
	err := uc.reportUseCase.GenerateReport(ctx, job.SurveyID)
	if err != nil {
		return fmt.Errorf("failed to generate report: %w", err)
	}

	return nil
}
