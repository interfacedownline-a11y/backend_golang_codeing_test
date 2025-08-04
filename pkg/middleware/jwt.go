package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"

	jwtutil "backend_golang_codeing_test/pkg/jwt"
	"backend_golang_codeing_test/pkg/redis"
	"backend_golang_codeing_test/pkg/response"
)

type AuthClaims struct {
	UserID   primitive.ObjectID
	Username string
	JTI      string
}

func JWTAuthMiddleware(jwtManager jwtutil.JWTService, readCache redis.Cache) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				return response.Unauthorized(c, "Missing or invalid Authorization header")
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := jwtManager.ParseToken(tokenStr)
			if err != nil {
				return response.Error(c, http.StatusUnauthorized, "Invalid token", err)
			}

			// ✅ Extract jti
			jti, ok := claims["jti"].(string)
			if !ok || jti == "" {
				return response.Unauthorized(c, "Missing or invalid jti")
			}

			// ✅ Extract user info (ObjectID)
			userIDStr, ok1 := claims["user_id"].(string)
			username, ok2 := claims["username"].(string)
			if !ok1 || !ok2 {
				return response.Unauthorized(c, "Invalid token claims")
			}

			userID, err := primitive.ObjectIDFromHex(userIDStr)
			if err != nil {
				return response.Unauthorized(c, "Invalid ObjectID format in token")
			}

			// ⏱️ Validate session key in Redis
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			sessionKey := fmt.Sprintf("jwt:user:%s:%s", userID.Hex(), jti)
			if _, err := readCache.Get(ctx, sessionKey); err != nil {
				return response.Unauthorized(c, "Token revoked or expired")
			}

			// 🧠 Set user context
			c.Set("user", AuthClaims{
				UserID:   userID,
				Username: username,
				JTI:      jti,
			})

			return next(c)
		}
	}
}
