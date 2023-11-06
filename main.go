package main

import (
	"github.com/IrvanWijayaSardam/SelfBank/config"
	"github.com/IrvanWijayaSardam/SelfBank/controller"
	"github.com/IrvanWijayaSardam/SelfBank/helper"
	"github.com/IrvanWijayaSardam/SelfBank/middleware"
	"github.com/IrvanWijayaSardam/SelfBank/repository"
	"github.com/IrvanWijayaSardam/SelfBank/routes"
	"github.com/IrvanWijayaSardam/SelfBank/service"
	"github.com/go-redis/redis"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

var (
	db          *gorm.DB      = config.SetupDatabaseConnection()
	client                    = openai.NewClient(config.EnvOpenAIKey())
	redisClient *redis.Client = config.ConnectRedis()

	userRepository         repository.UserRepository         = repository.NewUserRepository(db)
	depositRepository      repository.DepositRepository      = repository.NewDepositRepository(db)
	withdrawalRepository   repository.WithdrawalRepository   = repository.NewWithdrawalRepository(db)
	transactionRepository  repository.TransactionRepository  = repository.NewTransactionRepository(db)
	chatbotRepository      repository.ChatbotRepository      = repository.NewChatbotRepository(client)
	verificationRepository repository.VerificationRepository = repository.NewVerificationRepository(redisClient, db)

	authService         service.AuthService         = service.NewAuthService(userRepository)
	jwtService          service.JWTService          = service.NewJWTService()
	depositService      service.DepositService      = service.NewDepositService(depositRepository)
	withdrawalService   service.WithdrawalService   = service.NewWithdrawalService(withdrawalRepository)
	userService         service.UserService         = service.NewUserService(userRepository)
	transactionService  service.TransactionService  = service.NewTransactionService(transactionRepository)
	chatbotService      service.ChatbotService      = service.NewChatbotService(chatbotRepository)
	verificationService service.VerificationService = service.NewVerificationService(verificationRepository)
)

func main() {
	e := echo.New()
	e.Debug = true
	jwtMiddleware := middleware.AuthorizeJWT(jwtService)

	authController := controller.NewAuthController(authService, jwtService)
	depositController := controller.NewDepositController(depositService, jwtService)
	withdrawalController := controller.NewWithdrawalController(withdrawalService, userService, jwtService)
	userController := controller.NewUserController(userService, jwtService)
	transactionController := controller.NewTransactionController(transactionService, userService, jwtService)
	chatbotController := controller.NewChatbotController(chatbotService, jwtService)
	verificationController := controller.NewVerificationController(verificationService, jwtService)

	routes.RegisterRoutes(e, jwtService, authController)
	routes.DepositRoutes(e, depositService, depositController, jwtMiddleware)
	routes.MidtransRoutes(e, depositService, depositController, jwtMiddleware)
	routes.WithdrawalRoutes(e, withdrawalService, withdrawalController, jwtMiddleware)
	routes.UserRoutes(e, userService, userController, jwtMiddleware)
	routes.ProfileRoutes(e, userService, userController, jwtMiddleware)
	routes.TransactionRoutes(e, transactionService, transactionController, jwtMiddleware)
	routes.ImageRoutes(e, userController, jwtMiddleware)
	routes.ChatbotRoutes(e, chatbotController, jwtMiddleware)
	routes.VerificationRoutes(e, verificationService, verificationController, jwtMiddleware)

	logrus.Print(helper.GetCurrentTimeInLocation())
	e.Start(":8000")
}
