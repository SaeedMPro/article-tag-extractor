package mongodb

import (
	"context"
	"testing"
	"time"

	"github.com/SaeedMPro/article-tag-extractor/internal/domain/entity"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoUri = "mongodb://localhost:27017"

// Integration test helper (requires actual MongoDB instance)
func TestArticleRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This test requires a real MongoDB instance
	// You can run it with: go test -v -run TestArticleRepository_Integration
	// Make sure MongoDB is running on localhost:27017

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		t.Skipf("Skipping integration test: cannot connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	repo := NewArticleRepository(client, "test_db", "test_collection")

	// Test saving an article
	article := &entity.Article{
		Title:     "Integration Test Article",
		Body:      "This is an integration test article",
		Tags:      []string{"integration", "test", "article"},
		CreatedAt: time.Now(),
	}

	err = repo.SaveArticle(context.Background(), article)
	if err != nil {
		t.Errorf("Failed to save article: %v", err)
	}

	// Test getting top tags
	tags, err := repo.GetTopTags(context.Background(), 5)
	if err != nil {
		t.Errorf("Failed to get top tags: %v", err)
	}

	t.Logf("Retrieved %d tags from integration test", len(tags))

	// Clean up
	client.Database("test_db").Collection("test_collection").Drop(context.Background())
}
