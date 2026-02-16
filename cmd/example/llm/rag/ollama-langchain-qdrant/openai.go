package main

import (
	"github.com/tmc/langchaingo/llms/openai"
)

func NewOpenAI() (*openai.LLM, error) {
	opts := []openai.Option{
		openai.WithToken("YOUR-OPENAI-TOKEN"),
		openai.WithModel("gpt-4o-mini"),
		openai.WithEmbeddingModel("text-embedding-3-small"),
	}

	return openai.New(opts...)
}
