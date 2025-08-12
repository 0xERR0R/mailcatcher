.PHONY: build test test-e2e clean docker-build docker-run

# Build the mailcatcher binary
build:
	go build -o bin/mailcatcher ./cmd/mailcatcher


# Run tests
test:
	docker build -t mailcatcher:test .
	go test ./...
	go test -v ./tests/e2e/... -timeout 5m


# Clean build artifacts
clean:
	rm -rf bin/
	go clean -cache

# Build Docker image
docker-build:
	docker build -t mailcatcher:latest .

# Run mailcatcher with Docker Compose
docker-run:
	docker-compose up -d

# Stop mailcatcher with Docker Compose
docker-stop:
	docker-compose down

# Show logs
logs:
	docker-compose logs -f

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run


# Help target
help:
	@echo "Available targets:"
	@echo "  build           - Build the mailcatcher binary"
	@echo "  test            - Run unit tests"
	@echo "  clean           - Clean build artifacts"
	@echo "  docker-build    - Build Docker image"
	@echo "  docker-run      - Run with Docker Compose"
	@echo "  docker-stop     - Stop Docker Compose"
	@echo "  logs            - Show logs"
	@echo "  deps            - Install dependencies"
	@echo "  fmt             - Format code"
	@echo "  lint            - Lint code"
	@echo "  help            - Show this help"
