package main

import (
	"context"
	"log"
	"time"

	"github.com/SaeedMPro/article-tag-extractor/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	cfg := config.LoadConfig()
	log.Printf("config loaded: %v", cfg)

	mongoClient, err := mongo.NewClient(cfg.Database.URL)
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}

	// TODO: grpc server connection and graceful shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := mongoClient.Disconnect(ctx); err != nil {
		log.Printf("Mongo disconnect error: %v", err)
	}

	log.Println("Server stopped gracefully")
}
