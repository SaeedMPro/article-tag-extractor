package app

import (
	"context"
	"sync"
	"time"

	"github.com/SaeedMPro/article-tag-extractor/internal/domain/entity"
	"github.com/SaeedMPro/article-tag-extractor/internal/domain/port"
)

type ArticleService struct {
	Repo         port.ArticleRepository
	TagExtractor port.TagExtractor
}

func NewArticleService(repo port.ArticleRepository) *ArticleService {
	return &ArticleService{
		Repo:         repo,
		TagExtractor: NewTagExtractorService(),
	}
}

func NewArticleServiceWithExtractor(repo port.ArticleRepository, tagExtractor port.TagExtractor) *ArticleService {
	return &ArticleService{
		Repo:         repo,
		TagExtractor: tagExtractor,
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
			
			tags := s.TagExtractor.ExtractTags(a.Title, a.Body)
			a.Tags = tags

			article := &entity.Article{
				Title:     a.Title,
				Body:      a.Body,
				Tags:      tags,
				CreatedAt: time.Now(),
			}

			if err := s.Repo.SaveArticle(ctx, article); err == nil {
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
	return s.Repo.GetTopTags(ctx, limit)
}
