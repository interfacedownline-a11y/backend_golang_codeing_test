package middleware

import (
	"backend_golang_codeing_test/pkg/logger"
	"backend_golang_codeing_test/pkg/response"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RecoverMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			// if err := recover(); err != nil {
			// 	fmt.Printf("Recovered from panic: %v\n", err)

			// 	_ = c.JSON(http.StatusInternalServerError, map[string]string{
			// 		"message": "Internal Server Error",
			// 	})
			// }

			if err := recover(); err != nil {
				var e error
				if casted, ok := err.(error); ok {
					e = casted
				} else {
					e = fmt.Errorf("%v", err)
				}
				logger.Error("Internal Server Error : Panic", "code", e)
				response.Error(c, http.StatusInternalServerError, "Internal Server Error : Panic", e)
			}
		}()
		return next(c)
	}
}
