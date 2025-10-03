package app

import (
	"context"
	"sync"
	"time"

	"github.com/SaeedMPro/article-tag-extractor/internal/domain/entity"
	"github.com/SaeedMPro/article-tag-extractor/internal/domain/port"
)

type ArticleService struct {
	repo         port.ArticleRepository
	tagExtractor port.TagExtractor
}

func NewArticleService(repo port.ArticleRepository) *ArticleService {
	return &ArticleService{
		repo:         repo,
		tagExtractor: NewTagExtractorService(),
	}
}

func (s *ArticleService) ProcessArticles(ctx context.Context, articles []*entity.Article) (int, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	count := 0

	// extract tags concurrently and save articles with their tags
	for _, article := range articles {
		wg.Add(1)
		go func(a *entity.Article) {
			defer wg.Done()
			
			tags := s.tagExtractor.ExtractTags(a.Title, a.Body)
			a.Tags = tags

			article := &entity.Article{
				Title:     a.Title,
				Body:      a.Body,
				Tags:      tags,
				CreatedAt: time.Now(),
			}

			if err := s.repo.SaveArticle(ctx, article); err == nil {
				mu.Lock()
				count++
				mu.Unlock()
			}
		}(article)
	}
	wg.Wait()

	return count, nil
}

func (s *ArticleService) GetTopTags(ctx context.Context, limit int) ([]entity.TagFrequency, error) {
	return s.repo.GetTopTags(ctx, limit)
}
