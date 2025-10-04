package utils

import (
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Simple text",
			input:    "Hello World",
			expected: []string{"hello", "world"},
		},
		{
			name:     "Text with punctuation",
			input:    "Hello, World! How are you?",
			expected: []string{"hello", "world", "how", "are", "you"},
		},
		{
			name:     "Text with numbers and symbols",
			input:    "Go 1.21 is awesome! @golang #programming",
			expected: []string{"go", "is", "awesome", "golang", "programming"},
		},
		{
			name:     "Empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "Only punctuation",
			input:    "!@#$%^&*()",
			expected: []string{},
		},
		{
			name:     "Mixed case",
			input:    "Go Programming Language",
			expected: []string{"go", "programming", "language"},
		},
		{
			name:     "Text with special characters",
			input:    "API v2.0 & REST endpoints",
			expected: []string{"api", "v", "rest", "endpoints"},
		},
		{
			name:     "Text with multiple spaces",
			input:    "Multiple    spaces   between   words",
			expected: []string{"multiple", "spaces", "between", "words"},
		},
		{
			name:     "Text with tabs and newlines",
			input:    "Text\twith\nnewlines",
			expected: []string{"text", "with", "newlines"},
		},
		{
			name:     "Unicode text",
			input:    "Café naïve résumé",
			expected: []string{"caf", "na", "ve", "r", "sum"}, // Unicode characters are stripped by regex
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Tokenize(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d tokens, got %d", len(tt.expected), len(result))
				t.Errorf("Expected: %v", tt.expected)
				t.Errorf("Got: %v", result)
				return
			}

			for i, token := range result {
				if token != tt.expected[i] {
					t.Errorf("Expected token %d to be '%s', got '%s'", i, tt.expected[i], token)
				}
			}
		})
	}
}

func TestTokenizeAndStopWordIntegration(t *testing.T) {
	tests := []struct {
		name                 string
		input                string
		expectedTokens       int
		expectedNonStopWords int
	}{
		{
			name:                 "Text with many stopwords",
			input:                "The quick brown fox jumps over the lazy dog",
			expectedTokens:       9,
			expectedNonStopWords: 6, // quick, brown, fox, jumps, lazy, dog
		},
		{
			name:                 "Technical text",
			input:                "Go is a programming language developed by Google",
			expectedTokens:       8,
			expectedNonStopWords: 5, // go, programming, language, developed, google
		},
		{
			name:                 "Mixed content",
			input:                "The API provides REST endpoints for data access",
			expectedTokens:       8,
			expectedNonStopWords: 6, // api, provides, rest, endpoints, data, access
		},
		{
			name:                 "Only stopwords",
			input:                "The and or but not",
			expectedTokens:       5,
			expectedNonStopWords: 0,
		},
		{
			name:                 "No stopwords",
			input:                "Programming language development",
			expectedTokens:       3,
			expectedNonStopWords: 3, // programming, language, development
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := Tokenize(tt.input)

			if len(tokens) != tt.expectedTokens {
				t.Errorf("Expected %d tokens, got %d", tt.expectedTokens, len(tokens))
			}

			nonStopWords := 0
			var nonStopWordsList []string
			for _, token := range tokens {
				if !IsStopWord(token) {
					nonStopWords++
					nonStopWordsList = append(nonStopWordsList, token)
				}
			}

			if nonStopWords != tt.expectedNonStopWords {
				t.Errorf("Expected %d non-stop words, got %d. Words: %v", tt.expectedNonStopWords, nonStopWords, nonStopWordsList)
			}
		})
	}
}
