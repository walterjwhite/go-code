package ai

import (
	"context"
	"fmt"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

func init() {
	_ = llms.GenerateFromSinglePrompt
}

type AIClient struct {
	llm llms.Model
}

func NewAIClient() (*AIClient, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	llm, err := openai.New(
		openai.WithModel("gpt-3.5-turbo"),
		openai.WithToken(apiKey),
	)
	if err != nil {
		return nil, err
	}

	client := &AIClient{
		llm: llm,
	}

	return client, nil
}

func (client *AIClient) GenerateAnswer(ctx context.Context, question string) (string, error) {
	prompt := fmt.Sprintf("Please provide a professional and concise answer to the following job application question: %s", question)

	completion, err := llms.GenerateFromSinglePrompt(ctx, client.llm, prompt, llms.WithTemperature(0.7))
	if err != nil {
		return "", err
	}

	return completion, nil
}

func (client *AIClient) GenerateAnswerWithConstraints(ctx context.Context, question string, constraints string) (string, error) {
	prompt := fmt.Sprintf("Please provide a professional and concise answer to the following job application question: %s. %s", question, constraints)

	completion, err := llms.GenerateFromSinglePrompt(ctx, client.llm, prompt, llms.WithTemperature(0.7))
	if err != nil {
		return "", err
	}

	return completion, nil
}
