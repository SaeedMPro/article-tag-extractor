package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	Conn   *mongo.Client
	DBName string
}

func NewClient(uri, dbName string) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("mongo connect error: %w", err)
	}

	return &Client{
		Conn:   client,
		DBName: dbName,
	}, nil
}

func (c *Client) Collection(name string) *mongo.Collection {
	return c.Conn.Database(c.DBName).Collection(name)
}

func (c *Client) Disconnect(ctx context.Context) error {
	return c.Conn.Disconnect(ctx)
}
