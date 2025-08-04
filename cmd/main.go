package main

import (
	"backend_golang_codeing_test/config"
	"backend_golang_codeing_test/internal/server"
	"backend_golang_codeing_test/internal/task"
	"backend_golang_codeing_test/migrations"
	"backend_golang_codeing_test/pkg/database"
	"backend_golang_codeing_test/pkg/jwt"
	"backend_golang_codeing_test/pkg/logger"
	"backend_golang_codeing_test/pkg/redis"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"

	userRepository "backend_golang_codeing_test/internal/user/repository"
)

func main() {
	fmt.Println("main init!!")
	initLogger()

	cfg := loadConfig()

	db := setupDatabase(cfg)

	redisManager := setupRedis(cfg)

	writeCache := redis.NewCache(redisManager.GetWrite(), cfg.Redis.NameSpace)
	readCache := redis.NewCache(redisManager.GetRead(), cfg.Redis.NameSpace)

	jwtManager := setupJWT(cfg)

	setTask(db)
	runMigrations(db)

	e := startServer(db, cfg.Server.Port, writeCache, readCache, jwtManager)

	go func() {
		addr := fmt.Sprintf(":%d", cfg.Server.Port)
		if err := e.Start(addr); err != nil {
			log.Printf("Shutting down server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal("Failed to gracefully shutdown:", err)
	}
}

func initLogger() {
	logger.Init()
}

func loadConfig() *config.Config {
	return config.LoadConfig()
}

func setupDatabase(cfg *config.Config) *mongo.Database {
	return database.NewMongoDatabase(&cfg.Database)
}

func setupRedis(cfg *config.Config) redis.Redis {
	return redis.NewRedisManager(&cfg.Redis)
}

func setupJWT(cfg *config.Config) jwt.JWTService {
	return jwt.NewJWTManager(cfg.Jwt.Secret, cfg.Jwt.TTL)
}

func setTask(db *mongo.Database) {
	userRepo := userRepository.NewUserRepository(db)
	task.StartUserLoggerTask(userRepo)
}

func runMigrations(db *mongo.Database) {
	migrations.AutoMigrate(db)
	migrations.Seed(db)
}

func startServer(db *mongo.Database, port int, writeCache redis.Cache, readCache redis.Cache, jwtManager jwt.JWTService) *echo.Echo {
	e := server.InitRouter(db, writeCache, readCache, jwtManager)
	addr := fmt.Sprintf(":%d", port)
	e.Logger.Fatal(e.Start(addr))

	return e
}
