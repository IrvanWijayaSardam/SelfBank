package main

import (
	"github.com/IrvanWijayaSardam/SelfBank/config"
	"github.com/IrvanWijayaSardam/SelfBank/controller"
	"github.com/IrvanWijayaSardam/SelfBank/repository"
	"github.com/IrvanWijayaSardam/SelfBank/routes"
	"github.com/IrvanWijayaSardam/SelfBank/service"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

var (
	db *gorm.DB = config.SetupDatabaseConnection()

	userRepository repository.UserRepository = repository.NewUserRepository(db)

	authService service.AuthService = service.NewAuthService(userRepository)
	jwtService  service.JWTService  = service.NewJWTService()
)

func main() {
	e := echo.New()
	e.Debug = true

	authController := controller.NewAuthController(authService, jwtService)
	routes.RegisterRoutes(e, jwtService, authController)

	e.Start(":8001")
}
