package response

import (
	"net/http"
	"strings"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type BaseResponse struct {
	Status  string      `json:"status"`          // "success" or "error"
	Code    int         `json:"code"`            // HTTP status code
	Message string      `json:"message"`         // Human-readable message
	Data    interface{} `json:"data,omitempty"`  // Only for success
	Error   string      `json:"error,omitempty"` // Only for failure
}

func Success(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusOK, BaseResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	})
}

func Created(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusCreated, BaseResponse{
		Status:  "success",
		Code:    http.StatusCreated,
		Message: message,
		Data:    data,
	})
}

func Error(c echo.Context, code int, message string, err error) error {
	return c.JSON(code, BaseResponse{
		Status:  "error",
		Code:    code,
		Message: message,
		Error:   err.Error(),
	})
}

func ValidationError(c echo.Context, err error) error {
	fields := make(map[string]string)

	if verrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range verrs {
			jsonKey := strings.ToLower(e.Field())
			label := jsonKey

			var message string
			switch e.Tag() {
			case "required":
				message = label + " is required"
			case "min":
				message = label + " must be at least " + e.Param() + " characters"
			case "max":
				message = label + " must be at most " + e.Param() + " characters"
			default:
				message = label + " is invalid"
			}

			fields[jsonKey] = message
		}
	}

	return c.JSON(http.StatusUnprocessableEntity, echo.Map{
		"status":  "error",
		"code":    http.StatusUnprocessableEntity,
		"message": "Validation failed",
		"data":    fields,
	})
}

func BadRequest(c echo.Context, message string) error {
	return c.JSON(http.StatusBadRequest, BaseResponse{
		Status:  "error",
		Code:    http.StatusBadRequest,
		Message: message,
	})
}

func Unauthorized(c echo.Context, message string) error {
	var errStr string

	return c.JSON(http.StatusUnauthorized, BaseResponse{
		Status:  "error",
		Code:    http.StatusUnauthorized,
		Message: message,
		Error:   errStr,
	})
}
