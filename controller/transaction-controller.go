package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"github.com/IrvanWijayaSardam/SelfBank/dto"
	"github.com/IrvanWijayaSardam/SelfBank/entity"
	"github.com/IrvanWijayaSardam/SelfBank/helper"
	"github.com/IrvanWijayaSardam/SelfBank/service"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type TransactionController interface {
	Insert(context echo.Context) error
	All(context echo.Context) error
	FindTransactionByID(context echo.Context) error
	HandleMidtransNotification(context echo.Context) error
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
	errEnv := godotenv.Load()
	if errEnv != nil {
		panic("Failed to load env file")
	}

	MT_SERVER_KEY := os.Getenv("MT_SERVER_KEY")
	MT_CLIENT_KEY := os.Getenv("MT_CLIENT_KEY")

	token, err := c.jwtService.ValidateToken(authHeader)
	if err != nil {
		log.Println(err)
		response := helper.BuildErrorResponse("Token is not valid", err.Error(), nil)
		return context.JSON(http.StatusUnauthorized, response)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		context.Set("user", claims)

		userID, ok := claims["userid"].(string)
		if !ok {
			response := helper.BuildErrorResponse("Failed to process request", "UserID not found in claims", helper.EmptyObj{})
			return context.JSON(http.StatusBadRequest, response)
		}

		var TransactionDTO dto.TransactionDTO
		if err := context.Bind(&TransactionDTO); err != nil {
			response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
			return context.JSON(http.StatusBadRequest, response)
		}

		TransactionDTO.ID_User, _ = strconv.ParseUint(userID, 10, 64)
		Transaction := c.TransactionService.InsertTransaction(TransactionDTO)

		// Initialize Midtrans client with your server key and environment
		midtrans.ServerKey = MT_SERVER_KEY
		midtrans.ClientKey = MT_CLIENT_KEY
		midtrans.Environment = midtrans.Sandbox

		// Create a map that maps IdPayment to the corresponding bank name
		bankMap := map[string]string{
			"6": "bca",
			"7": "bri",
			"8": "bni",
			// Add more mappings as needed
		}

		if bank, ok := bankMap[TransactionDTO.PaymentType]; ok {
			var midtransBank midtrans.Bank

			switch bank {
			case "bca":
				midtransBank = midtrans.BankBca
			case "bri":
				midtransBank = midtrans.BankBri
			case "bni":
				midtransBank = midtrans.BankBni
			default:
				midtransBank = midtrans.BankBca
			}

			chargeReq := &coreapi.ChargeReq{
				PaymentType:  "bank_transfer",
				BankTransfer: &coreapi.BankTransferDetails{Bank: midtransBank},
				TransactionDetails: midtrans.TransactionDetails{
					OrderID:  strconv.FormatUint(Transaction.ID, 10),
					GrossAmt: int64(Transaction.Ammount),
				},
			}

			chargeResp, err := coreapi.ChargeTransaction(chargeReq)
			if err != nil {
				res := helper.BuildErrorResponse("Failed to charge transaction", err.Error(), helper.EmptyObj{})
				context.JSON(http.StatusInternalServerError, res)
				return err
			}

			var vaAccount string
			for _, va := range chargeResp.VaNumbers {
				if va.Bank == bank {
					vaAccount = va.VANumber
					break
				}
			}

			// Assuming you have a "response" map to store the response data
			response := make(map[string]interface{})
			response["va_account"] = vaAccount
			c.TransactionService.InsertPaymentToken(Transaction.ID, chargeResp.OrderID, vaAccount)

			// Respond with success message and the response data
			res := helper.BuildResponse(true, "Transaction inserted successfully!", response)
			return context.JSON(http.StatusCreated, res)
		} else {
			// Handle unsupported IdPayment values
			res := helper.BuildErrorResponse("Unsupported payment type", "Unsupported payment type", helper.EmptyObj{})
			return context.JSON(http.StatusBadRequest, res)
		}
	}

	response := helper.BuildErrorResponse("Failed to process request", "Invalid JWT token", helper.EmptyObj{})
	return context.JSON(http.StatusUnauthorized, response)
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
		response := helper.BuildErrorResponse("Token is not valid", err.Error(), nil)
		return context.JSON(http.StatusUnauthorized, response)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		context.Set("user", claims)

		userID, ok := claims["userid"].(string)
		if !ok {
			response := helper.BuildErrorResponse("Failed to process request", "UserID not found in claims", helper.EmptyObj{})
			return context.JSON(http.StatusBadRequest, response)
		}

		roleID, ok := claims["idrole"].(float64)
		if !ok {
			response := helper.BuildErrorResponse("Failed to process request", "IDRole not found in claims", helper.EmptyObj{})
			return context.JSON(http.StatusBadRequest, response)
		}
		switch roleID {
		case 1:
			Transactions, err := c.TransactionService.All(page, pageSize)
			if err != nil {
				response := helper.BuildErrorResponse("Failed to fetch data", err.Error(), helper.EmptyObj{})
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
				response := helper.BuildErrorResponse("Failed to fetch data", err.Error(), helper.EmptyObj{})
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
	response := helper.BuildErrorResponse("Failed to process request", "Invalid token claims", helper.EmptyObj{})
	return context.JSON(http.StatusBadRequest, response)

}

func (c *transactionController) FindTransactionByID(context echo.Context) error {
	id := context.Param("id")

	orderIDUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		// Handle the error when parsing orderID
		res := helper.BuildErrorResponse("Failed to parse order ID", err.Error(), helper.EmptyObj{})
		return context.JSON(http.StatusBadRequest, res)
	}

	Transaction := c.TransactionService.FindTransactionByID(orderIDUint)
	response := helper.BuildResponse(true, "OK!", Transaction)
	return context.JSON(http.StatusOK, response)
}

func (c *transactionController) HandleMidtransNotification(ctx echo.Context) error {
	var notificationPayload map[string]interface{}

	// 1. Parse JSON request body
	if err := json.NewDecoder(ctx.Request().Body).Decode(&notificationPayload); err != nil {
		// Handle the error when decoding the JSON payload
		res := helper.BuildErrorResponse("Failed to parse notification", err.Error(), helper.EmptyObj{})
		return ctx.JSON(http.StatusBadRequest, res)
	}

	// 2. Get order ID from the payload
	orderID, exists := notificationPayload["order_id"].(string)
	if !exists {
		// Handle the case when the key `order_id` is not found
		res := helper.BuildErrorResponse("Order ID not found in notification", "", helper.EmptyObj{})
		return ctx.JSON(http.StatusBadRequest, res)
	}

	transactionStatusResp, midErr := coreapi.CheckTransaction(orderID)
	if midErr != nil {
		// Handle the error when checking transaction status using Midtrans error type
		res := helper.BuildErrorResponse("Failed to check transaction status", midErr.Message, helper.EmptyObj{})
		return ctx.JSON(http.StatusInternalServerError, res)
	}

	// 4. Update transaction status in your database based on the response
	if transactionStatusResp != nil {
		status := ""
		switch transactionStatusResp.TransactionStatus {
		case "capture":
			if transactionStatusResp.FraudStatus == "challenge" {
				// Set transaction status on your database to 'challenge'
				status = "challenge"
			} else if transactionStatusResp.FraudStatus == "accept" {
				// Set transaction status on your database to 'success'
				status = "success"
			}
		case "settlement":
			orderIDUint, err := strconv.ParseUint(orderID, 10, 64)
			if err != nil {
				// Handle the error when parsing orderID
				res := helper.BuildErrorResponse("Failed to parse order ID", err.Error(), helper.EmptyObj{})
				return ctx.JSON(http.StatusBadRequest, res)
			}

			// Update the status of MasterJual to "2" in your database
			err = c.TransactionService.UpdateTransactionStatus(orderIDUint, 5)
			if err != nil {
				// Handle the error when updating the status
				res := helper.BuildErrorResponse("Failed to update transaction status", err.Error(), helper.EmptyObj{})
				return ctx.JSON(http.StatusInternalServerError, res)
			}
		case "deny":
			orderIDUint, err := strconv.ParseUint(orderID, 10, 64)
			if err != nil {
				// Handle the error when parsing orderID
				res := helper.BuildErrorResponse("Failed to parse order ID", err.Error(), helper.EmptyObj{})
				return ctx.JSON(http.StatusBadRequest, res)
			}

			// Update the status of MasterJual to "2" in your database
			err = c.TransactionService.UpdateTransactionStatus(orderIDUint, 4)
			if err != nil {
				// Handle the error when updating the status
				res := helper.BuildErrorResponse("Failed to update transaction status", err.Error(), helper.EmptyObj{})
				return ctx.JSON(http.StatusInternalServerError, res)
			}
		case "cancel", "expire":
			orderIDUint, err := strconv.ParseUint(orderID, 10, 64)
			if err != nil {
				// Handle the error when parsing orderID
				res := helper.BuildErrorResponse("Failed to parse order ID", err.Error(), helper.EmptyObj{})
				return ctx.JSON(http.StatusBadRequest, res)
			}

			// Update the status of MasterJual to "2" in your database
			err = c.TransactionService.UpdateTransactionStatus(orderIDUint, 3)
			if err != nil {
				// Handle the error when updating the status
				res := helper.BuildErrorResponse("Failed to update transaction status", err.Error(), helper.EmptyObj{})
				return ctx.JSON(http.StatusInternalServerError, res)
			}
		case "pending":
			orderIDUint, err := strconv.ParseUint(orderID, 10, 64)
			if err != nil {
				// Handle the error when parsing orderID
				res := helper.BuildErrorResponse("Failed to parse order ID", err.Error(), helper.EmptyObj{})
				return ctx.JSON(http.StatusBadRequest, res)
			}

			// Update status transaski ke 4 , status pending
			err = c.TransactionService.UpdateTransactionStatus(orderIDUint, 2)
			if err != nil {
				// Handle the error when updating the status
				res := helper.BuildErrorResponse("Failed to update transaction status", err.Error(), helper.EmptyObj{})
				return ctx.JSON(http.StatusInternalServerError, res)
			}
		}

		if status == "success" {
			orderIDUint, err := strconv.ParseUint(orderID, 10, 64)
			if err != nil {
				// Handle the error when parsing orderID
				res := helper.BuildErrorResponse("Failed to parse order ID", err.Error(), helper.EmptyObj{})
				return ctx.JSON(http.StatusBadRequest, res)
			}

			// Update status transaski ke 5 , status sukses
			err = c.TransactionService.UpdateTransactionStatus(orderIDUint, 5)
			if err != nil {
				// Handle the error when updating the status
				res := helper.BuildErrorResponse("Failed to update transaction status", err.Error(), helper.EmptyObj{})
				return ctx.JSON(http.StatusInternalServerError, res)
			}
		}
	}

	return ctx.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
