package mongodb

import (
	"context"
	"log"
	"time"

	"github.com/SaeedMPro/article-tag-extractor/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ArticleRepository struct {
	collection *mongo.Collection
}

func NewArticleRepository(client *mongo.Client, dbName, collectionName string) *ArticleRepository {
	db := client.Database(dbName)
	coll := db.Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "tags", Value: 1}},
	})
	if err != nil {
		log.Printf("failed to create index: %v\n", err)
	}

	return &ArticleRepository{
		collection: coll,
	}
}

func (r *ArticleRepository) SaveArticle(ctx context.Context, article *entity.Article) error {
	//TODO
	return nil
}

func (r *ArticleRepository) GetTopTags(ctx context.Context, limit int) ([]entity.TagFrequency, error) {
	//TODO
	return nil, nil
}
