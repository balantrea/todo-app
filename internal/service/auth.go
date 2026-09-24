package service

import (
	"crypto/sha1"
	"fmt"
	"os"
	"time"

	"github.com/balantrea/todo-app/internal/model"
	"github.com/balantrea/todo-app/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
)

var signingKey = os.Getenv("SIGNING_KEY")
var salt = os.Getenv("SALT")

const (
	tokenTTL = 12 * time.Hour
)

type AuthService struct {
	repo   repository.Authorization
	logger zerolog.Logger
}

type tokenClaims struct {
	jwt.RegisteredClaims
	UserId int `json:"user_id"`
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user model.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
		},
		user.Id,
	})

	tokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte(signingKey), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, err
	}

	return claims.UserId, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
