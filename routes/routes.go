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
