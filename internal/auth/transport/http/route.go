package http

import (
	"backend_golang_codeing_test/internal/auth/service"

	"github.com/labstack/echo/v4"
)

func RegisterAuthRoutes(g *echo.Group, authService service.AuthService) {
	handler := NewAuthHandler(authService)

	g.POST("/register", handler.Register)
	g.POST("/login", handler.Login)

}
