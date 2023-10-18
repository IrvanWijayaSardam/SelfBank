package middleware

import (
	"log"
	"net/http"

	"github.com/IrvanWijayaSardam/SelfBank/service"

	"github.com/IrvanWijayaSardam/SelfBank/helper"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func AuthorizeJWT(jwtService service.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				response := helper.BuildErrorResponse("No Token Found !")
				return c.JSON(http.StatusBadRequest, response)
			}

			tokenString := authHeader
			token, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				log.Println(err)
				response := helper.BuildErrorResponse("Token is not valid -" + err.Error())
				return c.JSON(http.StatusUnauthorized, response)
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				c.Set("user", claims)
				return next(c)
			}

			response := helper.BuildErrorResponse("Invalid Token Claims")
			return c.JSON(http.StatusUnauthorized, response)
		}
	}
}
