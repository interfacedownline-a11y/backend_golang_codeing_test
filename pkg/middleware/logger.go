package middleware

import (
	"time"

	"backend_golang_codeing_test/pkg/logger"

	"github.com/labstack/echo/v4"
)

func LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		stop := time.Now()

		status := c.Response().Status
		latency := stop.Sub(start).String()

		logFields := []interface{}{
			"method", c.Request().Method,
			"path", c.Path(),
			"ip", c.RealIP(),
			"status", status,
			"latency", latency,
			"user_agent", c.Request().UserAgent(),
		}

		if err != nil {
			logFields = append(logFields, "error", err.Error())
		}

		if err != nil || status < 200 || status >= 300 {
			logger.Error("HTTP Request", logFields...)
		} else {
			logger.Info("HTTP Request", logFields...)
		}

		return err
	}
}
