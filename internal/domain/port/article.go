package port

import (
	"context"

	"github.com/SaeedMPro/article-tag-extractor/internal/domain/entity"
)

// articleRepository defines the interface for article data operations
type ArticleRepository interface {
	SaveArticle(ctx context.Context, article *entity.Article) error
	GetTopTags(ctx context.Context, limit int) ([]entity.TagFrequency, error)
}

// tagExtractor defines the interface for tag extraction logic
type TagExtractor interface {
	ExtractTags(title, body string) []string
}

// articleService defines the interface for article business logic
type ArticleService interface {
	ProcessArticles(ctx context.Context, articles []entity.ProcessArticleRequest) (int, error)
	GetTopTags(ctx context.Context, limit int) ([]entity.TagFrequency, error)
}
