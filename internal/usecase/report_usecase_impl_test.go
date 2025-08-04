package usecase_test

import (
	"context"
	"errors"
	"github.com/rfanazhari/distributed-queue-processor/internal/usecase"
	"strings"
	"testing"
	"time"

	"github.com/rfanazhari/distributed-queue-processor/domain/entity"
)

// MockLockRepository is a manual mock for the LockRepository interface
type MockLockRepository struct {
	setLockFunc     func(ctx context.Context, key string, ttl time.Duration) (bool, error)
	releaseLockFunc func(ctx context.Context, key string) (bool, error)
}

func (m *MockLockRepository) SetLock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return m.setLockFunc(ctx, key, ttl)
}

func (m *MockLockRepository) ReleaseLock(ctx context.Context, key string) (bool, error) {
	return m.releaseLockFunc(ctx, key)
}

// MockQueueRepository is a manual mock for the QueueRepository interface
type MockQueueRepository struct {
	publishReportJobFunc  func(ctx context.Context, job entity.ReportJob) error
	consumeReportJobsFunc func(ctx context.Context, callback func(entity.ReportJob) error) error
	closeFunc             func() error
}

func (m *MockQueueRepository) PublishReportJob(ctx context.Context, job entity.ReportJob) error {
	return m.publishReportJobFunc(ctx, job)
}

func (m *MockQueueRepository) ConsumeReportJobs(ctx context.Context, callback func(entity.ReportJob) error) error {
	return m.consumeReportJobsFunc(ctx, callback)
}

func (m *MockQueueRepository) Close() error {
	return m.closeFunc()
}

func TestSubmitResponse_Success(t *testing.T) {
	// Setup mocks
	mockLockRepo := &MockLockRepository{
		setLockFunc: func(ctx context.Context, key string, ttl time.Duration) (bool, error) {
			// Verify the key and ttl
			if key != "report:lock:survey-123" || ttl != usecase.LockTTL {
				t.Errorf("Expected key=%s, ttl=%v, got key=%s, ttl=%v", "report:lock:survey-123", usecase.LockTTL, key, ttl)
			}
			return true, nil
		},
		releaseLockFunc: func(ctx context.Context, key string) (bool, error) {
			t.Errorf("ReleaseLock should not be called")
			return false, nil
		},
	}

	mockQueueRepo := &MockQueueRepository{
		publishReportJobFunc: func(ctx context.Context, job entity.ReportJob) error {
			// Verify the job
			if job.SurveyID != "survey-123" {
				t.Errorf("Expected SurveyID=%s, got %s", "survey-123", job.SurveyID)
			}
			return nil
		},
		consumeReportJobsFunc: func(ctx context.Context, callback func(entity.ReportJob) error) error {
			t.Errorf("ConsumeReportJobs should not be called")
			return nil
		},
		closeFunc: func() error {
			t.Errorf("Close should not be called")
			return nil
		},
	}

	// Create test data
	ctx := context.Background()
	response := entity.SurveyResponse{
		ID:        "resp-123",
		SurveyID:  "survey-123",
		Answers:   map[string]interface{}{"q1": "answer1"},
		CreatedAt: time.Now().Unix(),
	}

	// Create use case and call method
	uc := usecase.NewReportUseCase(mockLockRepo, mockQueueRepo)
	err := uc.SubmitResponse(ctx, response)

	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSubmitResponse_LockAlreadyAcquired(t *testing.T) {
	// Setup mocks
	mockLockRepo := &MockLockRepository{
		setLockFunc: func(ctx context.Context, key string, ttl time.Duration) (bool, error) {
			// Verify the key and ttl
			if key != "report:lock:survey-123" || ttl != usecase.LockTTL {
				t.Errorf("Expected key=%s, ttl=%v, got key=%s, ttl=%v", "report:lock:survey-123", usecase.LockTTL, key, ttl)
			}
			// Return false to simulate lock already acquired
			return false, nil
		},
		releaseLockFunc: func(ctx context.Context, key string) (bool, error) {
			t.Errorf("ReleaseLock should not be called")
			return false, nil
		},
	}

	mockQueueRepo := &MockQueueRepository{
		publishReportJobFunc: func(ctx context.Context, job entity.ReportJob) error {
			t.Errorf("PublishReportJob should not be called when lock is already acquired")
			return nil
		},
		consumeReportJobsFunc: func(ctx context.Context, callback func(entity.ReportJob) error) error {
			t.Errorf("ConsumeReportJobs should not be called")
			return nil
		},
		closeFunc: func() error {
			t.Errorf("Close should not be called")
			return nil
		},
	}

	// Create test data
	ctx := context.Background()
	response := entity.SurveyResponse{
		ID:        "resp-123",
		SurveyID:  "survey-123",
		Answers:   map[string]interface{}{"q1": "answer1"},
		CreatedAt: time.Now().Unix(),
	}

	// Create use case and call method
	uc := usecase.NewReportUseCase(mockLockRepo, mockQueueRepo)
	err := uc.SubmitResponse(ctx, response)

	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSubmitResponse_SetLockError(t *testing.T) {
	// Setup mocks
	expectedErr := errors.New("redis connection error")
	mockLockRepo := &MockLockRepository{
		setLockFunc: func(ctx context.Context, key string, ttl time.Duration) (bool, error) {
			// Verify the key and ttl
			if key != "report:lock:survey-123" || ttl != usecase.LockTTL {
				t.Errorf("Expected key=%s, ttl=%v, got key=%s, ttl=%v", "report:lock:survey-123", usecase.LockTTL, key, ttl)
			}
			// Return error to simulate lock error
			return false, expectedErr
		},
		releaseLockFunc: func(ctx context.Context, key string) (bool, error) {
			t.Errorf("ReleaseLock should not be called")
			return false, nil
		},
	}

	mockQueueRepo := &MockQueueRepository{
		publishReportJobFunc: func(ctx context.Context, job entity.ReportJob) error {
			t.Errorf("PublishReportJob should not be called when SetLock returns an error")
			return nil
		},
		consumeReportJobsFunc: func(ctx context.Context, callback func(entity.ReportJob) error) error {
			t.Errorf("ConsumeReportJobs should not be called")
			return nil
		},
		closeFunc: func() error {
			t.Errorf("Close should not be called")
			return nil
		},
	}

	// Create test data
	ctx := context.Background()
	response := entity.SurveyResponse{
		ID:        "resp-123",
		SurveyID:  "survey-123",
		Answers:   map[string]interface{}{"q1": "answer1"},
		CreatedAt: time.Now().Unix(),
	}

	// Create use case and call method
	uc := usecase.NewReportUseCase(mockLockRepo, mockQueueRepo)
	err := uc.SubmitResponse(ctx, response)

	// Assert results
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}

	if err != nil && !strings.Contains(err.Error(), "failed to set lock") {
		t.Errorf("Expected error to contain 'failed to set lock', got %v", err)
	}
}

func TestSubmitResponse_PublishJobError(t *testing.T) {
	// Setup mocks
	lockReleased := false
	mockLockRepo := &MockLockRepository{
		setLockFunc: func(ctx context.Context, key string, ttl time.Duration) (bool, error) {
			// Verify the key and ttl
			if key != "report:lock:survey-123" || ttl != usecase.LockTTL {
				t.Errorf("Expected key=%s, ttl=%v, got key=%s, ttl=%v", "report:lock:survey-123", usecase.LockTTL, key, ttl)
			}
			return true, nil
		},
		releaseLockFunc: func(ctx context.Context, key string) (bool, error) {
			if key != "report:lock:survey-123" {
				t.Errorf("Expected key=%s, got %s", "report:lock:survey-123", key)
			}
			lockReleased = true
			return true, nil
		},
	}

	expectedErr := errors.New("queue connection error")
	mockQueueRepo := &MockQueueRepository{
		publishReportJobFunc: func(ctx context.Context, job entity.ReportJob) error {
			// Verify the job
			if job.SurveyID != "survey-123" {
				t.Errorf("Expected SurveyID=%s, got %s", "survey-123", job.SurveyID)
			}
			// Return error to simulate publish error
			return expectedErr
		},
		consumeReportJobsFunc: func(ctx context.Context, callback func(entity.ReportJob) error) error {
			t.Errorf("ConsumeReportJobs should not be called")
			return nil
		},
		closeFunc: func() error {
			t.Errorf("Close should not be called")
			return nil
		},
	}

	// Create test data
	ctx := context.Background()
	response := entity.SurveyResponse{
		ID:        "resp-123",
		SurveyID:  "survey-123",
		Answers:   map[string]interface{}{"q1": "answer1"},
		CreatedAt: time.Now().Unix(),
	}

	// Create use case and call method
	uc := usecase.NewReportUseCase(mockLockRepo, mockQueueRepo)
	err := uc.SubmitResponse(ctx, response)

	// Assert results
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}

	if err != nil && !strings.Contains(err.Error(), "failed to publish report job") {
		t.Errorf("Expected error to contain 'failed to publish report job', got %v", err)
	}

	if !lockReleased {
		t.Errorf("Expected lock to be released when PublishReportJob fails")
	}
}

func TestGenerateReport(t *testing.T) {
	// Setup mocks
	mockLockRepo := &MockLockRepository{
		setLockFunc: func(ctx context.Context, key string, ttl time.Duration) (bool, error) {
			t.Errorf("SetLock should not be called")
			return false, nil
		},
		releaseLockFunc: func(ctx context.Context, key string) (bool, error) {
			t.Errorf("ReleaseLock should not be called")
			return false, nil
		},
	}

	mockQueueRepo := &MockQueueRepository{
		publishReportJobFunc: func(ctx context.Context, job entity.ReportJob) error {
			t.Errorf("PublishReportJob should not be called")
			return nil
		},
		consumeReportJobsFunc: func(ctx context.Context, callback func(entity.ReportJob) error) error {
			t.Errorf("ConsumeReportJobs should not be called")
			return nil
		},
		closeFunc: func() error {
			t.Errorf("Close should not be called")
			return nil
		},
	}

	// Create test data
	ctx := context.Background()
	surveyID := "survey-123"

	// Create use case and call method
	uc := usecase.NewReportUseCase(mockLockRepo, mockQueueRepo)
	err := uc.GenerateReport(ctx, surveyID)

	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
