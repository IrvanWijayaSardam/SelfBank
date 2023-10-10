package main

import (
	"github.com/IrvanWijayaSardam/SelfBank/config"
	"github.com/IrvanWijayaSardam/SelfBank/controller"
	"github.com/IrvanWijayaSardam/SelfBank/middleware"
	"github.com/IrvanWijayaSardam/SelfBank/repository"
	"github.com/IrvanWijayaSardam/SelfBank/routes"
	"github.com/IrvanWijayaSardam/SelfBank/service"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

var (
	db *gorm.DB = config.SetupDatabaseConnection()

	userRepository        repository.UserRepository        = repository.NewUserRepository(db)
	transactioNRepository repository.TransactionRepository = repository.NewTransactionRepository(db)

	authService       service.AuthService        = service.NewAuthService(userRepository)
	jwtService        service.JWTService         = service.NewJWTService()
	transacionService service.TransactionService = service.NewTransactionService(transactioNRepository)
)

func main() {
	e := echo.New()
	e.Debug = true
	jwtMiddleware := middleware.AuthorizeJWT(jwtService)

	authController := controller.NewAuthController(authService, jwtService)
	transactionController := controller.NewTransactionController(transacionService, jwtService)

	routes.RegisterRoutes(e, jwtService, authController)
	routes.TransactionRoutes(e, transacionService, transactionController, jwtMiddleware)

	e.Start(":8001")
}
