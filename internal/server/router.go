package server

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"

	"backend_golang_codeing_test/pkg/jwt"
	"backend_golang_codeing_test/pkg/middleware"
	"backend_golang_codeing_test/pkg/redis"
	"backend_golang_codeing_test/pkg/response"

	"backend_golang_codeing_test/internal/auth/service"
	userService "backend_golang_codeing_test/internal/user/service"

	userRepository "backend_golang_codeing_test/internal/user/repository"

	authHttp "backend_golang_codeing_test/internal/auth/transport/http"
	userHttp "backend_golang_codeing_test/internal/user/transport/http"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func InitRouter(db *mongo.Database, writeCache redis.Cache, readCache redis.Cache, jwt jwt.JWTService) *echo.Echo {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.RecoverMiddleware)
	e.Use(middleware.LoggerMiddleware)

	e.GET("/health", func(c echo.Context) error {
		return response.Success(c, "OK", nil)
	})

	//repository
	userRepo := userRepository.NewUserRepository(db)

	//service
	authService := service.NewAuthService(userRepo, writeCache, jwt)
	userService := userService.NewUserService(userRepo)

	api := e.Group("/api")

	registerPublicRoutes(api, authService)

	registerProtectedRoutes(api, jwt, readCache, userService)

	return e
}

func registerPublicRoutes(api *echo.Group, authSvc service.AuthService) {
	authGroup := api.Group("/auth")
	authHttp.RegisterAuthRoutes(authGroup, authSvc)
}

func registerProtectedRoutes(api *echo.Group, jwt jwt.JWTService, readCache redis.Cache, userSvc userService.UserService) {
	protected := api.Group("")
	protected.Use(middleware.JWTAuthMiddleware(jwt, readCache))

	userGroup := protected.Group("/users")
	userHttp.RegisterUserRoutes(userGroup, userSvc)

}
