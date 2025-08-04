package migrations

import (
	"backend_golang_codeing_test/internal/user/model"
	userRepo "backend_golang_codeing_test/internal/user/repository"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func Seed(db *mongo.Database) {
	ctx := context.Background()
	seedUsers(ctx, db)
}

func seedUsers(ctx context.Context, db *mongo.Database) {
	userRepo := userRepo.NewUserRepository(db)

	collection := db.Collection("users")
	count, err := collection.CountDocuments(ctx, bson.M{})

	if err != nil {
		log.Printf("Failed to count users: %v", err)
		return
	}
	if count > 0 {
		log.Println("Users already seeded, skipping...")
		return
	}

	users := []struct {
		Username string
		Email    string
		Password string
	}{
		{"admin", "admin@example.com", "123456789"},
		{"demo", "demo@example.com", "123456789"},
	}

	for _, u := range users {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash password for %s: %v", u.Username, err)
			continue
		}

		user := &model.User{
			Name:     u.Username,
			Email:    u.Email,
			Password: string(hashedPassword),
		}

		if err := userRepo.Create(ctx, user); err != nil {
			log.Printf("Failed to seed user %s: %v", u.Username, err)
		} else {
			log.Printf("Seeded user: %s", u.Username)
		}
	}
}
