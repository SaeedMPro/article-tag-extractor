# parameters
APP_NAME := article-tag-extractor
BIN_DIR := bin
PROTO_DIR := internal/proto
GEN_DIR := internal/pb

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BIN_DIR)/$(APP_NAME) ./cmd/...

# Run the app
.PHONY: run
run: build
	@echo "Running $(APP_NAME)..."
	@./$(BIN_DIR)/$(APP_NAME)

# Run all tests
.PHONY: test
test:
	@echo "Running all tests..."
	@go test ./... -v

.PHONY: test-unit
test-unit:
	@echo "Running unit tests..."
	@go test ./... -v -short
	@echo "Unit tests completed"

# Run specific package tests
.PHONY: test-app
test-app:
	@echo "Running app package tests..."
	@go test ./internal/app -v

.PHONY: test-utils
test-utils:
	@echo "Running utils package tests..."
	@go test ./utils -v

.PHONY: test-mongodb
test-mongodb:
	@echo "Running MongoDB package tests..."
	@go test ./internal/infra/mongodb -v

.PHONY: test-grpc
test-grpc:
	@echo "Running gRPC package tests..."
	@go test ./internal/infra/grpc -v

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf $(BIN_DIR)/*
	@rm -rf $(GEN_DIR)/*.pb.go $(GEN_DIR)/*.grpc.go

# Generate protobuf code
.PHONY: proto
proto:
	@echo "Generating protobuf code..."
	@protoc --go_out=$(GEN_DIR) --go-grpc_out=$(GEN_DIR) $(PROTO_DIR)/*.proto

# Start Docker Compose services
.PHONY: docker-up
docker-up:
	@echo "Starting services with Docker Compose..."
	docker-compose up --build

# Stop Docker Compose services
.PHONY: docker-down
docker-down:
	@echo "Stopping Docker Compose services..."
	docker-compose down

# Install Go dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Install protoc and protoc-gen-go
.PHONY: install-tools
install-tools:
	@echo "Installing protobuf tools..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	protoc --version

# Development workflow
.PHONY: dev
dev: deps proto build test-unit

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all              - Build everything (default)"
	@echo "  build            - Build the application"
	@echo "  run              - Build and run the application"
	@echo "  test             - Run all tests"
	@echo "  test-unit        - Run unit tests only"
	@echo "  test-app         - Run app package tests"
	@echo "  test-utils       - Run utils package tests"
	@echo "  test-mongodb     - Run MongoDB package tests"
	@echo "  test-grpc        - Run gRPC package tests"
	@echo "  clean            - Clean build artifacts"
	@echo "  proto            - Generate protobuf code"
	@echo "  docker-up        - Start services with Docker Compose"
	@echo "  docker-down      - Stop Docker Compose services"
	@echo "  deps             - Install dependencies"
	@echo "  install-tools    - Install protobuf tools"
	@echo "  dev              - Development workflow"
	@echo "  help             - Show this help"