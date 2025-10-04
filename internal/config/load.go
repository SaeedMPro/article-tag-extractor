package config

import (
	"os"
)

func LoadConfig() *Config {
	return &Config{
		Database: Database{
			URI:        getEnv("MONGODB_URL", "mongodb://localhost:27017"),
			DBName:     getEnv("MONGODB_DB_NAME", "article_db"),
		},
		Server: Server{
			GRPCPort: getEnv("GRPC_SERVER_PORT", "50051"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
