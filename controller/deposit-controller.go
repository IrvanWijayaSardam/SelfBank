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

type DepositController interface {
	Insert(context echo.Context) error
	All(context echo.Context) error
	FindDepositByID(context echo.Context) error
	HandleMidtransNotification(context echo.Context) error
	Refund(context echo.Context) error
}

type depositController struct {
	DepositService service.DepositService
	jwtService     service.JWTService
}

func NewDepositController(depositService service.DepositService, jwtService service.JWTService) DepositController {
	return &depositController{
		DepositService: depositService,
		jwtService:     jwtService,
	}
}

func (c *depositController) Insert(context echo.Context) error {
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

		var DepositDTO dto.DepositDTO
		if err := context.Bind(&DepositDTO); err != nil {
			response := helper.BuildErrorResponse("Failed to process request")
			return context.JSON(http.StatusBadRequest, response)
		}

		DepositDTO.ID_User, _ = strconv.ParseUint(userID, 10, 64)
		Deposit := c.DepositService.InsertDeposit(DepositDTO)

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

		if bank, ok := bankMap[DepositDTO.PaymentType]; ok {
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
					OrderID:  strconv.FormatUint(Deposit.ID, 10),
					GrossAmt: int64(Deposit.Amount),
				},
			}

			chargeResp, err := coreapi.ChargeTransaction(chargeReq)
			if err != nil {
				c.DepositService.UpdateDepositStatus(Deposit.ID, 3)
				res := helper.BuildErrorResponse("Failed to charge deposit")
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
			c.DepositService.InsertPaymentToken(Deposit.ID, chargeResp.OrderID, vaAccount)

			// Respond with success message and the response data
			res := helper.BuildResponse(true, "Deposit inserted successfully!", response)
			return context.JSON(http.StatusCreated, res)
		} else {
			res := helper.BuildErrorResponse("Unsupported payment type")
			return context.JSON(http.StatusBadRequest, res)
		}
	}

	response := helper.BuildErrorResponse("Invalid JWT token")
	return context.JSON(http.StatusUnauthorized, response)
}

func (c *depositController) All(context echo.Context) error {
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
			Deposits, err := c.DepositService.All(page, pageSize)
			if err != nil {
				response := helper.BuildErrorResponse("Failed to fetch data")
				return context.JSON(http.StatusInternalServerError, response)
			}

			total := c.DepositService.TotalDeposit()

			totalPages := (int(total) + pageSize - 1) / pageSize

			customResponse := struct {
				Status  bool                      `json:"status"`
				Message string                    `json:"message"`
				Errors  interface{}               `json:"errors"`
				Data    []entity.Deposit          `json:"data"`
				Paging  helper.PaginationResponse `json:"paging"`
			}{
				Status:  true,
				Message: "OK!",
				Errors:  nil,
				Data:    Deposits,
				Paging:  helper.PaginationResponse{TotalRecords: int(total), CurrentPage: page, TotalPages: totalPages},
			}

			return context.JSON(http.StatusOK, customResponse)
		case 2:
			userIDCnv, err := strconv.ParseUint(userID, 10, 64)
			if err != nil {
				fmt.Println("Conversion error:", err)
			}
			Deposits, err := c.DepositService.FindDepositByIDUser(userIDCnv, page, pageSize)
			if err != nil {
				response := helper.BuildErrorResponse("Failed to fetch data")
				return context.JSON(http.StatusInternalServerError, response)
			}

			total := c.DepositService.TotalDepositByUserID(userIDCnv)

			totalPages := (int(total) + pageSize - 1) / pageSize

			customResponse := struct {
				Status  bool                      `json:"status"`
				Message string                    `json:"message"`
				Errors  interface{}               `json:"errors"`
				Data    []entity.Deposit          `json:"data"`
				Paging  helper.PaginationResponse `json:"paging"`
			}{
				Status:  true,
				Message: "OK!",
				Errors:  nil,
				Data:    Deposits,
				Paging:  helper.PaginationResponse{TotalRecords: int(total), CurrentPage: page, TotalPages: totalPages},
			}

			return context.JSON(http.StatusOK, customResponse)
		}
	}
	response := helper.BuildErrorResponse("Invalid token claims")
	return context.JSON(http.StatusBadRequest, response)

}

func (c *depositController) Refund(context echo.Context) error {
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
		response := helper.BuildErrorResponse("Token is not valid")
		return context.JSON(http.StatusUnauthorized, response)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		context.Set("user", claims)

		var RefundDTO dto.RefundDTO
		if err := context.Bind(&RefundDTO); err != nil {
			response := helper.BuildErrorResponse("Failed to process request")
			return context.JSON(http.StatusBadRequest, response)
		}

		// Initialize Midtrans client with your server key and environment
		midtrans.ServerKey = MT_SERVER_KEY
		midtrans.ClientKey = MT_CLIENT_KEY
		midtrans.Environment = midtrans.Sandbox

		refundReq := &coreapi.RefundReq{
			RefundKey: "withdrawal22938928",
			Amount:    int64(RefundDTO.Amount),
			Reason:    RefundDTO.Reason,
		}

		orderIDStr := strconv.FormatUint(RefundDTO.OrderID, 10) // Convert the uint64 to a string

		refundResp, err := coreapi.DirectRefundTransaction(orderIDStr, refundReq)
		if err != nil {
			// Handle the error when refunding the deposit
			res := helper.BuildErrorResponse("Failed to refund deposit")
			context.JSON(http.StatusInternalServerError, res)
			return err
		}

		// Assuming you have a "response" map to store the response data
		response := make(map[string]interface{})
		response["refund_status"] = refundResp.StatusMessage
		response["refund_id"] = refundResp.ID

		// Respond with success message and the response data
		res := helper.BuildResponse(true, "Refund processed successfully!", response)
		return context.JSON(http.StatusOK, res)
	}

	response := helper.BuildErrorResponse("Invalid JWT token")
	return context.JSON(http.StatusUnauthorized, response)
}

func (c *depositController) FindDepositByID(context echo.Context) error {
	id := context.Param("id")

	orderIDUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		// Handle the error when parsing orderID
		res := helper.BuildErrorResponse("Failed to parse order ID")
		return context.JSON(http.StatusBadRequest, res)
	}

	Deposit := c.DepositService.FindDepositByID(orderIDUint)
	response := helper.BuildResponse(true, "OK!", Deposit)
	return context.JSON(http.StatusOK, response)
}

func (c *depositController) HandleMidtransNotification(ctx echo.Context) error {
	var notificationPayload map[string]interface{}

	// 1. Parse JSON request body
	if err := json.NewDecoder(ctx.Request().Body).Decode(&notificationPayload); err != nil {
		// Handle the error when decoding the JSON payload
		res := helper.BuildErrorResponse("Failed to parse notification")
		return ctx.JSON(http.StatusBadRequest, res)
	}

	// 2. Get order ID from the payload
	orderID, exists := notificationPayload["order_id"].(string)
	if !exists {
		// Handle the case when the key `order_id` is not found
		res := helper.BuildErrorResponse("Order ID not found in notification")
		return ctx.JSON(http.StatusBadRequest, res)
	}

	depositStatusResp, midErr := coreapi.CheckTransaction(orderID)
	if midErr != nil {
		// Handle the error when checking deposit status using Midtrans error type
		res := helper.BuildErrorResponse("Failed to check deposit status" + midErr.Message)
		return ctx.JSON(http.StatusInternalServerError, res)
	}

	// 4. Update deposit status in your database based on the response
	if depositStatusResp != nil {
		status := ""
		switch depositStatusResp.TransactionStatus {
		case "capture":
			if depositStatusResp.FraudStatus == "challenge" {
				// Set deposit status on your database to 'challenge'
				status = "challenge"
			} else if depositStatusResp.FraudStatus == "accept" {
				// Set deposit status on your database to 'success'
				status = "success"
			}
		case "settlement":
			orderIDUint, err := strconv.ParseUint(orderID, 10, 64)
			if err != nil {
				// Handle the error when parsing orderID
				res := helper.BuildErrorResponse("Failed to parse order ID" + err.Error())
				return ctx.JSON(http.StatusBadRequest, res)
			}

			// Update the status of MasterJual to "2" in your database
			err = c.DepositService.UpdateDepositStatus(orderIDUint, 5)
			if err != nil {
				// Handle the error when updating the status
				res := helper.BuildErrorResponse("Failed to update deposit status" + err.Error())
				return ctx.JSON(http.StatusInternalServerError, res)
			}
		case "deny":
			orderIDUint, err := strconv.ParseUint(orderID, 10, 64)
			if err != nil {
				// Handle the error when parsing orderID
				res := helper.BuildErrorResponse("Failed to parse order ID" + err.Error())
				return ctx.JSON(http.StatusBadRequest, res)
			}

			// Update the status of MasterJual to "2" in your database
			err = c.DepositService.UpdateDepositStatus(orderIDUint, 4)
			if err != nil {
				// Handle the error when updating the status
				res := helper.BuildErrorResponse("Failed to update deposit status" + err.Error())
				return ctx.JSON(http.StatusInternalServerError, res)
			}
		case "cancel", "expire":
			orderIDUint, err := strconv.ParseUint(orderID, 10, 64)
			if err != nil {
				// Handle the error when parsing orderID
				res := helper.BuildErrorResponse("Failed to parse order ID" + err.Error())
				return ctx.JSON(http.StatusBadRequest, res)
			}

			// Update the status of MasterJual to "2" in your database
			err = c.DepositService.UpdateDepositStatus(orderIDUint, 3)
			if err != nil {
				// Handle the error when updating the status
				res := helper.BuildErrorResponse("Failed to update deposit status" + err.Error())
				return ctx.JSON(http.StatusInternalServerError, res)
			}
		case "pending":
			orderIDUint, err := strconv.ParseUint(orderID, 10, 64)
			if err != nil {
				// Handle the error when parsing orderID
				res := helper.BuildErrorResponse("Failed to parse order ID" + err.Error())
				return ctx.JSON(http.StatusBadRequest, res)
			}

			// Update status transaski ke 4 , status pending
			err = c.DepositService.UpdateDepositStatus(orderIDUint, 2)
			if err != nil {
				// Handle the error when updating the status
				res := helper.BuildErrorResponse("Failed to update deposit status" + err.Error())
				return ctx.JSON(http.StatusInternalServerError, res)
			}
		}

		if status == "success" {
			orderIDUint, err := strconv.ParseUint(orderID, 10, 64)
			if err != nil {
				// Handle the error when parsing orderID
				res := helper.BuildErrorResponse("Failed to parse order ID" + err.Error())
				return ctx.JSON(http.StatusBadRequest, res)
			}

			// Update status transaski ke 5 , status sukses
			err = c.DepositService.UpdateDepositStatus(orderIDUint, 5)
			if err != nil {
				// Handle the error when updating the status
				res := helper.BuildErrorResponse("Failed to update deposit status" + err.Error())
				return ctx.JSON(http.StatusInternalServerError, res)
			}
		}
	}

	return ctx.JSON(http.StatusOK, map[string]string{"status": "ok"})
}