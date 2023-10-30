package controller

import (
	"net/http"

	"github.com/IrvanWijayaSardam/SelfBank/dto"
	"github.com/IrvanWijayaSardam/SelfBank/helper"
	"github.com/IrvanWijayaSardam/SelfBank/service"
	"github.com/labstack/echo/v4"
)

type ChatbotController interface {
	Request(ctx echo.Context) error
}

type chatbotController struct {
	ChatbotService service.ChatbotService
	jwtService     service.JWTService
}

func NewChatbotController(fundService service.ChatbotService, jwtService service.JWTService) ChatbotController {
	return &chatbotController{
		ChatbotService: fundService,
		jwtService:     jwtService,
	}
}

func (c *chatbotController) Request(ctx echo.Context) error {
	var request dto.ChatRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("error when parsing data"))
	}

	result, err := c.ChatbotService.Request(request)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("error when parsing data"))
	}

	return ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Chatbot Replied", result))
}
