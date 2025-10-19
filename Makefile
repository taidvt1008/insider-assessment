APP_NAME=insider-message-sender
GO_MAIN=cmd/server/main.go
SWAG_OUT=internal/docs

default: help

.PHONY: setup swag lint build run up down logs reset-db start-db clean redis db help

setup:
	@echo "📦 Installing dependencies..."
	go install github.com/swaggo/swag/cmd/swag@latest
	go mod tidy
	cp .env.example .env
	@echo "✅ Done."

swag:
	@echo "🧾 Generating Swagger docs..."
	swag init -g internal/api/router.go -o $(SWAG_OUT)
	@echo "✅ Swagger docs generated in $(SWAG_OUT)/"

lint:
	@echo "🔍 Running golangci-lint..."
	@golangci-lint run ./... --timeout 5m

build: swag
	@echo "🔨 Building binary..."
	go build -o $(APP_NAME) $(GO_MAIN)
	@echo "✅ Binary built: ./$(APP_NAME)"

run: start-db
	@echo "🚀 Running $(APP_NAME) locally..."
	go run $(GO_MAIN)

up:
	@echo "🐳 Starting full stack (Postgres, Redis, App)..."
	docker-compose --env-file .env.docker up --build

down:
	@echo "🧹 Stopping containers (keeping data)..."
	docker-compose down

logs:
	docker-compose logs -f --tail=50

reset-db:
	@echo "⚠️ Resetting PostgreSQL and Redis data (app will stay running)..."
	@echo "🧹 Clearing PostgreSQL schema..."
	docker-compose exec -T postgres psql -U postgres -d insider -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	docker-compose exec -T postgres psql -U postgres -d insider -f /docker-entrypoint-initdb.d/init.sql
	@echo "🧹 Flushing Redis cache..."
	docker-compose exec -T redis redis-cli FLUSHALL
	@echo "✅ PostgreSQL and Redis reset complete! App container is still running."


start-db:
	@echo "🐘 Starting PostgreSQL and ⚡ Redis in background..."
	docker-compose up -d postgres redis
	@echo ""
	@echo "✅ PostgreSQL and Redis are now running."
	@echo "👉 You can now run the app with: make run"
	@echo ""

stop-db:
	@echo "🧊 Stopping PostgreSQL and Redis containers..."
	docker-compose stop postgres redis
	@echo "✅ PostgreSQL and Redis stopped. App container (if any) still running."

clean:
	@echo "🧼 Cleaning build artifacts..."
	rm -f $(APP_NAME)
	rm -rf $(SWAG_OUT)/*
	@echo "✅ Clean complete."

redis:
	@echo "⚡ Connecting to Redis..."
	docker exec -it insider_redis redis-cli

db:
	@echo "🔍 Connecting to Postgres..."
	docker exec -it insider_postgres psql -U postgres -d insider

help:
	@echo "Available commands:"
	@echo "  make setup            - Install dependencies"
	@echo "  make swag             - Generate Swagger docs"
	@echo "  make lint             - Run golangci-lint to check code quality"
	@echo "  make build            - Build Go binary"
	@echo "  make run              - Run app locally (uses local Postgres/Redis)"
	@echo "  make up               - Build & run full stack (Docker)"
	@echo "  make down             - Stop containers but keep volumes"
	@echo "  make logs             - Tail Docker logs"
	@echo "  make start-db         - Start PostgreSQL + Redis (background) then prompt next step"
	@echo "  make stop-db          - Stop PostgreSQL + Redis"
	@echo "  make reset-db         - Reset PostgreSQL and Redis data (app will stay running)"
	@echo "  make clean            - Remove build artifacts"
	@echo "  make redis            - Connect to Redis"
	@echo "  make db               - Connect to Postgres"
	@echo ""
	@echo "Testing commands:"
	@echo "  make test-health      - Test health check endpoint"
	@echo "  make test-start       - Test scheduler start"
	@echo "  make test-stop        - Test scheduler stop"
	@echo "  make test-list-sent   - Test sent messages listing"
	@echo "  make test-list-failed - Test failed messages listing"

## 🧪 Quick API Tests
test-health:
	@echo "🏥 Testing health check endpoint..."
	@curl -s -X GET http://localhost:8080/health -H "Accept: application/json" | jq .

test-start:
	@echo "▶️ Testing scheduler START endpoint..."
	@curl -s -X POST http://localhost:8080/api/v1/scheduler/start -H "Content-Type: application/json" | jq .

test-stop:
	@echo "⏹️ Testing scheduler STOP endpoint..."
	@curl -s -X POST http://localhost:8080/api/v1/scheduler/stop -H "Content-Type: application/json" | jq .

test-list-sent:
	@echo "📬 Testing fetch sent messages endpoint..."
	@curl -s -X GET "http://localhost:8080/api/v1/messages/sent?limit=2&offset=1" -H "Accept: application/json" | jq .

test-list-failed:
	@echo "❌ Testing fetch failed messages endpoint..."
	@curl -s -X GET "http://localhost:8080/api/v1/messages/failed?limit=3" -H "Accept: application/json" | jq .

