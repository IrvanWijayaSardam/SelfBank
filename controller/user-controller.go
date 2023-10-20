package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/IrvanWijayaSardam/SelfBank/dto"
	"github.com/IrvanWijayaSardam/SelfBank/helper"
	"github.com/IrvanWijayaSardam/SelfBank/service"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type UserController interface {
	MyProfile(ctx echo.Context) error
	FileUpload(ctx echo.Context) error
}

type userController struct {
	userService service.UserService
	jwtService  service.JWTService
}

func NewUserController(userService service.UserService, jwtService service.JWTService) UserController {
	return &userController{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (c *userController) MyProfile(context echo.Context) error {
	authHeader := context.Request().Header.Get("Authorization")
	errEnv := godotenv.Load()
	if errEnv != nil {
		panic("Failed to load env file")
	}

	token, err := c.jwtService.ValidateToken(authHeader)
	if err != nil {
		log.Println(err)
		response := helper.BuildErrorResponse("Token is not valid")
		return context.JSON(http.StatusUnauthorized, response)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		context.Set("user", claims)

		userIDStr, ok := claims["userid"].(string)
		if !ok {
			response := helper.BuildErrorResponse("User ID not found in claims")
			return context.JSON(http.StatusBadRequest, response)
		}

		userID, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			response := helper.BuildErrorResponse("Failed to convert User ID to uint64")
			return context.JSON(http.StatusBadRequest, response)
		}

		user := c.userService.FindUser(userID)
		user.Balance = strconv.FormatInt(c.userService.GetSaldo(userID), 10)

		response := helper.BuildResponse(true, "OK!", user)
		return context.JSON(http.StatusOK, response)
	}

	response := helper.BuildErrorResponse("Invalid token claims")
	return context.JSON(http.StatusUnauthorized, response)
}

func (c *userController) FileUpload(context echo.Context) error {
	authHeader := context.Request().Header.Get("Authorization")
	errEnv := godotenv.Load()
	if errEnv != nil {
		panic("Failed to load env file")
	}

	token, err := c.jwtService.ValidateToken(authHeader)
	if err != nil {
		log.Println(err)
		response := helper.BuildErrorResponse("Token is not valid")
		return context.JSON(http.StatusUnauthorized, response)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		context.Set("user", claims)

		userIDStr, ok := claims["userid"].(string)
		if !ok {
			response := helper.BuildErrorResponse("User ID not found in claims")
			return context.JSON(http.StatusBadRequest, response)
		}

		userID, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			response := helper.BuildErrorResponse("Failed to convert User ID to uint64")
			return context.JSON(http.StatusBadRequest, response)
		}

		user := c.userService.FindUser(userID)

		formfile, err := context.FormFile("file")
		if err != nil {
			return context.JSON(
				http.StatusInternalServerError,
				dto.MediaDto{
					StatusCode: http.StatusInternalServerError,
					Message:    "error",
					Data:       map[string]interface{}{"data": "Select a file to upload"},
				})
		}

		// Open the uploaded file
		file, err := formfile.Open()
		if err != nil {
			return context.JSON(
				http.StatusInternalServerError,
				dto.MediaDto{
					StatusCode: http.StatusInternalServerError,
					Message:    "error",
					Data:       map[string]interface{}{"data": "Error opening uploaded file"},
				})
		}
		defer file.Close()

		// Pass the file to the service
		uploadUrl, err := service.NewMediaUpload().FileUpload(dto.File{File: file})
		if err != nil {
			response := helper.BuildErrorResponse("Failed to upload file")
			return context.JSON(http.StatusBadRequest, response)
		}
		user.Profile = uploadUrl
		c.userService.UpdateUser(user)

		response := helper.BuildResponse(true, "Image Successfully Uploaded", user)
		return context.JSON(http.StatusOK, response)
	}
	response := helper.BuildErrorResponse("Invalid token claims")
	return context.JSON(http.StatusUnauthorized, response)

}
