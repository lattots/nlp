package tokenizer

import (
	"fmt"
	"os"
	"strings"

	"github.com/neurosnap/sentences"

	"github.com/lattots/nlp/pkg/fin_lemmatizer"
)

type finnishTokenizer struct {
	sentenceTokenizer *sentences.DefaultSentenceTokenizer
	voikkoHandle      fin_lemmatizer.VoikkoHandle
}

func (t *finnishTokenizer) GetSentences(text string) []string {
	sens := t.sentenceTokenizer.Tokenize(text)
	result := make([]string, len(sens))
	for i, sen := range sens {
		result[i] = strings.TrimSpace(sen.Text)
	}
	return result
}

func (t *finnishTokenizer) GetTokens(text string) []string {
	// Text is converted to lower case for easier use later
	text = strings.ToLower(text)
	// Special characters and numbers are removed from text
	text = removeSpecialChars(text)
	// Text is split up to individual tokens
	tokens := splitText(text)
	// Slice of tokens is returned
	return tokens
}

func (t *finnishTokenizer) Lemmatize(tokens []string) ([]string, error) {
	lemmas, err := fin_lemmatizer.Batch(t.voikkoHandle, tokens)
	if err != nil {
		return nil, fmt.Errorf("error lemmatizing tokens: %s\n", err)
	}
	return lemmas, nil
}

func newFinnishTokenizer() (*finnishTokenizer, error) {
	voikkoHandle, err := fin_lemmatizer.InitVoikko()
	if err != nil {
		return nil, err
	}

	// Finnish sentence tokenizer training data is loaded from file.
	// Use data from https://github.com/neurosnap/sentences/data/
	const filepath = "../../data/sentence_tokenizer/finnish.json"
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	// Model is trained with the given data.
	training, err := sentences.LoadTraining(data)
	if err != nil {
		return nil, err
	}
	// Tokenizer instance is created with trained model.
	sentenceTokenizer := sentences.NewSentenceTokenizer(training)

	t := finnishTokenizer{
		sentenceTokenizer: sentenceTokenizer,
		voikkoHandle:      voikkoHandle,
	}
	return &t, nil
}
