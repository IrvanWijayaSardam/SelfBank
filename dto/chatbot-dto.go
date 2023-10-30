package dto

type ChatRequest struct {
	Message  string `json:"message" form:"message" validate:"required"`
	Response string
}
