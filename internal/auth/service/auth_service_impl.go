package service

import (
	"backend_golang_codeing_test/internal/auth/model"
	userModel "backend_golang_codeing_test/internal/user/model"
	"backend_golang_codeing_test/internal/user/repository"
	"backend_golang_codeing_test/pkg/jwt"
	"backend_golang_codeing_test/pkg/redis"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type authServiceImpl struct {
	userRepo repository.UserRepository
	cache    redis.Cache
	jwt      jwt.JWTService
}

func NewAuthService(userRepo repository.UserRepository, cache redis.Cache, jwt jwt.JWTService) AuthService {
	return &authServiceImpl{
		userRepo: userRepo,
		cache:    cache,
		jwt:      jwt,
	}
}

func (s *authServiceImpl) Register(ctx context.Context, input model.RegisterRequest) error {
	existsEmail, _ := s.userRepo.ExistsByEmail(ctx, input.Email)
	if existsEmail {
		return errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user := &userModel.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return errors.New("failed to create user")
	}

	return nil
}

func (s *authServiceImpl) Login(ctx context.Context, input model.LoginRequest) (*model.LoginResponse, error) {
	var user *userModel.User
	var err error

	user, err = s.userRepo.FindByEmail(ctx, input.Email)

	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	// 🔒 เช็ค password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid password")
	}

	// 🧹 เคลียร์ session เก่า
	if err := s.clearOldSessions(ctx, user.ID); err != nil {
		return nil, fmt.Errorf("failed to clear old sessions: %w", err)
	}

	// 🔐 สร้าง token ใหม่พร้อม claims
	tokenStr, claims, err := s.jwt.GenerateTokenWithClaims(user.ID, user.Name)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// ✅ ดึง jti อย่างปลอดภัย
	jtiRaw, ok := claims["jti"]
	if !ok {
		return nil, errors.New("missing jti in token claims")
	}
	jti, ok := jtiRaw.(string)
	if !ok {
		return nil, errors.New("invalid jti format in claims")
	}

	// 💾 เก็บ token ลง Redis
	key := fmt.Sprintf("jwt:user:%s:%s", user.ID.Hex(), jti)
	if err := s.cache.Set(ctx, key, "valid", s.jwt.TTL()); err != nil {
		return nil, errors.New("failed to store session")
	}

	return &model.LoginResponse{
		Token: tokenStr,
	}, nil
}

func (s *authServiceImpl) clearOldSessions(ctx context.Context, userID primitive.ObjectID) error {
	pattern := fmt.Sprintf("jwt:user:%s:*", userID.Hex())
	keys, err := s.cache.ScanKeys(ctx, pattern)
	if err != nil {
		return err
	}

	for _, fullKey := range keys {
		err := s.cache.DeleteRaw(ctx, fullKey)
		if err != nil {
			continue
		}
	}
	return nil
}
