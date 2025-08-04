package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rfanazhari/distributed-queue-processor/domain/entity"
	"github.com/rfanazhari/distributed-queue-processor/domain/repository"
)

const (
	queueName = "generate_report_queue"
)

// QueueRepository implements the repository.QueueRepository interface using RabbitMQ
type QueueRepository struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewQueueRepository creates a new RabbitMQ queue repository
func NewQueueRepository(url string) (repository.QueueRepository, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Declare the queue
	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	return &QueueRepository{
		conn:    conn,
		channel: ch,
	}, nil
}

// PublishReportJob publishes a report job to the queue
func (r *QueueRepository) PublishReportJob(ctx context.Context, job entity.ReportJob) error {
	body, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	err = r.channel.PublishWithContext(
		ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Make message persistent
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	return nil
}

// ConsumeReportJobs starts consuming report jobs from the queue
func (r *QueueRepository) ConsumeReportJobs(ctx context.Context, callback func(entity.ReportJob) error) error {
	msgs, err := r.channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}

				var job entity.ReportJob
				if err := json.Unmarshal(msg.Body, &job); err != nil {
					// Log error and acknowledge message to remove it from queue
					fmt.Printf("Error unmarshaling job: %v\n", err)
					msg.Ack(false)
					continue
				}

				// Process the job
				if err := callback(job); err != nil {
					// Log error but don't acknowledge to retry later
					fmt.Printf("Error processing job: %v\n", err)
					msg.Nack(false, true)
					continue
				}

				// Acknowledge the message
				msg.Ack(false)
			}
		}
	}()

	return nil
}

// Close closes the connection to RabbitMQ
func (r *QueueRepository) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
	return nil
}
