# Insider Message Sender

A Golang-based automatic message sending system that processes and sends messages from a database every 2 minutes using a custom scheduler implementation.

## 🚀 Features

- **Automatic Message Processing**: Sends 2 unsent messages every 2 minutes
- **Database Integration**: PostgreSQL for message storage with character limits
- **Redis Caching**: Caches messageId and sending time for tracking
- **RESTful API**: Start/stop scheduler and retrieve sent messages
- **Swagger Documentation**: Complete API documentation
- **Docker Support**: Full containerized deployment
- **Concurrent Processing**: Parallel message sending with goroutines
- **Retry Mechanism**: Automatic retry with exponential backoff for failed requests
- **Context-Aware Operations**: HTTP requests and cache operations respect cancellation
- **Graceful Shutdown**: Proper cleanup of connections on application exit
- **Production Ready**: Connection pooling, error handling, signal handling

## 📋 Requirements

- Golang
- Docker & Docker Compose

## 🛠️ Quick Start

### Option 1: Docker (Recommended)

```bash
# Clone the repository
git clone https://github.com/taidvt1008/insider-assessment.git
cd insider-assessment

# Setup environment
make setup

# ⚠️ IMPORTANT: Configure webhook URL for Docker
# Edit .env.docker and update WEBHOOK_URL to your actual webhook endpoint
# Example: WEBHOOK_URL=https://webhook.site/your-unique-url

# Start all services (PostgreSQL, Redis, App)
make up
```

### Option 2: Local Development

```bash
# Setup dependencies
make setup

# ⚠️ IMPORTANT: Configure webhook URL for local development
# Edit .env and update WEBHOOK_URL to your actual webhook endpoint
# Example: WEBHOOK_URL=https://webhook.site/your-unique-url

# Run the application
make run
```

## 🔧 Configuration

### Environment Files

After running `make setup`, you'll have two environment files:

- **`.env`** - For local development (`make run`)
- **`.env.docker`** - For Docker deployment (`make up`)

### Required Configuration

**⚠️ CRITICAL: Update WEBHOOK_URL in both files before running the application!**

#### For Local Development (`.env`):
```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=insider

# Redis Configuration
REDIS_ADDR=localhost:6379

# Webhook Configuration - UPDATE THIS!
WEBHOOK_URL=https://webhook.site/your-unique-url

# Application Configuration
SEND_INTERVAL=2m
SERVER_PORT=8080
```

#### For Docker Deployment (`.env.docker`):
```env
# Database Configuration
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=insider

# Redis Configuration
REDIS_ADDR=redis:6379

# Webhook Configuration - UPDATE THIS!
WEBHOOK_URL=https://webhook.site/your-unique-url

# Application Configuration
SEND_INTERVAL=2m
SERVER_PORT=8080
```

### Getting a Webhook URL

1. Visit [webhook.site](https://webhook.site)
2. Copy your unique URL
3. Update `WEBHOOK_URL` in both `.env` and `.env.docker` files

## 📊 Database Schema

```sql
-- Create enum type for message status
CREATE TYPE message_status AS ENUM ('pending', 'sent', 'failed');

CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    phone_number VARCHAR(20) NOT NULL CHECK (phone_number ~ '^\+[0-9]{10,15}$'),
    content TEXT NOT NULL CHECK (char_length(content) <= 160),
    status message_status DEFAULT 'pending',
    sent_at TIMESTAMPTZ
);

-- Create indexes for better performance
CREATE INDEX idx_messages_status ON messages(status);
CREATE INDEX idx_messages_sent_at ON messages(sent_at);
```

## 🎯 API Endpoints

### Health Check

#### Service Health Status
```bash
GET /health
```

**Response (Healthy):**
```json
{
  "status": "healthy",
  "timestamp": "2025-10-19T21:00:00Z",
  "services": {
    "database": "healthy",
    "scheduler": "running",
    "redis": "healthy"
  }
}
```

**Response (Unhealthy):**
```json
{
  "status": "unhealthy",
  "timestamp": "2025-10-19T21:00:00Z",
  "services": {
    "database": "unhealthy: connection refused",
    "scheduler": "stopped",
    "redis": "healthy"
  }
}
```

### Scheduler Control

#### Start Scheduler
```bash
POST /api/v1/scheduler/start
```

#### Stop Scheduler
```bash
POST /api/v1/scheduler/stop
```

### Message Management

#### Get Sent Messages (with pagination)
```bash
GET /api/v1/messages/sent?limit=10&offset=0
```

#### Get Failed Messages (with pagination)
```bash
GET /api/v1/messages/failed?limit=10&offset=0
```

### API Documentation
- **Swagger UI**: http://localhost:8080/swagger/index.html

## 🐳 Docker Commands

```bash
# Start all services
make up

# Stop services (keep data)
make down

# View logs
make logs

# Reset database and Redis
make reset-db

# Connect to PostgreSQL
make db

# Connect to Redis
make redis
```

## 🧪 Testing

```bash
# Test health check
curl http://localhost:8080/health

# Test scheduler start
make test-start

# Test scheduler stop
make test-stop

# Test message listing
make test-list
```

## 📁 Project Structure

```
insider-assessment/
├── cmd/server/         # Application entry point
├── internal/
│   ├── api/            # HTTP handlers and routing
│   ├── cache/          # Redis client implementation
│   ├── config/         # Configuration management
│   ├── constants/      # Application constants
│   ├── docs/           # Swagger documentation
│   ├── model/          # Data models and DTOs
│   ├── repository/     # Database access layer
│   └── scheduler/      # Background job scheduler
├── scripts/            # Database initialization
├── docker-compose.yml  # Multi-container setup
├── Dockerfile          # Application container
└── Makefile            # Development commands
```

## 🔄 How It Works

1. **Startup**: Application automatically starts the scheduler on deployment
2. **Processing**: Every 2 minutes, the scheduler:
   - Fetches 2 unsent messages from the database
   - Sends them concurrently to the webhook URL
   - Marks successful messages as "sent" in the database
   - Caches messageId and timestamp in Redis
3. **API Control**: Use REST endpoints to start/stop the scheduler
4. **Monitoring**: Retrieve sent messages with pagination support

## 📋 Constants

The application uses predefined constants for message status values:

```go
const (
    MessageStatusPending = "pending"
    MessageStatusSent    = "sent"
    MessageStatusFailed  = "failed"
)
```

This ensures type safety and prevents magic strings throughout the codebase.

## 🚦 Message Flow

```
Database (pending) → Scheduler → Webhook → Database (sent) + Redis (cache)
```

## 🛡️ Error Handling

- **Network Failures**: Automatic retry with exponential backoff (3 attempts)
- **Database Errors**: Graceful degradation with logging
- **Redis Failures**: Non-blocking cache operations
- **Invalid Messages**: Content length validation (160 chars max)
- **Connection Leaks**: Proper response body reading to prevent leaks
- **Context Cancellation**: All operations respect context cancellation

## 📈 Performance Features

- **Connection Pooling**: Optimized database and HTTP connections
- **Concurrent Processing**: Parallel message sending
- **Resource Management**: Proper cleanup and memory management
- **Retry Strategy**: Exponential backoff (1s → 2s → 4s) for failed requests
- **Context Management**: Efficient cancellation of long-running operations
- **Graceful Shutdown**: Signal handling for clean application termination

## 🔍 Monitoring & Logging

The application provides comprehensive logging for:
- Scheduler start/stop events
- Message processing status with retry attempts
- Database operations
- Webhook responses
- Error conditions and connection cleanup
- Graceful shutdown process

## 🤝 Development

### Available Commands

```bash
make setup       # Install dependencies and setup environment
make swag        # Generate Swagger documentation
make lint        # Run code quality checks
make build       # Build the application binary
make run         # Run locally with database services
make clean       # Clean build artifacts
```

### Testing Commands

```bash
make test-start     # Test scheduler start endpoint
make test-stop      # Test scheduler stop endpoint
make test-list      # Test get sent messages endpoint
make test-failed    # Test get failed messages endpoint
make test-graceful  # Test graceful shutdown functionality
```

### Code Quality

- **Linting**: golangci-lint integration
- **Documentation**: Swagger/OpenAPI 3.0
- **Testing**: Unit and integration test support
- **Formatting**: gofmt and goimports

## 🔍 Known Issues & Future Improvements

### Current Limitations

#### 1. **Logging System**
- **Issue**: Using standard `log` package instead of structured logging
- **Problem**: No log levels (DEBUG, INFO, WARN, ERROR), no structured fields
- **Impact**: Difficult to debug and monitor in production
- **Future**: Consider implementing structured logging with [logrus](https://github.com/sirupsen/logrus) or [zap](https://github.com/uber-go/zap)

#### 2. **Redis Caching Strategy**
- **Issue**: Redis only caches messageId + timestamp, not full message data
- **Problem**: API endpoints still query database for sent messages
- **Impact**: Potential performance bottleneck for high-volume message retrieval
- **Future**: Implement Redis caching for sent messages list with TTL and cache invalidation

#### 3. **Additional Improvements**
- **Metrics & Monitoring**: Add Prometheus metrics for request rates, error rates, processing times
- **Rate Limiting**: Implement per-client rate limiting for API endpoints
- **Database Migrations**: Add proper migration system for schema changes
- **Configuration Validation**: Validate configuration on startup
- **Circuit Breaker**: Add circuit breaker pattern for webhook calls
- **Message Queuing**: Consider message queue (RabbitMQ/Kafka) for high-volume scenarios

### Design Decisions

#### **Why Standard Logging?**
- **Simplicity**: Standard `log` package is sufficient for MVP and development
- **No Dependencies**: Avoids external logging dependencies
- **Easy Migration**: Can easily upgrade to structured logging later

#### **Why Limited Redis Usage?**
- **Requirements Clarity**: Project requirements didn't specify caching strategy for message lists
- **Data Consistency**: Database remains source of truth for message data
- **Implementation Complexity**: Full caching requires cache invalidation logic

## 📝 License

This project is part of the Insider Assessment.

---

**Built with ❤️ using Go, PostgreSQL, Redis, and Docker**
