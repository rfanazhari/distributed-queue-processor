package repository

import (
	"context"
	"github.com/rfanazhari/distributed-queue-processor/domain/entity"
)

// QueueRepository defines the interface for message queue operations
type QueueRepository interface {
	// PublishReportJob publishes a report job to the queue
	PublishReportJob(ctx context.Context, job entity.ReportJob) error

	// ConsumeReportJobs starts consuming report jobs from the queue
	// The callback function is called for each job received
	ConsumeReportJobs(ctx context.Context, callback func(entity.ReportJob) error) error

	// Close closes the connection to the message queue
	Close() error
}
