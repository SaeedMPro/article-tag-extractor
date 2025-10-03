package mongodb

import (
	"context"

	"github.com/SaeedMPro/article-tag-extractor/internal/domain/entity"
	"go.mongodb.org/mongo-driver/mongo"
)

type ArticleRepository struct {
	collection *mongo.Collection
}

func NewArticleRepository(client *mongo.Client, dbName string, collection string) *ArticleRepository {
	return &ArticleRepository{
		collection: client.Database(dbName).Collection(collection),
	}
}

func (r *ArticleRepository) SaveArticle(ctx context.Context, a *entity.Article) error {
	//TODO
	return nil
}

func (r *ArticleRepository) GetTopTags(ctx context.Context, limit int) ([]entity.TagFrequency, error) {
	//TODO
	return nil, nil
}
