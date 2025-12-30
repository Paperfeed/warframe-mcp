.PHONY: build test test-verbose clean install

# Build the binary
build:
	go build -o warframe-mcp

# Build for Windows
build-windows:
	GOOS=windows GOARCH=amd64 go build -o warframe-mcp.exe

# Run all tests
test:
	go test ./...

# Run tests with verbose output
test-verbose:
	go test ./... -v

# Run tests with coverage
test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Run specific package tests
test-worldstate:
	go test ./internal/worldstate -v

test-market:
	go test ./internal/market -v

test-recommender:
	go test ./internal/recommender -v

test-config:
	go test ./internal/config -v

test-alecaframe:
	go test ./internal/alecaframe -v

# Test Alecaframe API integration
test-alecaframe-api:
	@echo "Testing Alecaframe API endpoints..."
	@if [ -z "$$ALECAFRAME_USER_HASH" ] || [ -z "$$ALECAFRAME_PUBLIC_TOKEN" ]; then \
		echo "Error: Set ALECAFRAME_USER_HASH and ALECAFRAME_PUBLIC_TOKEN environment variables"; \
		echo "See TESTING_ALECAFRAME.md for details"; \
		exit 1; \
	fi
	go run cmd/test-alecaframe/main.go -endpoint relics
	@echo ""
	go run cmd/test-alecaframe/main.go -endpoint stats

# Clean build artifacts
clean:
	rm -f warframe-mcp warframe-mcp.exe coverage.out coverage.html

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Run all checks (fmt, test, build)
check: fmt test build

# Help
help:
	@echo "Available targets:"
	@echo "  build                - Build the binary"
	@echo "  build-windows        - Build Windows binary"
	@echo "  test                 - Run all tests"
	@echo "  test-verbose         - Run tests with verbose output"
	@echo "  test-coverage        - Run tests with coverage report"
	@echo "  test-alecaframe-api  - Test Alecaframe API integration (requires credentials)"
	@echo "  clean                - Remove build artifacts"
	@echo "  deps                 - Download and tidy dependencies"
	@echo "  fmt                  - Format code"
	@echo "  lint                 - Run linter (requires golangci-lint)"
	@echo "  check                - Run fmt, test, and build"
	@echo "  help                 - Show this help message"
