package tokenizer

import (
	"regexp"
	"strings"
)

func removeSpecialChars(text string) string {
	// All whitespace characters, including newline, are converted to spaces.
	newline := regexp.MustCompile(`\s`)
	noNewLines := newline.ReplaceAllString(text, " ")

	// Special characters and numbers are removed from all words.
	special := regexp.MustCompile(`[^\p{L} ]`)
	noSpecial := special.ReplaceAllString(noNewLines, "")

	return noSpecial
}

func splitText(text string) []string {
	// Text is split up at whitespaces
	result := strings.Fields(text)
	return result
}
