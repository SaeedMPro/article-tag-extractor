package app

import (
	"github.com/SaeedMPro/article-tag-extractor/utils"
)

type TagExtractorService struct {
}

func NewTagExtractorService() *TagExtractorService {
	return &TagExtractorService{}
}

func (t *TagExtractorService) ExtractTags(title, body string) []string {
	content := title + " " + body

	// simple tokenization by splitting on spaces and punctuation
	tokens := utils.Tokenize(content)

	// count word frequencies and filter stop words
	wordCount := make(map[string]int)
	for _, token := range tokens {
		if !utils.IsStopWord(token) && len(token) > 2 {
			wordCount[token]++
		}
	}

	// convert to slice and sort by frequency
	type wordFreq struct {
		word  string
		count int
	}
	var wordFreqs []wordFreq
	for word, count := range wordCount {
		wordFreqs = append(wordFreqs, wordFreq{word, count})
	}

	tags := []string{}
	for i := 0; i < 10 && i < len(wordFreqs); i++ {
		tags = append(tags, wordFreqs[i].word)
	}

	return tags
}
