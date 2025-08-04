package http

import (
	"backend_golang_codeing_test/internal/user/service"
	"backend_golang_codeing_test/internal/user/transport/http/dto"
	"backend_golang_codeing_test/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) GetAll(c echo.Context) error {
	users, err := h.service.GetAllUser((c.Request().Context()))

	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to fetch users", err)
	}
	return response.Success(c, "Fetched users successfully", users)
}

func (h *UserHandler) GetByID(c echo.Context) error {
	id := c.Param("id") // รับ ID จาก URL เช่น /api/users/:id

	user, err := h.service.FindByID(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "User not found", err)
	}

	return response.Success(c, "Fetched user successfully", user)
}

func (h *UserHandler) Update(c echo.Context) error {
	var req dto.UpdateUserRequest

	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request payload")
	}

	if err := c.Validate(&req); err != nil {
		return response.ValidationError(c, err)
	}

	if req.Email == "" && req.Name == "" {
		return response.BadRequest(c, "Nothing to update")
	}

	if err := h.service.Update(c.Request().Context(), req); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to update user", err)
	}

	return response.Success(c, "User updated successfully", nil)
}

func (h *UserHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "User ID is required")
	}

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to delete user", err)
	}

	return response.Success(c, "User deleted successfully", nil)
}
