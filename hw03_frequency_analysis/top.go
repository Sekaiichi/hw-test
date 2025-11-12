package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var splitPattern = regexp.MustCompile(`[\s!.,?]+`)

type frequencyWords struct {
	word      string
	frequency int
}

func Top10(str string) []string {
	// Place your code here
	countWords := wordFrequencyAnalysis(str)
	sort.Slice(countWords, func(i, j int) bool {
		return countWords[i].frequency > countWords[j].frequency
	})

	return topResult(countWords, 10)
}

func wordFrequencyAnalysis(str string) []frequencyWords {
	words := wordsCounter(str)
	wordsFreq := make([]frequencyWords, 0, len(words))
	for word, count := range words {
		wordsFreq = append(wordsFreq, frequencyWords{word: word, frequency: count})
	}

	return wordsFreq
}

func wordsCounter(str string) map[string]int {
	wordCount := make(map[string]int)
	for _, word := range splitPattern.Split(str, -1) {
		if word == "" || word == "-" {
			continue
		}
		word = strings.ToLower(word)
		wordCount[word]++
	}

	return wordCount
}

func topResult(words []frequencyWords, n int) []string {
	if n > len(words) {
		n = len(words)
	}
	result := make([]string, 0, n)
	for _, word := range words[:n] {
		result = append(result, word.word)
	}
	return result
}
