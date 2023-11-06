package controller

import (
	"net/http"

	"github.com/IrvanWijayaSardam/SelfBank/dto"
	"github.com/IrvanWijayaSardam/SelfBank/helper"
	"github.com/IrvanWijayaSardam/SelfBank/service"
	"github.com/labstack/echo/v4"
)

type VerificationController interface {
	SendVerificationEmail(ctx echo.Context) error
	ValidateVerification(ctx echo.Context) error
}

type verificationController struct {
	verificationService service.VerificationService
	jwtService          service.JWTService
}

func NewVerificationController(verifService service.VerificationService, jwtService service.JWTService) VerificationController {
	return &verificationController{
		verificationService: verifService,
		jwtService:          jwtService,
	}
}

func (c *verificationController) SendVerificationEmail(ctx echo.Context) error {
	var verificationDTO dto.VerificationDTO
	if err := ctx.Bind(&verificationDTO); err != nil {
		response := helper.BuildErrorResponse("Failed to process request")
		return ctx.JSON(http.StatusBadRequest, response)
	}

	sendEmail := c.verificationService.SendVerificationEmail(verificationDTO.Email)
	if sendEmail != nil {
		response := helper.BuildErrorResponse("Failed Sending Email Verification" + sendEmail.Error())
		return ctx.JSON(http.StatusBadGateway, response)
	}

	response := helper.BuildOkResponse(true, "Please Check Your Email To Do A Verification")
	return ctx.JSON(http.StatusOK, response)
}

func (c *verificationController) ValidateVerification(ctx echo.Context) error {
	otp := ctx.QueryParam("otp")

	verifyOTP := c.verificationService.VerifyOtp(otp)
	if verifyOTP != true {
		response := helper.BuildErrorResponse("Incorrect / Expired OTP")
		return ctx.JSON(http.StatusBadGateway, response)
	}

	response := helper.BuildOkResponse(true, "Account verified")
	return ctx.JSON(http.StatusOK, response)
}
