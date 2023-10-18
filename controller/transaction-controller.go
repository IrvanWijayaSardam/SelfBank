package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"

	"github.com/IrvanWijayaSardam/SelfBank/dto"
	"github.com/IrvanWijayaSardam/SelfBank/entity"
	"github.com/IrvanWijayaSardam/SelfBank/helper"
	"github.com/IrvanWijayaSardam/SelfBank/service"
)

type TransactionController interface {
	Insert(context echo.Context) error
	All(context echo.Context) error
	FindTransactionByID(context echo.Context) error
}

type transactionController struct {
	TransactionService service.TransactionService
	jwtService         service.JWTService
}

func NewTransactionController(transactionService service.TransactionService, jwtService service.JWTService) TransactionController {
	return &transactionController{
		TransactionService: transactionService,
		jwtService:         jwtService,
	}
}

func (c *transactionController) Insert(context echo.Context) error {
	authHeader := context.Request().Header.Get("Authorization")
	token, err := c.jwtService.ValidateToken(authHeader)

	if err != nil {
		log.Println(err)
		response := helper.BuildErrorResponse("Token is not valid")
		return context.JSON(http.StatusUnauthorized, response)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		context.Set("user", claims)

		userID, ok := claims["userid"].(string)
		if !ok {
			response := helper.BuildErrorResponse("UserID not found in claims")
			return context.JSON(http.StatusBadRequest, response)
		}

		accountNumber, ok := claims["accountnumber"].(string)
		if !ok {
			response := helper.BuildErrorResponse("Account Number not found in claims")
			return context.JSON(http.StatusBadRequest, response)
		}
		var TransactionDTO dto.TransactionDTO
		if err := context.Bind(&TransactionDTO); err != nil {
			response := helper.BuildErrorResponse("Failed to process request " + err.Error())
			return context.JSON(http.StatusBadRequest, response)
		}

		validateTo := c.TransactionService.ValidateAccNumber(TransactionDTO.TransactionTo)
		if validateTo == false {
			response := helper.BuildErrorResponse("Nomor Rekening Tujuan Tidak Valid")
			return context.JSON(http.StatusBadRequest, response)
		}

		TransactionDTO.ID_User, _ = strconv.ParseUint(userID, 10, 64)
		TransactionDTO.TransactionFrom, _ = strconv.ParseUint(accountNumber, 10, 64)
		Transaction := c.TransactionService.InsertTransaction(TransactionDTO)

		res := helper.BuildResponse(true, "Transaction Success", Transaction)
		return context.JSON(http.StatusCreated, res)
	} else {
		res := helper.BuildErrorResponse("There's Something Wrong")
		return context.JSON(http.StatusBadRequest, res)
	}
}

func (c *transactionController) All(context echo.Context) error {
	pageParam := context.QueryParam("page")
	pageSizeParam := context.QueryParam("pageSize")

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
	authHeader := context.Request().Header.Get("Authorization")

	token, err := c.jwtService.ValidateToken(authHeader)
	if err != nil {
		log.Println(err)
		response := helper.BuildErrorResponse("Token is not valid")
		return context.JSON(http.StatusUnauthorized, response)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		context.Set("user", claims)

		userID, ok := claims["userid"].(string)
		if !ok {
			response := helper.BuildErrorResponse("UserID not found in claims")
			return context.JSON(http.StatusBadRequest, response)
		}

		roleID, ok := claims["idrole"].(float64)
		if !ok {
			response := helper.BuildErrorResponse("IDRole not found in claims")
			return context.JSON(http.StatusBadRequest, response)
		}
		switch roleID {
		case 1:
			Transactions, err := c.TransactionService.All(page, pageSize)
			if err != nil {
				response := helper.BuildErrorResponse("Failed to fetch data")
				return context.JSON(http.StatusInternalServerError, response)
			}

			total := c.TransactionService.TotalTransaction()

			totalPages := (int(total) + pageSize - 1) / pageSize

			customResponse := struct {
				Status  bool                      `json:"status"`
				Message string                    `json:"message"`
				Errors  interface{}               `json:"errors"`
				Data    []entity.Transaction      `json:"data"`
				Paging  helper.PaginationResponse `json:"paging"`
			}{
				Status:  true,
				Message: "OK!",
				Errors:  nil,
				Data:    Transactions,
				Paging:  helper.PaginationResponse{TotalRecords: int(total), CurrentPage: page, TotalPages: totalPages},
			}

			return context.JSON(http.StatusOK, customResponse)
		case 2:
			userIDCnv, err := strconv.ParseUint(userID, 10, 64)
			if err != nil {
				fmt.Println("Conversion error:", err)
			}
			Transactions, err := c.TransactionService.FindTransactionByIDUser(userIDCnv, page, pageSize)
			if err != nil {
				response := helper.BuildErrorResponse("Failed to fetch data")
				return context.JSON(http.StatusInternalServerError, response)
			}

			total := c.TransactionService.TotalTransactionByUserID(userIDCnv)

			totalPages := (int(total) + pageSize - 1) / pageSize

			customResponse := struct {
				Status  bool                      `json:"status"`
				Message string                    `json:"message"`
				Errors  interface{}               `json:"errors"`
				Data    []entity.Transaction      `json:"data"`
				Paging  helper.PaginationResponse `json:"paging"`
			}{
				Status:  true,
				Message: "OK!",
				Errors:  nil,
				Data:    Transactions,
				Paging:  helper.PaginationResponse{TotalRecords: int(total), CurrentPage: page, TotalPages: totalPages},
			}

			return context.JSON(http.StatusOK, customResponse)
		}
	}
	response := helper.BuildErrorResponse("Invalid token claims")
	return context.JSON(http.StatusBadRequest, response)

}

func (c *transactionController) FindTransactionByID(context echo.Context) error {
	id := context.Param("id")

	orderIDUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		// Handle the error when parsing orderID
		res := helper.BuildErrorResponse("Failed to parse order ID")
		return context.JSON(http.StatusBadRequest, res)
	}

	Transaction := c.TransactionService.FindTransactionByID(orderIDUint)
	response := helper.BuildResponse(true, "OK!", Transaction)
	return context.JSON(http.StatusOK, response)
}
