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
	// errEnv := godotenv.Load()
	// if errEnv != nil {
	// 	panic("Failed to load env file")
	// }

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
					OrderID:  Deposit.ID,
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

			response := make(map[string]interface{})
			response["va_account"] = vaAccount
			c.DepositService.InsertPaymentToken(Deposit.ID, chargeResp.TransactionID, vaAccount, "-")

			res := helper.BuildResponse(true, "Deposit inserted successfully!", response)
			return context.JSON(http.StatusCreated, res)
		} else if DepositDTO.PaymentType == "10" {
			chargeReq := &coreapi.ChargeReq{
				Gopay: &coreapi.GopayDetails{
					EnableCallback: true,
				},
				PaymentType: "gopay",
				TransactionDetails: midtrans.TransactionDetails{
					OrderID:  Deposit.ID,
					GrossAmt: int64(Deposit.Amount),
				},
			}

			chargeResp, err := coreapi.ChargeTransaction(chargeReq)
			if err != nil {
				res := helper.BuildErrorResponse("Failed to charge transaction")
				return context.JSON(http.StatusInternalServerError, res)
			}
			response := make(map[string]interface{})

			if len(chargeResp.Actions) > 0 {
				for _, action := range chargeResp.Actions {
					if action.Name == "deeplink-redirect" {
						deepLinkURL := action.URL
						response["callback_url"] = deepLinkURL
						c.DepositService.InsertPaymentToken(Deposit.ID, chargeResp.TransactionID, "-", deepLinkURL)
						break
					}
				}
			}
			res := helper.BuildResponse(true, "Transaction inserted successfully!", response)
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
	authHeader := context.Request().Header.Get("Authorization")
	pageParam := context.QueryParam("page")
	pageSizeParam := context.QueryParam("pageSize")
	exportTo := context.QueryParam("exportTo")

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

			var depositResponses []dto.DepositResponse

			for _, deposit := range Deposits {
				status := ""
				switch deposit.Status {
				case 1:
					status = "Created"
				case 2:
					status = "Pending"
				case 3:
					status = "Cancelled"
				case 4:
					status = "Denied"
				case 5:
					status = "Paid"
				default:
					status = "Created"
				}

				paymentInfo := c.DepositService.FindPaymentInfoById(deposit.ID)

				depositResponse := dto.DepositResponse{
					Id_deposit:      deposit.ID,
					Id_user:         deposit.ID_User,
					Virtual_account: paymentInfo.VirtualAcc,
					Url_callback:    paymentInfo.CallbackUrl,
					Amount:          deposit.Amount,
					Status:          status,
					Date:            deposit.Date,
				}
				depositResponses = append(depositResponses, depositResponse)
			}
			if exportTo == "pdf" {
				pdfBuffer, err := c.DepositService.GenerateDepositPDF(depositResponses)
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

			total := c.DepositService.TotalDeposit()

			totalPages := (int(total) + pageSize - 1) / pageSize

			customResponse := struct {
				Status  bool                      `json:"status"`
				Message string                    `json:"message"`
				Data    []dto.DepositResponse     `json:"data"`
				Paging  helper.PaginationResponse `json:"paging"`
			}{
				Status:  true,
				Message: "OK!",
				Data:    depositResponses,
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
			var depositResponses []dto.DepositResponse

			for _, deposit := range Deposits {
				status := ""
				switch deposit.Status {
				case 1:
					status = "Created"
				case 2:
					status = "Pending"
				case 3:
					status = "Cancelled"
				case 4:
					status = "Denied"
				case 5:
					status = "Paid"
				default:
					status = "Created"
				}

				paymentInfo := c.DepositService.FindPaymentInfoById(deposit.ID)

				depositResponse := dto.DepositResponse{
					Id_deposit:      deposit.ID,
					Id_user:         deposit.ID_User,
					Virtual_account: paymentInfo.VirtualAcc,
					Url_callback:    paymentInfo.CallbackUrl,
					Amount:          deposit.Amount,
					Status:          status,
					Date:            deposit.Date,
				}
				depositResponses = append(depositResponses, depositResponse)
			}

			if exportTo == "pdf" {
				pdfBuffer, err := c.DepositService.GenerateDepositPDF(depositResponses)
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
			total := c.DepositService.TotalDepositByUserID(userIDCnv)

			totalPages := (int(total) + pageSize - 1) / pageSize

			customResponse := struct {
				Status  bool                      `json:"status"`
				Message string                    `json:"message"`
				Data    []dto.DepositResponse     `json:"data"`
				Paging  helper.PaginationResponse `json:"paging"`
			}{
				Status:  true,
				Message: "OK!",
				Data:    depositResponses,
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

	Deposit := c.DepositService.FindDepositByID(id)

	status := ""
	switch Deposit.Status {
	case 1:
		status = "Created"
	case 2:
		status = "Pending"
	case 3:
		status = "Cancelled"
	case 4:
		status = "Denied"
	case 5:
		status = "Paid"
	default:
		status = "Created"
	}

	paymentInfo := c.DepositService.FindPaymentInfoById(Deposit.ID)

	depositResponse := dto.DepositResponse{
		Id_deposit:      Deposit.ID,
		Id_user:         Deposit.ID_User,
		Virtual_account: paymentInfo.VirtualAcc,
		Url_callback:    paymentInfo.CallbackUrl,
		Amount:          Deposit.Amount,
		Status:          status,
	}

	response := helper.BuildResponse(true, "OK!", depositResponse)
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

	if depositStatusResp != nil {
		status := ""
		switch depositStatusResp.TransactionStatus {
		case "capture":
			if depositStatusResp.FraudStatus == "challenge" {
				status = "challenge"
			} else if depositStatusResp.FraudStatus == "accept" {
				status = "success"
			}
		case "settlement":
			err := c.DepositService.UpdateDepositStatus(orderID, 5)
			if err != nil {
				// Handle the error when updating the status
				res := helper.BuildErrorResponse("Failed to update deposit status" + err.Error())
				return ctx.JSON(http.StatusInternalServerError, res)
			}
		case "deny":
			err := c.DepositService.UpdateDepositStatus(orderID, 4)
			if err != nil {
				// Handle the error when updating the status
				res := helper.BuildErrorResponse("Failed to update deposit status" + err.Error())
				return ctx.JSON(http.StatusInternalServerError, res)
			}
		case "cancel", "expire":
			err := c.DepositService.UpdateDepositStatus(orderID, 3)
			if err != nil {
				// Handle the error when updating the status
				res := helper.BuildErrorResponse("Failed to update deposit status" + err.Error())
				return ctx.JSON(http.StatusInternalServerError, res)
			}
		case "pending":
			err := c.DepositService.UpdateDepositStatus(orderID, 2)
			if err != nil {
				// Handle the error when updating the status
				res := helper.BuildErrorResponse("Failed to update deposit status" + err.Error())
				return ctx.JSON(http.StatusInternalServerError, res)
			}
		}

		if status == "success" {
			err := c.DepositService.UpdateDepositStatus(orderID, 5)
			if err != nil {
				// Handle the error when updating the status
				res := helper.BuildErrorResponse("Failed to update deposit status" + err.Error())
				return ctx.JSON(http.StatusInternalServerError, res)
			}
		}
	}

	return ctx.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
