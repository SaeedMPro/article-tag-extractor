# Go parameters
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

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test ./... -v

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