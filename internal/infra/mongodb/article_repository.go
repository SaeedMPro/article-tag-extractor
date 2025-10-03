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
	_, err := r.collection.InsertOne(ctx, article)
	return err
}

func (r *ArticleRepository) GetTopTags(ctx context.Context, limit int) ([]entity.TagFrequency, error) {
	pipeline := mongo.Pipeline{
		//unwind the tags array:
		{{Key: "$unwind", Value: "$tags"}},

		//group by tags and calculate frequency:
		{{
			Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$tags"},
				{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
			},
		}},

		//sort by frequency in desc order:
		{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},

		//limit the results:
		{{Key: "$limit", Value: limit}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var topTags []entity.TagFrequency
	for cursor.Next(ctx) {
		var result struct {
			Tag   string `bson:"_id"`
			Count int    `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		topTags = append(topTags, entity.TagFrequency{
			Tag:       result.Tag,
			Frequency: result.Count,
		})
	}
	return topTags, nil
}
