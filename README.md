# Article Tag Extractor

A high-performance, concurrent microservice for extracting keywords (tags) from article content using gRPC and MongoDB. The service processes articles in parallel, extracts meaningful tags using advanced text processing, and provides tag frequency analytics.

## Features

- **Concurrent Processing**: Process multiple articles simultaneously using goroutines
- **Tag Extraction**: Text normalization and stop-word filtering
- **gRPC API**: High-performance gRPC interface for microservice communication
- **MongoDB Integration**: Efficient storage and aggregation for tag analytics
- **Docker Support**: Complete containerization with Docker Compose
- **Comprehensive Testing**: Proper unit tests
- **Graceful Shutdown**: Proper resource cleanup and connection management

## Requirements

- Go 1.24+
- MongoDB 4.4+
- Docker & Docker Compose (optional)
- protoc (for protocol buffer generation)

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   gRPC Client   │───▶│  gRPC Server    │───▶│  Article Service│
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                        │
                       ┌─────────────────┐              │
                       │  MongoDB        │◀─────────────┘
                       │  - Articles     │
                       │  - Tag Freq     │
                       └─────────────────┘
```

### Project Structure

```
article-tag-extractor/
├── cmd/                    # Application entry point
│   └── main.go
├── internal/
│   ├── app/               # Business logic layer
│   │   ├── article_service.go
│   │   ├── tag_extractor_service.go
│   │   └── *_test.go
│   ├── config/            # Configuration management
│   ├── domain/            # Domain entities and interfaces
│   │   ├── entity/
│   │   └── port/
│   ├── infra/             # Infrastructure layer
│   │   ├── grpc/          # gRPC server implementation
│   │   └── mongodb/       # MongoDB repository
│   └── proto/             # Protocol buffer definitions
├── utils/                 # Utility functions
├── docker-compose.yaml    # Docker services
├── Dockerfile            # Container definition
├── Makefile              # Build automation
└── README.md
```

## Quick Start

### Using Docker Compose (Recommended)

```bash
# Clone the repository
git clone <repository-url>
cd article-tag-extractor

# Start all services
make docker-up

# The service will be available at localhost:50051
```

### Manual Setup

```bash
# Install dependencies
make deps

# Generate protocol buffer code
make proto

# Build the application
make build

# Start MongoDB (if not using Docker)
# Make sure MongoDB is running on localhost:27017

# Run the application
make run
```

## Configuration

The service can be configured using environment variables:

```bash
# MongoDB Configuration
export MONGODB_URL="mongodb://localhost:27017"
export MONGODB_DB_NAME="article_db"

# Server Configuration
export GRPC_SERVER_PORT="50051"
```

## API Usage

### gRPC Service Definition

```protobuf
service ArticleService {
  rpc ProcessArticles(ProcessArticlesRequest) returns (ProcessArticlesResponse);
  rpc GetTopTags(GetTopTagsRequest) returns (GetTopTagsResponse);
}
```

### Example Usage with grpcurl

```bash
# Process articles
grpcurl -plaintext -d '{
  "articles": [
    {
      "title": "Go Programming Language",
      "body": "Go is a programming language developed by Google for building scalable applications."
    },
    {
      "title": "Microservices Architecture", 
      "body": "Microservices enable building distributed systems using containerization and orchestration."
    }
  ]
}' localhost:50051 article.ArticleService/ProcessArticles

# Get top tags
grpcurl -plaintext -d '{"limit": 5}' localhost:50051 article.ArticleService/GetTopTags
```

## Testing

### Run Tests

```bash
# Run all tests
make test

# Run unit tests only (fast)
make test-unit

# Run specific package tests
make test-app
make test-utils
make test-mongodb # needed mongodb on mongodb://localhost:27017 (you can change uri in unit test)
make test-grpc
```

### Available Make Commands

```bash
make help              # Show all available commands
...
```

## 🔍 Tag Extraction Algorithm

The service uses a sophisticated tag extraction algorithm:

1. **Text Normalization**:
   - Convert to lowercase
   - Remove punctuation and special characters
   - Split into words

2. **Stop-word Filtering**:
   - Remove common English stop-words (the, and, is, etc.)
   - Filter out words shorter than 3 characters

3. **Frequency Analysis**:
   - Count word occurrences

4. **Concurrent Processing**:
   - Each article processed in separate goroutine
   - Parallel tag extraction and database storage


## Troubleshooting

### Common Issues

- **Protocol Buffer Generation Failed**
   ```bash
   # Install protoc tools
   make install-tools
   
   # Regenerate protobuf code
   make proto
   ```

## Acknowledgments

- [gRPC](https://grpc.io/) for high-performance RPC framework
- [MongoDB](https://www.mongodb.com/) for document database
- [Go](https://golang.org/) for the programming language
- [Protocol Buffers](https://developers.google.com/protocol-buffers) for serialization



