package repository

import (
	"context"

	"github.com/IrvanWijayaSardam/SelfBank/dto"
	"github.com/sashabaranov/go-openai"
)

type ChatbotRepository interface {
	Request(request dto.ChatRequest) (string, error)
}

type chatbotRepository struct {
	ai *openai.Client
}

func NewChatbotRepository(ai *openai.Client) ChatbotRepository {
	return &chatbotRepository{ai: ai}
}

func (repository *chatbotRepository) Request(request dto.ChatRequest) (string, error) {
	ctx := context.TODO()

	chatMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: request.Message,
	}

	chatRequest := openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{chatMessage},
	}

	resp, err := repository.ai.CreateChatCompletion(ctx, chatRequest)
	if err != nil {
		return "", err
	}

	reply := resp.Choices[0].Message.Content
	return reply, nil
}
