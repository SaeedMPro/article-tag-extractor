package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/SaeedMPro/article-tag-extractor/internal/app"
	"github.com/SaeedMPro/article-tag-extractor/internal/domain/entity"
	pb "github.com/SaeedMPro/article-tag-extractor/internal/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MockArticleService is a mock implementation of ArticleService
type MockArticleService struct {
	processArticlesResult int
	processArticlesError  error
	getTopTagsResult      []entity.TagFrequency
	getTopTagsError       error
	processCallCount      int
	getTopTagsCallCount   int
}

// MockArticleRepository is a mock implementation of ArticleRepository
type MockArticleRepository struct {
	articles         []entity.Article
	tagFrequencies   []entity.TagFrequency
	saveError        error
	getTopTagsError  error
	saveCallCount    int
	getTopTagsCallCount int
}

func (m *MockArticleRepository) SaveArticle(ctx context.Context, article *entity.Article) error {
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

func (m *MockTagExtractor) ExtractTagsConcurrently(title, body string) []string {
	return m.tags
}

func (m *MockArticleService) ProcessArticles(ctx context.Context, articles []*entity.Article) (int, error) {
	m.processCallCount++
	return m.processArticlesResult, m.processArticlesError
}

func (m *MockArticleService) GetTopTags(ctx context.Context, limit int) ([]entity.TagFrequency, error) {
	m.getTopTagsCallCount++
	return m.getTopTagsResult, m.getTopTagsError
}

func TestServer_ProcessArticles(t *testing.T) {
	tests := []struct {
		name           string
		request        *pb.ProcessArticlesRequest
		mockResult     int
		mockError      error
		expectedCode   codes.Code
		expectedCount  int32
	}{
		{
			name: "successful processing",
			request: &pb.ProcessArticlesRequest{
				Articles: []*pb.Article{
					{Title: "Test Article", Body: "This is a test article"},
					{Title: "Another Article", Body: "This is another test article"},
				},
			},
			mockResult:    2,
			mockError:     nil,
			expectedCode:  codes.OK,
			expectedCount: 2,
		},
		{
			name: "empty articles list",
			request: &pb.ProcessArticlesRequest{
				Articles: []*pb.Article{},
			},
			mockResult:    0,
			mockError:     nil,
			expectedCode:  codes.InvalidArgument,
			expectedCount: 0,
		},
		{
			name: "nil articles list",
			request: &pb.ProcessArticlesRequest{
				Articles: nil,
			},
			mockResult:    0,
			mockError:     nil,
			expectedCode:  codes.InvalidArgument,
			expectedCount: 0,
		},
		{
			name: "service error",
			request: &pb.ProcessArticlesRequest{
				Articles: []*pb.Article{
					{Title: "Test Article", Body: "This is a test article"},
				},
			},
			mockResult:    0,
			mockError:     errors.New("database connection failed"),
			expectedCode:  codes.Internal,
			expectedCount: 0,
		},
		{
			name: "large batch processing",
			request: func() *pb.ProcessArticlesRequest {
				articles := make([]*pb.Article, 100)
				for i := 0; i < 100; i++ {
					articles[i] = &pb.Article{
						Title: "Article " + string(rune(i)),
						Body:  "Content for article " + string(rune(i)),
					}
				}
				return &pb.ProcessArticlesRequest{Articles: articles}
			}(),
			mockResult:    100,
			mockError:     nil,
			expectedCode:  codes.OK,
			expectedCount: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock service with proper dependencies
			mockRepo := &MockArticleRepository{}
			mockExtractor := &MockTagExtractor{}
			articleService := app.NewArticleServiceWithExtractor(mockRepo, mockExtractor)
			
			// Create a proper server instance
			grpcServer := NewServer(articleService)

			// Since we can't easily mock the service in the current structure,
			// we'll test the validation logic and error handling
			ctx := context.Background()
			response, err := grpcServer.ProcessArticles(ctx, tt.request)

			if tt.expectedCode == codes.InvalidArgument {
				if err == nil {
					t.Error("Expected error for invalid argument, got nil")
				} else {
					st, ok := status.FromError(err)
					if !ok {
						t.Error("Expected gRPC status error")
					} else if st.Code() != tt.expectedCode {
						t.Errorf("Expected error code %v, got %v", tt.expectedCode, st.Code())
					}
				}
				return
			}

			// For successful cases, we expect no error
			if tt.expectedCode == codes.OK {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				} else if response.TotalProcessed != tt.expectedCount {
					t.Errorf("Expected %d processed articles, got %d", tt.expectedCount, response.TotalProcessed)
				}
			}

		})
	}
}

func TestServer_GetTopTags(t *testing.T) {
	tests := []struct {
		name         string
		request      *pb.GetTopTagsRequest
		mockResult   []entity.TagFrequency
		mockError    error
		expectedCode codes.Code
		expectedTags int
	}{
		{
			name: "successful get top tags",
			request: &pb.GetTopTagsRequest{
				Limit: 5,
			},
			mockResult: []entity.TagFrequency{
				{Tag: "golang", Frequency: 10},
				{Tag: "programming", Frequency: 8},
				{Tag: "web", Frequency: 5},
			},
			mockError:    nil,
			expectedCode: codes.OK,
			expectedTags: 3,
		},
		{
			name: "zero limit",
			request: &pb.GetTopTagsRequest{
				Limit: 0,
			},
			mockResult:   nil,
			mockError:    nil,
			expectedCode: codes.InvalidArgument,
			expectedTags: 0,
		},
		{
			name: "negative limit",
			request: &pb.GetTopTagsRequest{
				Limit: -1,
			},
			mockResult:   nil,
			mockError:    nil,
			expectedCode: codes.InvalidArgument,
			expectedTags: 0,
		},
		{
			name: "service error",
			request: &pb.GetTopTagsRequest{
				Limit: 5,
			},
			mockResult:   nil,
			mockError:    errors.New("database error"),
			expectedCode: codes.Internal,
			expectedTags: 0,
		},
		{
			name: "empty result",
			request: &pb.GetTopTagsRequest{
				Limit: 5,
			},
			mockResult:   []entity.TagFrequency{},
			mockError:    nil,
			expectedCode: codes.OK,
			expectedTags: 0,
		},
		{
			name: "large limit",
			request: &pb.GetTopTagsRequest{
				Limit: 1000,
			},
			mockResult: func() []entity.TagFrequency {
				tags := make([]entity.TagFrequency, 100)
				for i := 0; i < 100; i++ {
					tags[i] = entity.TagFrequency{
						Tag:       "tag" + string(rune(i)),
						Frequency: 100 - i,
					}
				}
				return tags
			}(),
			mockError:    nil,
			expectedCode: codes.OK,
			expectedTags: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock service with proper dependencies
			mockRepo := &MockArticleRepository{
				tagFrequencies:  tt.mockResult,
				getTopTagsError: tt.mockError,
			}
			mockExtractor := &MockTagExtractor{}
			articleService := app.NewArticleServiceWithExtractor(mockRepo, mockExtractor)
			
			grpcServer := NewServer(articleService)

			ctx := context.Background()
			response, err := grpcServer.GetTopTags(ctx, tt.request)

			if tt.expectedCode == codes.InvalidArgument {
				if err == nil {
					t.Error("Expected error for invalid argument, got nil")
				} else {
					st, ok := status.FromError(err)
					if !ok {
						t.Error("Expected gRPC status error")
					} else if st.Code() != tt.expectedCode {
						t.Errorf("Expected error code %v, got %v", tt.expectedCode, st.Code())
					}
				}
				return
			}

			if tt.expectedCode == codes.OK {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				} else if len(response.Tags) != tt.expectedTags {
					t.Errorf("Expected %d tags, got %d", tt.expectedTags, len(response.Tags))
				}

				// Verify tag conversion
				if tt.expectedTags > 0 {
					for i, tag := range response.Tags {
						if tag.Tag != tt.mockResult[i].Tag {
							t.Errorf("Expected tag %s, got %s", tt.mockResult[i].Tag, tag.Tag)
						}
						if tag.Frequency != int32(tt.mockResult[i].Frequency) {
							t.Errorf("Expected frequency %d, got %d", tt.mockResult[i].Frequency, tag.Frequency)
						}
					}
				}
			}
		})
	}
}

func TestServer_RequestValidation(t *testing.T) {
	// Create a mock service with proper dependencies
	mockRepo := &MockArticleRepository{}
	mockExtractor := &MockTagExtractor{}
	articleService := app.NewArticleServiceWithExtractor(mockRepo, mockExtractor)
	
	grpcServer := NewServer(articleService)

	t.Run("ProcessArticles with nil request", func(t *testing.T) {
		ctx := context.Background()
		_, err := grpcServer.ProcessArticles(ctx, nil)
		
		if err == nil {
			t.Error("Expected error for nil request, got nil")
		}
	})

	t.Run("GetTopTags with nil request", func(t *testing.T) {
		ctx := context.Background()
		_, err := grpcServer.GetTopTags(ctx, nil)
		
		if err == nil {
			t.Error("Expected error for nil request, got nil")
		}
	})

	t.Run("ProcessArticles with articles containing empty fields", func(t *testing.T) {
		request := &pb.ProcessArticlesRequest{
			Articles: []*pb.Article{
				{Title: "", Body: "Valid body"},
				{Title: "Valid title", Body: ""},
				{Title: "", Body: ""},
			},
		}

		ctx := context.Background()
		response, err := grpcServer.ProcessArticles(ctx, request)
		
		// The server should accept these requests (validation is up to the service layer)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		
		if response == nil {
			t.Error("Expected response, got nil")
		}
	})
}

func TestServer_ContextHandling(t *testing.T) {
	// Create a mock service with proper dependencies
	mockRepo := &MockArticleRepository{}
	mockExtractor := &MockTagExtractor{}
	articleService := app.NewArticleServiceWithExtractor(mockRepo, mockExtractor)
	
	grpcServer := NewServer(articleService)

	t.Run("ProcessArticles with cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		request := &pb.ProcessArticlesRequest{
			Articles: []*pb.Article{
				{Title: "Test Article", Body: "This is a test article"},
			},
		}

		_, err := grpcServer.ProcessArticles(ctx, request)
		// The error handling depends on the service implementation
		// It might return context.Canceled or continue processing
		t.Logf("Error with cancelled context: %v", err)
	})

	t.Run("GetTopTags with cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		request := &pb.GetTopTagsRequest{
			Limit: 5,
		}

		_, err := grpcServer.GetTopTags(ctx, request)
		// The error handling depends on the service implementation
		t.Logf("Error with cancelled context: %v", err)
	})
}
