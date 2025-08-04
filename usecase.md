# Survey Report Generation System

## Project Description

The Survey Report Generation System is a robust application built using Golang, Redis, and RabbitMQ that follows clean architecture principles. This system efficiently handles survey response submissions and generates reports asynchronously through a message queue architecture.

The application is structured according to clean architecture with distinct layers:
- Domain Layer: Contains business entities, repository interfaces, and use case interfaces
- Infrastructure Layer: Implements repository interfaces using Redis and RabbitMQ
- Delivery Layer: Handles HTTP requests and responses
- Use Case Layer: Implements business logic

Key features include:
- Survey response submission via HTTP API
- Debouncing of report generation jobs using Redis locks
- Asynchronous report generation using RabbitMQ
- Graceful shutdown handling
- Docker support for running multiple instances

## Issue Identification

While the system is well-designed, there are several potential issues that could arise:

1. **Concurrency and Race Conditions**: When multiple instances are running, there's a risk of race conditions when processing the same survey ID simultaneously.

2. **Lock Expiration Timing**: The Redis lock has a TTL of 30 seconds. If report generation takes longer than this time, another instance might acquire the lock and trigger duplicate report generation.

3. **Message Queue Reliability**: If RabbitMQ experiences downtime or connectivity issues, report generation jobs could be lost.

4. **Scaling Challenges**: As the system scales, the single Redis and RabbitMQ instances could become bottlenecks.

5. **Error Handling and Recovery**: The system needs robust error handling for failed report generation attempts.

## Solutions

To address these issues, the following solutions have been implemented:

1. **Redis Lock Mechanism**: The application uses Redis SETNX command to implement a distributed lock mechanism that prevents multiple instances from processing the same survey simultaneously. When a survey response is submitted:
   - The system tries to acquire a lock with key `report:lock:{survey_id}`
   - If the lock is acquired, a report generation job is published to RabbitMQ
   - If the lock already exists, the job is skipped (debounced)
   - The lock has a TTL of 30 seconds and expires automatically

2. **Asynchronous Processing with RabbitMQ**: The application uses RabbitMQ to handle asynchronous report generation:
   - The publisher sends report generation jobs to the `generate_report_queue`
   - The consumer processes jobs from the queue and generates reports
   - Messages are acknowledged only after successful processing

3. **Docker Containerization**: The application supports Docker deployment with multiple instances, allowing for horizontal scaling and improved reliability.

4. **Graceful Shutdown**: The system implements graceful shutdown handling to ensure that in-progress tasks are completed before the application terminates.

5. **Comprehensive Testing**: Test scripts are provided to verify the system's behavior in various scenarios, including:
   - Testing a single instance
   - Testing multiple instances
   - Testing Redis lock mechanism
   - Testing debounce functionality

These solutions ensure that the Survey Report Generation System is robust, scalable, and reliable, even when handling high volumes of survey submissions across multiple application instances.