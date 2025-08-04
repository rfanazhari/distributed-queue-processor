package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/rfanazhari/distributed-queue-processor/domain/entity"
	"github.com/rfanazhari/distributed-queue-processor/domain/repository"
)

const (
	// LockTTL is the time-to-live for the Redis lock
	LockTTL = 30 * time.Second

	// LockKeyPrefix is the prefix for the Redis lock key
	LockKeyPrefix = "report:lock:"
)

// reportUseCase implements the ReportUseCase interface
type reportUseCase struct {
	lockRepo  repository.LockRepository
	queueRepo repository.QueueRepository
}

// NewReportUseCase creates a new report use case
func NewReportUseCase(
	lockRepo repository.LockRepository,
	queueRepo repository.QueueRepository,
) ReportUseCase {
	return &reportUseCase{
		lockRepo:  lockRepo,
		queueRepo: queueRepo,
	}
}

// SubmitResponse handles a new survey response submission
func (uc *reportUseCase) SubmitResponse(ctx context.Context, response entity.SurveyResponse) error {
	// Create lock key using survey ID
	lockKey := fmt.Sprintf("%s%s", LockKeyPrefix, response.SurveyID)

	// Try to acquire lock
	locked, err := uc.lockRepo.SetLock(ctx, lockKey, LockTTL)
	if err != nil {
		return fmt.Errorf("failed to set lock: %w", err)
	}

	// If lock was not acquired, it means a job is already scheduled
	if !locked {
		// Skip publishing a new job
		return nil
	}

	// Create and publish report job
	job := entity.ReportJob{
		SurveyID: response.SurveyID,
	}

	if err := uc.queueRepo.PublishReportJob(ctx, job); err != nil {
		// If publishing fails, release the lock
		uc.lockRepo.ReleaseLock(ctx, lockKey)
		return fmt.Errorf("failed to publish report job: %w", err)
	}

	return nil
}

// GenerateReport generates a report for the given survey ID
func (uc *reportUseCase) GenerateReport(ctx context.Context, surveyID string) error {
	// This is where the actual report generation logic would go
	// For this example, we'll just simulate some work

	fmt.Printf("Generating report for survey ID: %s\n", surveyID)

	// Simulate some work
	time.Sleep(2 * time.Second)

	fmt.Printf("Report generation completed for survey ID: %s\n", surveyID)

	return nil
}
