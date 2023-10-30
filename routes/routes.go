package routes

import (
	"github.com/IrvanWijayaSardam/SelfBank/controller"
	"github.com/IrvanWijayaSardam/SelfBank/service"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, jwtService service.JWTService, authController controller.AuthController) {
	registerRoutes := e.Group("/api/auth")

	registerRoutes.POST("/login", authController.Login)
	registerRoutes.POST("/register", authController.Register)
}

func DepositRoutes(e *echo.Echo, depositService service.DepositService,
	depositController controller.DepositController, jwtMiddleware echo.MiddlewareFunc) {
	depositRoutes := e.Group("/api/deposit")

	depositRoutes.Use(jwtMiddleware)
	depositRoutes.POST("/", depositController.Insert)
	depositRoutes.GET("/", depositController.All)
	depositRoutes.GET("/:id", depositController.FindDepositByID)

}

func WithdrawalRoutes(e *echo.Echo, withdrawalService service.WithdrawalService,
	withdrawalController controller.WithdrawalController, jwtMiddleware echo.MiddlewareFunc) {
	withdrawalRoutes := e.Group("/api/withdrawal")

	withdrawalRoutes.Use(jwtMiddleware)

	withdrawalRoutes.POST("/", withdrawalController.Insert)
	withdrawalRoutes.GET("/", withdrawalController.All)
	withdrawalRoutes.GET("/:id", withdrawalController.FindWithdrawalByID)

}

func TransactionRoutes(e *echo.Echo, transactionService service.TransactionService,
	transactionController controller.TransactionController, jwtMiddleware echo.MiddlewareFunc) {
	trxRoutes := e.Group("/api/transaction")

	trxRoutes.Use(jwtMiddleware)

	trxRoutes.POST("/", transactionController.Insert)
	trxRoutes.GET("/", transactionController.All)
	trxRoutes.GET("/:id", transactionController.FindTransactionByID)
}

func ProfileRoutes(e *echo.Echo, userService service.UserService,
	userController controller.UserController, jwtMiddleware echo.MiddlewareFunc) {
	profileRoutes := e.Group("/api/profile")

	profileRoutes.GET("/", userController.MyProfile)
	profileRoutes.PUT("/", userController.UpdateProfile)
	profileRoutes.DELETE("/:id", userController.DeleteUser)

}

func UserRoutes(e *echo.Echo, userService service.UserService,
	userController controller.UserController, jwtMiddleware echo.MiddlewareFunc) {
	profileRoutes := e.Group("/api/user")

	profileRoutes.GET("/", userController.All)
	profileRoutes.PUT("/", userController.UpdateProfile)
	profileRoutes.DELETE("/:id", userController.DeleteUser)

}

func ImageRoutes(e *echo.Echo, userController controller.UserController, jwtMiddleware echo.MiddlewareFunc) {
	imageRoutes := e.Group("/api/cdn/images")

	imageRoutes.POST("/file", userController.FileUpload)

}

func MidtransRoutes(e *echo.Echo, transactionService service.DepositService,
	transactionController controller.DepositController, jwtMiddleware echo.MiddlewareFunc) {
	authRoutes := e.Group("/api/midtrans/notifications")

	authRoutes.POST("/", transactionController.HandleMidtransNotification)

}
