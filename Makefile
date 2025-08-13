.PHONY: build test lint install clean help

# Build the binary
build:
	go build -o bin/goroutine-defer-guard ./cmd/goroutine-defer-guard

# Install the binary globally
	go install ./cmd/goroutine-defer-guard

# Run tests
test:
	go test ./...

# Run the linter on this project
lint:
	go run ./cmd/goroutine-defer-guard -root="$(shell pwd)" ./...

# Run the linter with verbose output
lint-verbose:
	go run ./cmd/goroutine-defer-guard -root="$(shell pwd)" -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Initialize dependencies
deps:
	go mod tidy
	go mod download

# Show help
help:
	@echo "Available targets:"
	@echo "  build        - Build the binary"
	@echo "  install      - Install the binary globally"
	@echo "  test         - Run all tests"
	@echo "  lint         - Run the linter on this project"
	@echo "  lint-verbose - Run the linter with verbose output"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Initialize dependencies"
	@echo "  help         - Show this help message"
