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

	userRepository       repository.UserRepository       = repository.NewUserRepository(db)
	depositRepository    repository.DepositRepository    = repository.NewDepositRepository(db)
	withdrawalRepository repository.WithdrawalRepository = repository.NewWithdrawalRepository(db)

	authService       service.AuthService       = service.NewAuthService(userRepository)
	jwtService        service.JWTService        = service.NewJWTService()
	depositService    service.DepositService    = service.NewDepositService(depositRepository)
	withdrawalService service.WithdrawalService = service.NewWithdrawalService(withdrawalRepository)
)

func main() {
	e := echo.New()
	e.Debug = true
	jwtMiddleware := middleware.AuthorizeJWT(jwtService)

	authController := controller.NewAuthController(authService, jwtService)
	depositController := controller.NewDepositController(depositService, jwtService)
	withdrawalController := controller.NewWithdrawalController(withdrawalService, jwtService)

	routes.RegisterRoutes(e, jwtService, authController)
	routes.DepositRoutes(e, depositService, depositController, jwtMiddleware)
	routes.MidtransRoutes(e, depositService, depositController, jwtMiddleware)
	routes.WithdrawalRoutes(e, withdrawalService, withdrawalController, jwtMiddleware)

	e.Start(":8001")
}
