package dto

type VerificationDTO struct {
	Email string `json:"email" form:"email" binding:"required,email"`
}
