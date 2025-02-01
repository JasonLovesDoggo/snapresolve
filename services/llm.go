package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

type LLMService struct {
	client *openai.Client
}

func NewLLMService() *LLMService {
	return &LLMService{}
}

func (l *LLMService) Init(apiKey string) {
	l.client = openai.NewClient(apiKey)
}

func (l *LLMService) Analyze(imagePath, prompt string) (string, error) {
	if l.client == nil {
		return "", fmt.Errorf("OpenAI client not initialized")
	}

	imgBytes, err := os.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to read image: %w", err)
	}

	b64Img := base64.StdEncoding.EncodeToString(imgBytes)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := l.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4VisionPreview,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleUser,
					MultiContent: []openai.ChatMessagePart{
						{
							Type: openai.ChatMessagePartTypeText,
							Text: prompt,
						},
						{
							Type: openai.ChatMessagePartTypeImageURL,
							ImageURL: &openai.ChatMessageImageURL{
								URL:    fmt.Sprintf("data:image/png;base64,%s", b64Img),
								Detail: openai.ImageURLDetailHigh,
							},
						},
					},
				},
			},
			MaxTokens: 500,
		},
	)

	if err != nil {
		return "", fmt.Errorf("OpenAI API error: %w", err)
	}

	return resp.Choices[0].Message.Content, nil
}
