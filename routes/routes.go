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

func DepositRoutes(e *echo.Echo, depositService service.DepositService,
	depositController controller.DepositController, jwtMiddleware echo.MiddlewareFunc) {
	authRoutes := e.Group("/api/deposit")

	authRoutes.Use(jwtMiddleware)
	authRoutes.POST("/", depositController.Insert)
	authRoutes.GET("/", depositController.All)
	authRoutes.GET("/:id", depositController.FindDepositByID)

}

func WithdrawalRoutes(e *echo.Echo, withdrawalService service.WithdrawalService,
	withdrawalController controller.WithdrawalController, jwtMiddleware echo.MiddlewareFunc) {
	authRoutes := e.Group("/api/withdrawal")

	authRoutes.Use(jwtMiddleware)

	authRoutes.POST("/", withdrawalController.Insert)
	authRoutes.GET("/", withdrawalController.All)
	authRoutes.GET("/:id", withdrawalController.FindWithdrawalByID)

}

func MidtransRoutes(e *echo.Echo, transactionService service.DepositService,
	transactionController controller.DepositController, jwtMiddleware echo.MiddlewareFunc) {
	authRoutes := e.Group("/api/midtrans/notifications")

	authRoutes.POST("/", transactionController.HandleMidtransNotification)

}
