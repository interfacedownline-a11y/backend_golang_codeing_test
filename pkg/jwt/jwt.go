package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JWTService interface {
	GenerateToken(userID primitive.ObjectID, username string) (string, error)
	GenerateTokenWithClaims(userID primitive.ObjectID, username string) (string, jwt.MapClaims, error)
	TTL() time.Duration
	ParseToken(tokenStr string) (jwt.MapClaims, error)
}

type JWTImpl struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTManager(secret string, ttl time.Duration) JWTService {
	return &JWTImpl{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

func (j *JWTImpl) GenerateToken(userID primitive.ObjectID, username string) (string, error) {
	now := time.Now()
	jti := uuid.New().String()

	claims := jwt.MapClaims{
		"jti":      jti,
		"sub":      userID.Hex(), // แปลงเป็น string
		"username": username,
		"user_id":  userID.Hex(),
		"iat":      now.Unix(),
		"exp":      now.Add(j.ttl).Unix(),
		"nbf":      now.Unix(),
		"iss":      "golang-coding-test-api",
		"aud":      "golang-coding-test-client",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWTImpl) GenerateTokenWithClaims(userID primitive.ObjectID, username string) (string, jwt.MapClaims, error) {
	now := time.Now()
	jti := uuid.New().String()

	claims := jwt.MapClaims{
		"jti":      jti,
		"sub":      userID.Hex(),
		"username": username,
		"user_id":  userID.Hex(),
		"iat":      now.Unix(),
		"exp":      now.Add(j.ttl).Unix(),
		"nbf":      now.Unix(),
		"iss":      "golang-coding-test-api",
		"aud":      "golang-coding-test-client",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(j.secret)
	if err != nil {
		return "", nil, err
	}

	return signed, claims, nil
}

func (j *JWTImpl) TTL() time.Duration {
	return j.ttl
}

func (j *JWTImpl) ParseToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	return claims, nil
}
