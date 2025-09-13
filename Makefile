# Makefile for Zeus project

.PHONY: help build run stop clean dev prod logs

# Default target
help:
	@echo "Available commands:"
	@echo "  dev     - Start development environment"
	@echo "  prod    - Start production environment"
	@echo "  build   - Build Docker image"
	@echo "  run     - Run the application"
	@echo "  stop    - Stop all services"
	@echo "  clean   - Clean up containers and images"
	@echo "  logs    - Show logs"
	@echo "  test    - Run tests"

# Development environment
dev:
	docker-compose up --build

# Production environment
prod:
	docker-compose -f docker-compose.prod.yml up --build -d

# Build Docker image
build:
	docker build -t zeus:latest .

# Build production image
build-prod:
	docker build -f Dockerfile.prod -t zeus:prod .

# Run the application
run:
	docker-compose up

# Stop all services
stop:
	docker-compose down
	docker-compose -f docker-compose.prod.yml down

# Clean up
clean:
	docker-compose down -v --rmi all
	docker-compose -f docker-compose.prod.yml down -v --rmi all
	docker system prune -f

# Show logs
logs:
	docker-compose logs -f

# Show production logs
logs-prod:
	docker-compose -f docker-compose.prod.yml logs -f

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Install dependencies
deps:
	go mod download
	go mod tidy
