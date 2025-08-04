package http

import (
	"backend_golang_codeing_test/internal/user/service"

	"github.com/labstack/echo/v4"
)

func RegisterUserRoutes(g *echo.Group, s service.UserService) {
	h := NewUserHandler(s)
	g.GET("", h.GetAll)
	g.GET("/:id", h.GetByID)

	g.PUT("", h.Update)
	g.DELETE("/:id", h.Delete)

}
