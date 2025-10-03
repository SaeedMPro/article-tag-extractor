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

	// remove stop words
	tags := make([]string, 0)
	for _, token := range tokens {
		if !utils.IsStopWord(token) {
			tags = append(tags, token)
		}
	}
	return tags
}

func (t *TagExtractorService) ExtractTagsConcurrently(title, body string) []string {
	content := title + " " + body
	tokens := utils.Tokenize(content)

	var wg sync.WaitGroup
	tagsChan := make(chan string, len(tokens))

	for _, token := range tokens {
		wg.Add(1)
		go func(t string) {
			defer wg.Done()
			if !utils.IsStopWord(t) {
				tagsChan <- t
			}
		}(token)
	}

	wg.Wait()
	close(tagsChan)

	tags := make([]string, 0, len(tagsChan))
	for tag := range tagsChan {
		tags = append(tags, tag)
	}
	return tags
}