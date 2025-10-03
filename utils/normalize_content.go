package utils

import (
	"regexp"
	"strings"
)

var stopWords = map[string]bool{
	"a": true, "an": true, "and": true, "are": true, "as": true, "at": true,
	"be": true, "by": true, "for": true, "from": true, "he": true,
	"in": true, "is": true, "it": true, "its": true, "of": true, "on": true,
	"that": true, "the": true, "to": true, "was": true, "will": true, "with": true,
	"this": true, "these": true, "they": true, "them": true, "their": true,
	"there": true, "then": true, "than": true, "or": true, "but": true, "not": true,
	"have": true, "had": true, "has": true, "having": true, "been": true,
	"being": true, "do": true, "does": true, "did": true, "doing": true,
	"can": true, "could": true, "should": true, "would": true, "may": true,
	"might": true, "must": true, "shall": true, "i": true, "you": true, "we": true,
	"us": true, "our": true, "my": true, "me": true, "him": true, "her": true,
	"his": true, "she": true, "all": true, "any": true, "both": true, "each": true,
	"few": true, "more": true, "most": true, "other": true, "some": true,
	"such": true, "no": true, "nor": true, "so": true, "too": true, "very": true,
}

func IsStopWord(word string) bool {
	return stopWords[word]
}

// Tokenize splits the input text into lowercase words, removing punctuation.
func Tokenize(text string) []string {
	words := []string{}

	reg := regexp.MustCompile(`[^a-zA-Z\s]`)
	content := reg.ReplaceAllString(text, " ")

	for _, word := range strings.Fields(content) {
		words = append(words, strings.ToLower(word))
	}
	return words
}
