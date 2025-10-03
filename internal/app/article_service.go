package app

import (
	"context"

	"github.com/SaeedMPro/article-tag-extractor/internal/domain/entity"
	"github.com/SaeedMPro/article-tag-extractor/internal/domain/port"
)

type ArticleService struct {
	repo port.ArticleRepository
}

func NewArticleService(repo port.ArticleRepository) *ArticleService {
	return &ArticleService{repo: repo}
}

func (s *ArticleService) ProcessArticles(ctx context.Context, articles []*entity.Article) (int, error) {
	count := 0

	//TODO: implement processing article

	return count, nil
}

func (s *ArticleService) GetTopTags(ctx context.Context, limit int) ([]entity.TagFrequency, error) {
	return s.repo.GetTopTags(ctx, limit)
}
