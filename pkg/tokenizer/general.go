package tokenizer

import (
	"fmt"
	"os"
	"strings"

	"github.com/aaaton/golem/v4"
	"github.com/aaaton/golem/v4/dicts/de"
	"github.com/aaaton/golem/v4/dicts/en"
	"github.com/aaaton/golem/v4/dicts/es"
	"github.com/aaaton/golem/v4/dicts/fr"
	"github.com/aaaton/golem/v4/dicts/it"
	"github.com/aaaton/golem/v4/dicts/sv"
	"github.com/neurosnap/sentences"
)

// Tokenizer for English, Swedish, French, German, Spanish and Italian
type generalTokenizer struct {
	sentenceTokenizer *sentences.DefaultSentenceTokenizer
	lemmatizer        *golem.Lemmatizer
}

func (t *generalTokenizer) GetSentences(text string) []string {
	sens := t.sentenceTokenizer.Tokenize(text)
	result := make([]string, len(sens))
	for i, sen := range sens {
		result[i] = strings.TrimSpace(sen.Text)
	}
	return result
}

func (t *generalTokenizer) GetTokens(text string) []string {
	// Text is converted to lower case for easier use later
	text = strings.ToLower(text)
	// Special characters and numbers are removed from text
	text = removeSpecialChars(text)
	// Text is split up to individual tokens
	tokens := splitText(text)
	// Slice of tokens is returned
	return tokens
}

func (t *generalTokenizer) Lemmatize(tokens []string) ([]string, error) {
	lemmas := make([]string, len(tokens))
	for i, token := range tokens {
		lemma := t.lemmatizer.Lemma(token)
		if lemma == "" {
			lemmas[i] = token
		} else {
			lemmas[i] = lemma
		}
	}
	return lemmas, nil
}

func newGeneralTokenizer(language string) (*generalTokenizer, error) {
	const sentenceTokenizerPrefix = "../../data/sentence_tokenizer/"
	var sentenceTokenizerFilepath string
	var lemmatizer *golem.Lemmatizer
	switch language {
	case "eng":
		sentenceTokenizerFilepath = sentenceTokenizerPrefix + "english.json"
		var err error
		lemmatizer, err = golem.New(en.New())
		if err != nil {
			return nil, fmt.Errorf("error loading lemmatizer: %s\n", err)
		}
	case "swe":
		sentenceTokenizerFilepath = sentenceTokenizerPrefix + "swedish.json"
		var err error
		lemmatizer, err = golem.New(sv.New())
		if err != nil {
			return nil, fmt.Errorf("error loading lemmatizer: %s\n", err)
		}
	case "fra":
		sentenceTokenizerFilepath = sentenceTokenizerPrefix + "french.json"
		var err error
		lemmatizer, err = golem.New(fr.New())
		if err != nil {
			return nil, fmt.Errorf("error loading lemmatizer: %s\n", err)
		}
	case "deu":
		sentenceTokenizerFilepath = sentenceTokenizerPrefix + "german.json"
		var err error
		lemmatizer, err = golem.New(de.New())
		if err != nil {
			return nil, fmt.Errorf("error loading lemmatizer: %s\n", err)
		}
	case "ita":
		sentenceTokenizerFilepath = sentenceTokenizerPrefix + "italian.json"
		var err error
		lemmatizer, err = golem.New(it.New())
		if err != nil {
			return nil, fmt.Errorf("error loading lemmatizer: %s\n", err)
		}
	case "spa":
		sentenceTokenizerFilepath = sentenceTokenizerPrefix + "spanish.json"
		var err error
		lemmatizer, err = golem.New(es.New())
		if err != nil {
			return nil, fmt.Errorf("error loading lemmatizer: %s\n", err)
		}
	}

	// Sentence tokenizer training data is loaded from file.
	// Use data from https://github.com/neurosnap/sentences/data/
	data, err := os.ReadFile(sentenceTokenizerFilepath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %s\n", err)
	}
	// Model is trained with the given data.
	training, err := sentences.LoadTraining(data)
	if err != nil {
		return nil, fmt.Errorf("error loading training data: %s\n", err)
	}
	// Tokenizer instance is created with trained model.
	sentenceTokenizer := sentences.NewSentenceTokenizer(training)

	t := &generalTokenizer{
		sentenceTokenizer: sentenceTokenizer,
		lemmatizer:        lemmatizer,
	}
	return t, nil
}
