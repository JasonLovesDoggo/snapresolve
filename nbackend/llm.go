package backend

import (
	"context"
	"encoding/base64"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

type LLMService struct {
	client *openai.Client
}

func NewLLMService() *LLMService {
	return &LLMService{}
}

func (l *LLMService) Initialize(apiKey string) {
	l.client = openai.NewClient(apiKey)
}

func (l *LLMService) AnalyzeScreenshot(imagePath string, prompt string) (string, error) {
	imageData, _ := os.ReadFile(imagePath)
	b64 := base64.StdEncoding.EncodeToString(imageData)

	resp, err := l.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4VisionPreview,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleUser,
					MultiContent: []openai.ChatMessagePart{
						{Type: openai.ChatMessagePartTypeText, Text: prompt},
						{Type: openai.ChatMessagePartTypeImageURL, ImageURL: &openai.ChatMessageImageURL{
							URL: "data:image/png;base64," + b64,
						}},
					},
				},
			},
		},
	)

	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}
