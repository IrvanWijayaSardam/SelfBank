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
	UpdateProfile(ctx echo.Context) error
	DeleteUser(ctx echo.Context) error
	All(ctx echo.Context) error
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

func (c *userController) All(ctx echo.Context) error {
	authHeader := ctx.Request().Header.Get("Authorization")
	pageParam := ctx.QueryParam("page")
	pageSizeParam := ctx.QueryParam("pageSize")

	defaultPage := 1
	defaultPageSize := 10

	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		page = defaultPage
	}

	pageSize, err := strconv.Atoi(pageSizeParam)
	if err != nil || pageSize < 1 {
		pageSize = defaultPageSize
	}

	token, err := c.jwtService.ValidateToken(authHeader)
	if err != nil {
		log.Println(err)
		response := helper.BuildErrorResponse("Token is not valid")
		return ctx.JSON(http.StatusUnauthorized, response)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ctx.Set("user", claims)

		roleID, ok := claims["idrole"].(float64)
		if !ok {
			response := helper.BuildErrorResponse("IDRole not found in claims")
			return ctx.JSON(http.StatusBadRequest, response)
		}
		switch roleID {
		case 1:
			response := helper.BuildErrorResponse("Unauthorized")
			return ctx.JSON(http.StatusUnauthorized, response)
		case 2:
			users, err := c.userService.All(page, pageSize)
			if err != nil {
				response := helper.BuildErrorResponse("Failed to fetch data")
				return ctx.JSON(http.StatusInternalServerError, response)
			}
			response := helper.BuildResponse(true, "OK!", users)
			return ctx.JSON(http.StatusOK, response)
		default:
			response := helper.BuildErrorResponse("Unauthorized")
			return ctx.JSON(http.StatusUnauthorized, response)
		}
	}
	response := helper.BuildErrorResponse("Invalid token claims")
	return ctx.JSON(http.StatusUnauthorized, response)
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

func (c *userController) UpdateProfile(context echo.Context) error {
	var updateUserDTO dto.UserUpdateDTO
	if err := context.Bind(&updateUserDTO); err != nil {
		response := helper.BuildErrorResponse("Failed to process request" + err.Error())
		return context.JSON(http.StatusBadRequest, response)
	}
	authHeader := context.Request().Header.Get("Authorization")
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
		user.Namadepan = updateUserDTO.Namadepan
		user.Namabelakang = updateUserDTO.Namabelakang
		user.Username = updateUserDTO.Username
		user.Telephone = updateUserDTO.Telephone
		user.Jk = updateUserDTO.Jk
		if updateUserDTO.Password != "" {
			user.Password = helper.HashAndSalt([]byte(updateUserDTO.Password))
		}

		c.userService.UpdateUser(user)
		response := helper.BuildResponse(true, "OK!", user)
		return context.JSON(http.StatusOK, response)
	} else {
		response := helper.BuildErrorResponse("Invalid token claims")
		return context.JSON(http.StatusUnauthorized, response)
	}
}

func (c *userController) DeleteUser(context echo.Context) error {
	authHeader := context.Request().Header.Get("Authorization")
	id := context.Param("id")

	token, err := c.jwtService.ValidateToken(authHeader)
	if err != nil {
		log.Println(err)
		response := helper.BuildErrorResponse("Token is not valid")
		return context.JSON(http.StatusUnauthorized, response)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		context.Set("user", claims)

		roleID, ok := claims["idrole"].(float64)
		if !ok {
			response := helper.BuildErrorResponse("IDRole not found in claims")
			return context.JSON(http.StatusBadRequest, response)
		}
		switch roleID {
		case 1:
			response := helper.BuildErrorResponse("Unauthorized")
			return context.JSON(http.StatusUnauthorized, response)
		case 2:
			res := c.userService.DeleteUser(helper.StringToUint64(id))
			if res {
				response := helper.BuildOkResponse(res, "Users Succesfully Deleted !"+id)
				return context.JSON(http.StatusOK, response)
			} else {
				response := helper.BuildErrorResponse("Failed to update user")
				return context.JSON(http.StatusBadRequest, response)
			}

		default:
			response := helper.BuildErrorResponse("Unauthorized")
			return context.JSON(http.StatusUnauthorized, response)
		}
	}
	response := helper.BuildErrorResponse("Invalid token claims")
	return context.JSON(http.StatusUnauthorized, response)

}

func (c *userController) FileUpload(context echo.Context) error {

	authHeader := context.Request().Header.Get("Authorization")

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
