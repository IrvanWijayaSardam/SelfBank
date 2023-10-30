package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"

	"github.com/IrvanWijayaSardam/SelfBank/dto"
	"github.com/IrvanWijayaSardam/SelfBank/helper"
	"github.com/IrvanWijayaSardam/SelfBank/service"
)

type WithdrawalController interface {
	Insert(context echo.Context) error
	All(context echo.Context) error
	FindWithdrawalByID(context echo.Context) error
}

type withdrawalController struct {
	WithdrawalService service.WithdrawalService
	UserService       service.UserService
	jwtService        service.JWTService
}

func NewWithdrawalController(withdrawalService service.WithdrawalService, userService service.UserService, jwtService service.JWTService) WithdrawalController {
	return &withdrawalController{
		WithdrawalService: withdrawalService,
		UserService:       userService,
		jwtService:        jwtService,
	}
}

func (c *withdrawalController) Insert(context echo.Context) error {
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

		var WithdrawalDTO dto.WithdrawalDTO
		if err := context.Bind(&WithdrawalDTO); err != nil {
			response := helper.BuildErrorResponse("Failed to process request")
			return context.JSON(http.StatusBadRequest, response)
		}

		WithdrawalDTO.ID_User, _ = strconv.ParseUint(userID, 10, 64)
		currentSaldo := c.UserService.GetSaldo(WithdrawalDTO.ID_User)

		if currentSaldo >= int64(WithdrawalDTO.Amount) {
			Withdrawal := c.WithdrawalService.InsertWithdrawal(WithdrawalDTO)
			res := helper.BuildResponse(true, "Withdrawal Success", Withdrawal)
			return context.JSON(http.StatusCreated, res)
		} else {
			res := helper.BuildErrorResponse("Cannot continue withdrawal because your balance is insufficient")
			return context.JSON(http.StatusBadRequest, res)
		}

	} else {
		res := helper.BuildErrorResponse("There's Something Wrong Contact Developer")
		return context.JSON(http.StatusBadRequest, res)
	}
}

func (c *withdrawalController) All(context echo.Context) error {
	pageParam := context.QueryParam("page")
	pageSizeParam := context.QueryParam("pageSize")
	exportTo := context.QueryParam("exportTo")

	var withdrawalResponse []dto.WithdrawalResponseDTO

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
			Withdrawals, err := c.WithdrawalService.All(page, pageSize)
			if err != nil {
				response := helper.BuildErrorResponse("Failed to fetch data")
				return context.JSON(http.StatusInternalServerError, response)
			}

			for _, transaction := range Withdrawals {
				response := dto.WithdrawalResponseDTO{
					ID:     transaction.ID,
					IDUser: transaction.ID_User,
					Date:   helper.ConvertUnixtime(transaction.Date).Format("2006-01-02 15:04:05"),
					Amount: transaction.Amount,
					Status: transaction.Status,
					To:     transaction.To,
				}
				withdrawalResponse = append(withdrawalResponse, response)
			}

			if exportTo == "pdf" {
				pdfBuffer, err := c.WithdrawalService.GenerateWithdrawalPDF(Withdrawals)
				if err != nil {
					response := helper.BuildErrorResponse("Failed to generate PDF")
					return context.JSON(http.StatusInternalServerError, response)
				}

				pdfFileName := "transactions.pdf"

				// Set the response headers to force download
				context.Response().Header().Set("Content-Disposition", "attachment; filename="+pdfFileName)
				context.Response().Header().Set("Content-Type", "application/pdf")

				// Write the PDF from the buffer to the response writer
				_, err = pdfBuffer.WriteTo(context.Response())
				if err != nil {
					response := helper.BuildErrorResponse("Failed to write PDF to response")
					return context.JSON(http.StatusInternalServerError, response)
				}
			}

			total := c.WithdrawalService.TotalWithdrawal()

			totalPages := (int(total) + pageSize - 1) / pageSize

			customResponse := struct {
				Status  bool                        `json:"status"`
				Message string                      `json:"message"`
				Data    []dto.WithdrawalResponseDTO `json:"data"`
				Paging  helper.PaginationResponse   `json:"paging"`
			}{
				Status:  true,
				Message: "OK!",
				Data:    withdrawalResponse,
				Paging:  helper.PaginationResponse{TotalRecords: int(total), CurrentPage: page, TotalPages: totalPages},
			}

			return context.JSON(http.StatusOK, customResponse)
		case 2:
			userIDCnv, err := strconv.ParseUint(userID, 10, 64)
			if err != nil {
				fmt.Println("Conversion error:", err)
			}
			Withdrawals, err := c.WithdrawalService.FindWithdrawalByIDUser(userIDCnv, page, pageSize)
			if err != nil {
				response := helper.BuildErrorResponse("Failed to fetch data")
				return context.JSON(http.StatusInternalServerError, response)
			}
			for _, transaction := range Withdrawals {
				response := dto.WithdrawalResponseDTO{
					ID:     transaction.ID,
					IDUser: transaction.ID_User,
					Date:   helper.ConvertUnixtime(transaction.Date).Format("2006-01-02 15:04:05"),
					Amount: transaction.Amount,
					To:     transaction.To,
					Status: transaction.Status,
				}
				withdrawalResponse = append(withdrawalResponse, response)
			}

			if exportTo == "pdf" {
				pdfBuffer, err := c.WithdrawalService.GenerateWithdrawalPDF(Withdrawals)
				if err != nil {
					response := helper.BuildErrorResponse("Failed to generate PDF")
					return context.JSON(http.StatusInternalServerError, response)
				}

				pdfFileName := "transactions.pdf"

				// Set the response headers to force download
				context.Response().Header().Set("Content-Disposition", "attachment; filename="+pdfFileName)
				context.Response().Header().Set("Content-Type", "application/pdf")

				// Write the PDF from the buffer to the response writer
				_, err = pdfBuffer.WriteTo(context.Response())
				if err != nil {
					response := helper.BuildErrorResponse("Failed to write PDF to response")
					return context.JSON(http.StatusInternalServerError, response)
				}
			}

			total := c.WithdrawalService.TotalWithdrawalByUserID(userIDCnv)

			totalPages := (int(total) + pageSize - 1) / pageSize

			customResponse := struct {
				Status  bool                        `json:"status"`
				Message string                      `json:"message"`
				Data    []dto.WithdrawalResponseDTO `json:"data"`
				Paging  helper.PaginationResponse   `json:"paging"`
			}{
				Status:  true,
				Message: "OK!",
				Data:    withdrawalResponse,
				Paging:  helper.PaginationResponse{TotalRecords: int(total), CurrentPage: page, TotalPages: totalPages},
			}

			return context.JSON(http.StatusOK, customResponse)
		}
	}
	response := helper.BuildErrorResponse("Invalid token claims")
	return context.JSON(http.StatusBadRequest, response)

}

func (c *withdrawalController) FindWithdrawalByID(context echo.Context) error {
	id := context.Param("id")

	orderIDUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		// Handle the error when parsing orderID
		res := helper.BuildErrorResponse("Failed to parse order ID")
		return context.JSON(http.StatusBadRequest, res)
	}

	Withdrawal := c.WithdrawalService.FindWithdrawalByID(orderIDUint)
	if Withdrawal.ID == 0 {
		res := helper.BuildErrorResponse("Withdrawal not found")
		return context.JSON(http.StatusNotFound, res)
	} else {
		var data = dto.WithdrawalResponseDTO{
			ID:     Withdrawal.ID,
			IDUser: Withdrawal.ID_User,
			Date:   helper.ConvertUnixtime(Withdrawal.Date).Format("2006-01-02 15:04:05"),
			Amount: Withdrawal.Amount,
			To:     Withdrawal.To,
			Status: Withdrawal.Status,
		}

		response := helper.BuildResponse(true, "OK!", data)
		return context.JSON(http.StatusOK, response)
	}
}
