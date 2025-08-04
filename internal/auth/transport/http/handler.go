package http

import (
	"backend_golang_codeing_test/internal/auth/model"
	"backend_golang_codeing_test/internal/auth/service"
	"backend_golang_codeing_test/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req model.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request")
	}

	if err := c.Validate(&req); err != nil {
		return response.ValidationError(c, err)
	}

	if err := h.authService.Register(c.Request().Context(), req); err != nil {
		return response.Error(c, http.StatusConflict, "Register failed", err)
	}

	return response.Created(c, "Register successful", nil)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req model.LoginRequest

	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request")
	}

	if err := c.Validate(&req); err != nil {
		return response.ValidationError(c, err)
	}

	loginResp, err := h.authService.Login(c.Request().Context(), req)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "invalid password" {
			return response.Error(c, http.StatusBadRequest, "Invalid email/username or password", err)
		}
		return response.Error(c, http.StatusInternalServerError, "Login failed", err)
	}

	return response.Success(c, "Login successful", loginResp)
}
