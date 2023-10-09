package controller

import (
	"net/http"
	"strconv"

	"github.com/IrvanWijayaSardam/SelfBank/dto"
	"github.com/IrvanWijayaSardam/SelfBank/entity"
	"github.com/IrvanWijayaSardam/SelfBank/helper"
	"github.com/IrvanWijayaSardam/SelfBank/service"

	"github.com/labstack/echo/v4"
)

type AuthController interface {
	Login(ctx echo.Context) error
	Register(ctx echo.Context) error
}

type authController struct {
	authService service.AuthService
	jwtService  service.JWTService
}

func NewAuthController(authService service.AuthService, jwtService service.JWTService) AuthController {
	return &authController{
		authService: authService,
		jwtService:  jwtService,
	}
}

func (c *authController) Login(ctx echo.Context) error {
	var loginDTO dto.LoginDTO
	if err := ctx.Bind(&loginDTO); err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		return ctx.JSON(http.StatusBadRequest, response)
	}

	authResult := c.authService.VerifyCredential(loginDTO.Email, loginDTO.Password)
	if v, ok := authResult.(entity.User); ok {
		generatedToken, _ := c.jwtService.GenerateToken(strconv.FormatUint(v.ID, 10), v.Namadepan, v.Email, v.Telephone, v.Jk, v.IdRole)
		v.Token = generatedToken
		response := helper.BuildResponse(true, "OK", v)
		return ctx.JSON(http.StatusOK, response)
	}

	response := helper.BuildErrorResponse("Please check your credentials", "Invalid Credential", helper.EmptyObj{})
	return ctx.JSON(http.StatusUnauthorized, response)
}

func (c *authController) Register(ctx echo.Context) error {
	var registerDTO dto.RegisterDTO
	if err := ctx.Bind(&registerDTO); err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		return ctx.JSON(http.StatusBadRequest, response)
	}

	if !c.authService.IsDuplicateEmail(registerDTO.Email) {
		response := helper.BuildErrorResponse("Failed to process request", "Duplicate email", helper.EmptyObj{})
		return ctx.JSON(http.StatusConflict, response)
	}

	createdUser := c.authService.CreateUser(registerDTO)
	token, _ := c.jwtService.GenerateToken(strconv.FormatUint(createdUser.ID, 10), createdUser.Namadepan, createdUser.Email, createdUser.Telephone, createdUser.Jk, createdUser.IdRole)
	createdUser.Token = token
	response := helper.BuildResponse(true, "OK!", createdUser)
	return ctx.JSON(http.StatusCreated, response)
}
