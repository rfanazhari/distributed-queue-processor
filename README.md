# Distributed Queue Processor

A robust distributed processing system built with Go, Redis, and RabbitMQ that efficiently handles survey response submissions and generates reports asynchronously.

## Overview

This project demonstrates a distributed queue processing architecture that ensures reliable and efficient processing of survey responses across multiple application instances. It implements debouncing mechanisms using Redis locks and asynchronous processing with RabbitMQ to prevent duplicate work and ensure scalability.

## Features

- **Survey Response Submission API**: Simple HTTP API for submitting survey responses
- **Distributed Locking**: Uses Redis to implement distributed locks, preventing multiple instances from processing the same survey simultaneously
- **Asynchronous Processing**: Leverages RabbitMQ for reliable message queuing and asynchronous report generation
- **Debouncing Mechanism**: Prevents duplicate report generation for the same survey within a configurable time window
- **Horizontal Scaling**: Supports running multiple instances for improved throughput and reliability
- **Graceful Shutdown**: Ensures in-progress tasks are completed before application termination

## Architecture

The application follows clean architecture principles with distinct layers:

- **Domain Layer**: Contains business entities, repository interfaces, and use case interfaces
- **Infrastructure Layer**: Implements repository interfaces using Redis and RabbitMQ
- **Delivery Layer**: Handles HTTP requests and responses
- **Use Case Layer**: Implements business logic

## Project Structure

```
.
├── Dockerfile                 # Docker configuration for building the application
├── README.md                  # Project documentation
├── docker-compose.yml         # Docker Compose configuration for running multiple instances
├── docs/                      # Documentation files
│   └── claude.md              # Additional documentation
├── domain/                    # Domain layer
│   ├── entity/                # Business entities
│   │   └── survey.go          # Survey entity definitions
│   └── repository/            # Repository interfaces
│       ├── lock_repository.go # Interface for distributed locking
│       └── queue_repository.go # Interface for message queuing
├── go.mod                     # Go module definition
├── go.sum                     # Go module checksums
├── internal/                  # Internal application code
│   ├── delivery/              # Delivery layer
│   │   └── http/              # HTTP delivery implementation
│   │       └── handler.go     # HTTP request handlers
│   ├── infrastructure/        # Infrastructure layer
│   │   ├── rabbitmq/          # RabbitMQ implementation
│   │   │   └── queue_repository.go # RabbitMQ queue repository
│   │   └── redis/             # Redis implementation
│   │       └── lock_repository.go # Redis lock repository
│   └── usecase/               # Use case layer
│       ├── report_usecase.go          # Report use case interface
│       ├── report_usecase_impl.go     # Report use case implementation
│       └── report_worker_usecase_impl.go # Report worker use case implementation
├── main.go                    # Application entry point
├── test_api.ps1               # PowerShell script for testing the API
├── test_api.sh                # Shell script for testing the API
├── test_multiple_instances.ps1 # PowerShell script for testing multiple instances
└── usecase.md                 # Use case documentation
```

## Prerequisites

- Go 1.24 or higher
- Docker and Docker Compose
- Redis
- RabbitMQ

## Running the Application

### Using Docker Compose

The easiest way to run the application is using Docker Compose:

```bash
docker-compose up
```

This will start:
- Redis for distributed locking
- RabbitMQ for message queuing
- Three instances of the application on ports 8080, 8081, and 8082

### Manual Setup

1. Start Redis:
```bash
docker run -d -p 6379:6379 redis:alpine
```

2. Start RabbitMQ:
```bash
docker run -d -p 5672:5672 -p 15672:15672 rabbitmq:3-management-alpine
```

3. Build and run the application:
```bash
go build -o survey-app .
./survey-app
```

## Environment Variables

- `REDIS_ADDR`: Redis server address (default: "localhost:6379")
- `RABBITMQ_URL`: RabbitMQ connection URL (default: "amqp://guest:guest@localhost:5672/")
- `HTTP_ADDR`: HTTP server address (default: ":8080")

## API Endpoints

### Submit Survey Response

```
POST /api/survey/submit
```

Request body:
```json
{
  "survey_id": "survey123",
  "answers": {
    "question1": "answer1",
    "question2": "answer2",
    "question3": "answer3"
  }
}
```

Response:
```json
{
  "message": "Response submitted successfully",
  "id": "20250804211345.123456"
}
```

## Testing

The repository includes PowerShell scripts for testing the application:

- `test_api.ps1`: Tests the survey submission API on a single instance
- `test_multiple_instances.ps1`: Tests the distributed locking and debouncing mechanisms across multiple instances

To run the tests:

```powershell
.\test_api.ps1
.\test_multiple_instances.ps1
```

## Key Implementation Details

1. **Redis Lock Mechanism**: Uses Redis SETNX command to implement a distributed lock with key `report:lock:{survey_id}` and a TTL of 30 seconds.

2. **Asynchronous Processing**: Uses RabbitMQ to handle asynchronous report generation:
    - Publisher sends jobs to the `generate_report_queue`
    - Consumer processes jobs and generates reports
    - Messages are acknowledged only after successful processing

3. **Graceful Shutdown**: Implements proper shutdown handling to ensure in-progress tasks are completed before termination.

## License

This project is open source and available under the MIT license.