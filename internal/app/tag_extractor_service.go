package app

import (
	"sync"

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
	
	// sort by frequency (descending) and then by word (ascending) for consistency
	for i := 0; i < len(wordFreqs)-1; i++ {
		for j := i + 1; j < len(wordFreqs); j++ {
			if wordFreqs[i].count < wordFreqs[j].count || 
			   (wordFreqs[i].count == wordFreqs[j].count && wordFreqs[i].word > wordFreqs[j].word) {
				wordFreqs[i], wordFreqs[j] = wordFreqs[j], wordFreqs[i]
			}
		}
	}
	
	// extract top 10 tags
	maxTags := 10
	if len(wordFreqs) < maxTags {
		maxTags = len(wordFreqs)
	}
	
	var tags []string
	for i := 0; i < maxTags; i++ {
		tags = append(tags, wordFreqs[i].word)
	}
	
	return tags
}

func (t *TagExtractorService) ExtractTagsConcurrently(title, body string) []string {
	content := title + " " + body
	tokens := utils.Tokenize(content)

	// count word frequencies and filter stop words
	wordCount := make(map[string]int)
	var mu sync.Mutex
	
	var wg sync.WaitGroup
	for _, token := range tokens {
		wg.Add(1)
		go func(t string) {
			defer wg.Done()
			if !utils.IsStopWord(t) && len(t) > 2 {
				mu.Lock()
				wordCount[t]++
				mu.Unlock()
			}
		}(token)
	}

	wg.Wait()

	// convert to slice and sort by frequency
	type wordFreq struct {
		word  string
		count int
	}
	
	var wordFreqs []wordFreq
	for word, count := range wordCount {
		wordFreqs = append(wordFreqs, wordFreq{word, count})
	}
	
	// sort by frequency (descending) and then by word (ascending) for consistency
	for i := 0; i < len(wordFreqs)-1; i++ {
		for j := i + 1; j < len(wordFreqs); j++ {
			if wordFreqs[i].count < wordFreqs[j].count || 
			   (wordFreqs[i].count == wordFreqs[j].count && wordFreqs[i].word > wordFreqs[j].word) {
				wordFreqs[i], wordFreqs[j] = wordFreqs[j], wordFreqs[i]
			}
		}
	}
	
	// extract top 10 tags
	maxTags := 10
	if len(wordFreqs) < maxTags {
		maxTags = len(wordFreqs)
	}
	
	var tags []string
	for i := 0; i < maxTags; i++ {
		tags = append(tags, wordFreqs[i].word)
	}
	
	return tags
}