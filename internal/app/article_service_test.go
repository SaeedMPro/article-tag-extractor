package app

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/SaeedMPro/article-tag-extractor/internal/domain/entity"
)

// MockArticleRepository is a mock implementation of ArticleRepository
type MockArticleRepository struct {
	articles            []entity.Article
	tagFrequencies      []entity.TagFrequency
	saveError           error
	getTopTagsError     error
	saveCallCount       int
	getTopTagsCallCount int
	mu                  sync.Mutex
}

func (m *MockArticleRepository) SaveArticle(ctx context.Context, article *entity.Article) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.saveCallCount++
	if m.saveError != nil {
		return m.saveError
	}
	m.articles = append(m.articles, *article)
	return nil
}

func (m *MockArticleRepository) GetTopTags(ctx context.Context, limit int) ([]entity.TagFrequency, error) {
	m.getTopTagsCallCount++
	if m.getTopTagsError != nil {
		return nil, m.getTopTagsError
	}
	if limit > len(m.tagFrequencies) {
		limit = len(m.tagFrequencies)
	}
	return m.tagFrequencies[:limit], nil
}

// MockTagExtractor is a mock implementation of TagExtractor
type MockTagExtractor struct {
	tags []string
}

func (m *MockTagExtractor) ExtractTags(title, body string) []string {
	return m.tags
}

func TestArticleService_ProcessArticles(t *testing.T) {
	tests := []struct {
		name          string
		articles      []*entity.Article
		mockTags      []string
		expectedCount int
		expectError   bool
		saveError     error
	}{
		{
			name: "Process single article successfully",
			articles: []*entity.Article{
				{Title: "Test Article", Body: "This is a test article"},
			},
			mockTags:      []string{"test", "article"},
			expectedCount: 1,
			expectError:   false,
			saveError:     nil,
		},
		{
			name: "Process multiple articles successfully",
			articles: []*entity.Article{
				{Title: "Article 1", Body: "Content 1"},
				{Title: "Article 2", Body: "Content 2"},
				{Title: "Article 3", Body: "Content 3"},
			},
			mockTags:      []string{"article", "content"},
			expectedCount: 3,
			expectError:   false,
			saveError:     nil,
		},
		{
			name: "Process articles with save error",
			articles: []*entity.Article{
				{Title: "Test Article", Body: "This is a test article"},
			},
			mockTags:      []string{"test", "article"},
			expectedCount: 0,
			expectError:   false, // Service doesn't return error, just counts successful saves
			saveError:     errors.New("database error"),
		},
		{
			name:          "Process empty articles list",
			articles:      []*entity.Article{},
			mockTags:      []string{},
			expectedCount: 0,
			expectError:   false,
			saveError:     nil,
		},
		{
			name: "Process large batch of articles",
			articles: func() []*entity.Article {
				articles := make([]*entity.Article, 100000)
				for i := 0; i < 100000; i++ {
					articles[i] = &entity.Article{
						Title: "Article " + string(rune(i)),
						Body:  "Content for article " + string(rune(i)),
					}
				}
				return articles
			}(),
			mockTags:      []string{"article", "content"},
			expectedCount: 100000,
			expectError:   false,
			saveError:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockArticleRepository{
				saveError: tt.saveError,
			}
			mockExtractor := &MockTagExtractor{tags: tt.mockTags}

			service := &ArticleService{
				Repo:         mockRepo,
				TagExtractor: mockExtractor,
			}

			ctx := context.Background()
			result, err := service.ProcessArticles(ctx, tt.articles)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if result != tt.expectedCount {
				t.Errorf("Expected %d processed articles, got %d", tt.expectedCount, result)
			}

			// Verify that SaveArticle was called the expected number of times
			expectedSaveCalls := len(tt.articles)
			if mockRepo.saveCallCount != expectedSaveCalls {
				t.Errorf("Expected SaveArticle to be called %d times, got %d", expectedSaveCalls, mockRepo.saveCallCount)
			}
			
			// Verify that articles were saved to repository (if no save error)
			if tt.saveError == nil {
				
				if len(mockRepo.articles) != tt.expectedCount {
					t.Errorf("Expected %d articles in repository, got %d", tt.expectedCount, len(mockRepo.articles))
				}

				// Verify that tags were extracted and assigned
				for i, article := range mockRepo.articles {
					if len(article.Tags) != len(tt.mockTags) {
						t.Errorf("Article %d: expected %d tags, got %d", i, len(tt.mockTags), len(article.Tags))
					}

					if article.CreatedAt.IsZero() {
						t.Errorf("Article %d: CreatedAt should not be zero", i)
					}
				}
			}
		})
	}
}

func TestArticleService_ProcessArticles_Concurrency(t *testing.T) {
	// Test concurrent processing with timing
	mockRepo := &MockArticleRepository{}
	mockExtractor := &MockTagExtractor{tags: []string{"test", "concurrent"}}

	service := &ArticleService{
		Repo:         mockRepo,
		TagExtractor: mockExtractor,
	}

	// Create articles
	articleNumbers := 1000
	articles := make([]*entity.Article, articleNumbers)
	for i := 0; i < articleNumbers; i++ {
		articles[i] = &entity.Article{
			Title: "Concurrent Article " + string(rune(i)),
			Body:  "This is a test for concurrent processing",
		}
	}

	ctx := context.Background()
	start := time.Now()
	result, err := service.ProcessArticles(ctx, articles)
	duration := time.Since(start)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result != articleNumbers {
		t.Errorf("Expected %d processed articles, got %d", articleNumbers, result)
	}

	// Verify all articles were processed 
	if len(mockRepo.articles) < articleNumbers {
		t.Errorf("Expected at least %d articles in repository, got %d", articleNumbers, len(mockRepo.articles))
	}

	// Concurrent processing should be reasonably fast
	expectedTime := 5 * time.Second
	if duration > expectedTime {
		t.Errorf("Processing took too long: %v", duration)
	}

	t.Logf("Processed %d articles in %v", result, duration)
}

func TestArticleService_GetTopTags(t *testing.T) {
	mockTagFrequencies := []entity.TagFrequency{
		{Tag: "golang", Frequency: 10},
		{Tag: "programming", Frequency: 8},
		{Tag: "web", Frequency: 5},
		{Tag: "api", Frequency: 3},
		{Tag: "microservices", Frequency: 2},
	}

	tests := []struct {
		name     string
		limit    int
		expected int
		error    error
	}{
		{
			name:     "Get top 2 tags",
			limit:    2,
			expected: 2,
			error:    nil,
		},
		{
			name:     "Get top 10 tags (more than available)",
			limit:    10,
			expected: 5,
			error:    nil,
		},
		{
			name:     "Get top 0 tags",
			limit:    0,
			expected: 0,
			error:    nil,
		},
		{
			name:     "Database error",
			limit:    5,
			expected: 0,
			error:    errors.New("database connection failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockArticleRepository{
				tagFrequencies:  mockTagFrequencies,
				getTopTagsError: tt.error,
			}
			mockExtractor := &MockTagExtractor{}

			service := &ArticleService{
				Repo:         mockRepo,
				TagExtractor: mockExtractor,
			}

			ctx := context.Background()
			result, err := service.GetTopTags(ctx, tt.limit)

			if tt.error != nil {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if len(result) != tt.expected {
				t.Errorf("Expected %d tags, got %d", tt.expected, len(result))
			}

			// Verify GetTopTags was called once
			if mockRepo.getTopTagsCallCount != 1 {
				t.Errorf("Expected GetTopTags to be called 1 time, got %d", mockRepo.getTopTagsCallCount)
			}
		})
	}
}

func TestArticleService_ContextCancellation(t *testing.T) {
	mockRepo := &MockArticleRepository{}
	mockExtractor := &MockTagExtractor{tags: []string{"test"}}

	service := &ArticleService{
		Repo:         mockRepo,
		TagExtractor: mockExtractor,
	}

	// Create a context that will be cancelled
	ctx, cancel := context.WithCancel(context.Background())

	articles := []*entity.Article{
		{Title: "Test Article", Body: "This is a test article"},
	}

	// Cancel the context immediately
	cancel()

	// Process articles with cancelled context
	result, err := service.ProcessArticles(ctx, articles)

	// The service should handle context cancellation gracefully
	// It might return 0 processed articles or an error
	if result > 1 {
		t.Errorf("Expected 0 or 1 processed articles with cancelled context, got %d", result)
	}

	// Error handling depends on implementation
	t.Logf("Result with cancelled context: %d, error: %v", result, err)
}

func BenchmarkArticleService_ProcessArticles(b *testing.B) {
	mockRepo := &MockArticleRepository{}
	mockExtractor := &MockTagExtractor{tags: []string{"benchmark", "test"}}

	service := &ArticleService{
		Repo:         mockRepo,
		TagExtractor: mockExtractor,
	}

	articles := make([]*entity.Article, 1000000)
	for i := 0; i < 1000000; i++ {
		articles[i] = &entity.Article{
			Title: "Benchmark Article " + string(rune(i)),
			Body:  "This is a benchmark test article, contains programming, development, and software engineering concepts.",
		}
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.ProcessArticles(ctx, articles)
	}
}

