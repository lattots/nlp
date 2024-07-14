# NLP

## What is NLP?

NLP is a natural language processing library for Go. It is designed to make text processing easier for machine learning use cases.

## What can it do?

- Sentence tokenization:
  - NLP can perform sentence tokenization (i.e. split text into a slice of sentences) thanks to [this](https://github.com/neurosnap/sentences) excellent library by [neurosnap](https://github.com/neurosnap)
- Text tokenization:
  - NLP can normalize text and split it up to a slice of tokens / words
  - Example:
    - "I run a lot" -> {"i", "run", "a", "lot"}
- Lemmatization:
  - NLP can perform lemmatization (i.e. convert words to their base forms) thanks to [this](https://github.com/aaaton/golem) Go lemmatization library and [this](https://github.com/voikko/corevoikko) Finnish language library written in C

## How to use it?

### Installation

TODO: Write installation instructions

### Example use

```go
package main

import (
	"fmt"
	"log"

	"github.com/lattots/nlp/pkg/tokenizer"
)

func main() {
	language := "eng" // Target language for tokenizer
	tok, err := tokenizer.New(language)
	if err != nil {
		log.Fatalln("error creating tokenizer:", err)
	}

	exampleText := "This is some example text. It is used to showcase package tokenizer."

	tokens := tok.GetTokens(exampleText)
	fmt.Println(tokens)

    lemmas, err := tok.Lemmatize(tokens)
	if err != nil {
		log.Fatalln("error lemmatizing tokens:", err)
    }
	fmt.Println(lemmas)

	sentences := tok.GetSentences(exampleText)
	fmt.Println(sentences)
}
```
