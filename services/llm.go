package services

import (
	"context"
	"fmt"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/llms/openai"
)

type Provider string

const (
	ProviderOpenAI Provider = "openai"
	ProviderGemini Provider = "gemini"
)

type LLMService struct {
	llm    llms.Model
	prompt string
}

func NewLLMService(provider Provider, apiKey, prompt string) (*LLMService, error) {
	var llmInstance llms.Model
	var err error

	switch provider {
	case ProviderOpenAI:
		llmInstance, err = openai.New(openai.WithToken(apiKey),
			openai.WithModel("gpt-4o"))
		if err != nil {
			return nil, fmt.Errorf("failed to initialize OpenAI: %w", err)
		}

	case ProviderGemini:
		llmInstance, err = googleai.New(context.Background(), googleai.WithAPIKey(apiKey),
			googleai.WithDefaultModel("gemini-1.5-pro"))
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Gemini: %w", err)
		}

	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	return &LLMService{
		llm:    llmInstance,
		prompt: prompt,
	}, nil
}

func (l *LLMService) Analyze(imagePath string) (string, error) {
	imgContent, err := l.readImageContent(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to read image: %w", err)
	}

	ctx := context.Background()
	parts := []llms.ContentPart{

		llms.TextPart(l.prompt),
		llms.BinaryPart("image/png", imgContent),
	}

	completion, err := l.llm.GenerateContent(ctx, []llms.MessageContent{
		{
			Parts: parts,
			Role:  llms.ChatMessageTypeHuman,
		},
	})

	if err != nil {
		return "", fmt.Errorf("error getting response from LLM: %v", err)
	}

	result := completion.Choices[0].Content
	fmt.Println(result)
	return result, nil
}

func (l *LLMService) readImageContent(imagePath string) ([]byte, error) {
	imgBytes, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, err
	}

	return imgBytes, nil
}
