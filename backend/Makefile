.PHONY: build run dev test clean tidy

# Binary name
BINARY=api

# Build the application
build:
	go build -o bin/$(BINARY) ./cmd/api

# Run the application
run: build
	./bin/$(BINARY)

# Run with hot reload (requires air: go install github.com/air-verse/air@latest)
dev:
	air

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Tidy dependencies
tidy:
	go mod tidy

# Download dependencies
deps:
	go mod download

# Generate sqlc code (when you add sqlc)
sqlc:
	sqlc generate

# Lint (requires golangci-lint)
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...
