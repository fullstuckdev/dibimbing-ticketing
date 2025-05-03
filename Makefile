# Makefile for Ticketing System

.PHONY: run build test dev clean help

# Default target
help:
	@echo "Available commands:"
	@echo "  make run      - Run the application"
	@echo "  make dev      - Run with hot reload using Air"
	@echo "  make build    - Build the application"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make help     - Show this help message"

# Run the application
run:
	go run main.go

# Run with hot reload using Air
dev:
	$(shell go env GOPATH)/bin/air

# Build the application
build:
	go build -o ticketing-system main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f ticketing-system
	rm -rf tmp/ 