package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/SaeedMPro/article-tag-extractor/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	Conn   *mongo.Client
	DBName string
}

func NewClient(dbConfig config.Database) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// create mongo client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbConfig.URL))
	if err != nil {
		return nil, fmt.Errorf("mongo connect error: %w", err)
	}

	// ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("mongo ping error: %w", err)
	}

	return &Client{
		Conn:   client,
		DBName: dbConfig.DBName,
	}, nil
}

func (c *Client) Collection(name string) *mongo.Collection {
	return c.Conn.Database(c.DBName).Collection(name)
}

func (c *Client) Disconnect(ctx context.Context) error {
	return c.Conn.Disconnect(ctx)
}
