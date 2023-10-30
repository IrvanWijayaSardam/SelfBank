package service

import (
	"github.com/IrvanWijayaSardam/SelfBank/dto"
	"github.com/IrvanWijayaSardam/SelfBank/repository"
)

type ChatbotService interface {
	Request(request dto.ChatRequest) (string, error)
}

type chatbotService struct {
	ChatbotRepository repository.ChatbotRepository
}

func NewChatbotService(fundRep repository.ChatbotRepository) ChatbotService {
	return &chatbotService{
		ChatbotRepository: fundRep,
	}
}

func (service *chatbotService) Request(request dto.ChatRequest) (string, error) {
	return service.ChatbotRepository.Request(request)
}
