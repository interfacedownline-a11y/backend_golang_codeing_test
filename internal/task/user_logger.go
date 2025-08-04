package task

import (
	"backend_golang_codeing_test/internal/user/repository"
	"backend_golang_codeing_test/pkg/logger"
	"context"
	"time"
)

func StartUserLoggerTask(repo repository.UserRepository) {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for range ticker.C {
			users, err := repo.GetAllUser(context.Background())
			if err != nil {
				logger.Error("Failed to fetch user count", "error", err)
				continue
			}
			logger.Info("User count", "count", len(users))
		}
	}()
}
