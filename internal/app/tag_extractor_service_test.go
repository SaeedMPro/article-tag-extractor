package app

import (
	"strings"
	"testing"

	"github.com/SaeedMPro/article-tag-extractor/utils"
)

func TestTagExtractorService_Combined(t *testing.T) {
	extractor := NewTagExtractorService()

	tests := []struct {
		name     string
		title    string
		body     string
		expected int // minimum expected tags
	}{
		{
			name:     "Simple article with programming content",
			title:    "Go Programming Language",
			body:     "Go is a programming language developed by Google. Go is fast and efficient for building scalable applications.",
			expected: 4,
		},
		{
			name:     "Article with many stop words",
			title:    "The Art of Programming",
			body:     "The art of programming is a skill that requires practice and dedication. Programming is fun and rewarding for developers.",
			expected: 5,
		},
		{
			name:     "Empty content",
			title:    "",
			body:     "",
			expected: 0,
		},
		{
			name:     "Only stop words",
			title:    "The and or but",
			body:     "This is a test with only common words that are stopwords",
			expected: 2,
		},
		{
			name:     "Technical article",
			title:    "Microservices Architecture",
			body:     "Microservices architecture enables building scalable distributed systems using containerization and orchestration tools like Kubernetes.",
			expected: 6,
		},
		{
			name:     "Article with punctuation and special characters",
			title:    "API Development @2024",
			body:     "REST APIs are essential for modern web development. JSON, HTTP, and authentication are key concepts!",
			expected: 6,
		},
		{
			name:     "Single word title",
			title:    "Programming",
			body:     "This is about programming",
			expected: 2,
		},
		{
			name:     "Numbers and symbols",
			title:    "Go 1.21 Features",
			body:     "Go 1.21 introduced new features including generics and improved error handling.",
			expected: 6,
		},
		{
			name:     "Very long content",
			title:    "Long Article",
			body:     "This is a very long article with many words repeated multiple times to test the tag extraction algorithm with extensive content that should produce meaningful tags from the most frequent words in the text content.",
			expected: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tags := extractor.ExtractTags(tt.title, tt.body)

			// Check minimum expected tags
			if len(tags) < tt.expected {
				t.Errorf("Expected at least %d tags, got %d", tt.expected, len(tags))
			}

			// Validate tag content
			tagMap := make(map[string]bool)
			for i, tag := range tags {
				// No stop words
				if utils.IsStopWord(tag) {
					t.Errorf("Tag '%s' is a stop word", tag)
				}
				// Not empty
				if tag == "" {
					t.Errorf("Tag %d is empty", i)
				}
				// Lowercase
				if tag != strings.ToLower(tag) {
					t.Errorf("Tag '%s' should be lowercase", tag)
				}
				// No duplicates
				if tagMap[tag] {
					t.Errorf("Duplicate tag found: %s", tag)
				}
				tagMap[tag] = true
			}
		})
	}
}
