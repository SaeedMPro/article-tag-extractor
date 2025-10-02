package entity

import (
	"time"
)

type Article struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Title     string    `bson:"title" json:"title"`
	Body      string    `bson:"body" json:"body"`
	Tags      []string  `bson:"tags" json:"tags"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

type TagFrequency struct {
	Tag       string `bson:"_id" json:"tag"`
	Frequency int    `bson:"frequency" json:"frequency"`
}

type ProcessArticleRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}
