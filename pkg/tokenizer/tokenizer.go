package tokenizer

import (
	"errors"
)

type Tokenizer interface {
	GetTokens(text string) []string
	GetSentences(text string) []string
	Lemmatize(tokens []string) ([]string, error)
}

func New(language string) (Tokenizer, error) {
	var tokenizer Tokenizer
	switch language {
	case "fin":
		var err error
		tokenizer, err = newFinnishTokenizer()
		if err != nil {
			return nil, err
		}
		return tokenizer, nil
	case "eng", "swe", "fra", "deu", "ita", "spa":
		var err error
		tokenizer, err = newGeneralTokenizer(language)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("invalid language")
	}
	return tokenizer, nil
}
