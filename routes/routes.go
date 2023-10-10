package routes

import (
	"github.com/IrvanWijayaSardam/SelfBank/controller"
	"github.com/IrvanWijayaSardam/SelfBank/service"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, jwtService service.JWTService, authController controller.AuthController) {
	authRoutes := e.Group("/api/auth")

	authRoutes.POST("/login", authController.Login)
	authRoutes.POST("/register", authController.Register)
}

func TransactionRoutes(e *echo.Echo, transactionService service.TransactionService,
	transactionController controller.TransactionController, jwtMiddleware echo.MiddlewareFunc) {
	authRoutes := e.Group("/api/transaction")

	authRoutes.Use(jwtMiddleware)

	authRoutes.POST("/", transactionController.Insert)
	authRoutes.GET("/", transactionController.All)
	authRoutes.GET("/:id", transactionController.FindTransactionByID)

}

func MidtransRoutes(e *echo.Echo, transactionService service.TransactionService,
	transactionController controller.TransactionController, jwtMiddleware echo.MiddlewareFunc) {
	authRoutes := e.Group("/api/midtrans/notification")

	authRoutes.POST("/", transactionController.HandleMidtransNotification)

}
