package main

import (
	"github.com/tmc/langchaingo/llms/ollama"
)

func NewOllama() (*ollama.LLM, error) {
	opts := []ollama.Option{
		ollama.WithModel(*modelFlag),
		ollama.WithServerURL(*ollamaUrlFlag),
	}

	return ollama.New(opts...)
}
