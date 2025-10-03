package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SaeedMPro/article-tag-extractor/internal/app"
	"github.com/SaeedMPro/article-tag-extractor/internal/config"
	"github.com/SaeedMPro/article-tag-extractor/internal/infra/grpc"
	"github.com/SaeedMPro/article-tag-extractor/internal/infra/mongodb"
)

func main() {
	// load config
	cfg := config.LoadConfig()
	log.Printf("config loaded: %v", cfg)

	// connect to mongo
	db, err := mongodb.NewClient(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := db.Disconnect(ctx); err != nil {
			log.Printf("mongo disconnect error: %v", err)
		}
	}()

	// create repo & service & grpc server
	articleRepo := mongodb.NewArticleRepository(db.Conn, cfg.Database.DBName, "articles")
	articleService := app.NewArticleService(articleRepo)
	grpcServer := grpc.NewServer(articleService)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.GRPCPort))
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", cfg.Server.GRPCPort, err)
	}

	// run gRPC server
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("gRPC server running on port %s", cfg.Server.GRPCPort)
		if err := grpcServer.Serve(listener); err != nil {
			serverErr <- err
		}
	}()

	// handle gracefully shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case <-stop:
		log.Println("shutting down gracefully...")
		grpcServer.GracefulStop()
	case err := <-serverErr:
		log.Fatalf("grpc server error: %v", err)
	}

	log.Println("Server stopped gracefully")
}
