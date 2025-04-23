package service

import (
	"context"
	"film-library/internal/model"
	"film-library/internal/repository"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const tokenTTL = 24 * time.Hour

type AuthService struct {
	repo      repository.AuthRepository
	secretKey []byte
}

func NewAuthService(repo repository.AuthRepository, secret string) *AuthService {
	return &AuthService{
		repo:      repo,
		secretKey: []byte(secret),
	}
}

func (s *AuthService) CreateUser(ctx context.Context, user model.User) (string, error) {
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = hashedPassword

	if err := s.repo.CreateUser(ctx, &user); err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id":  user.ID,
			"username": user.Username,
			"exp":      time.Now().Add(tokenTTL).Unix(),
		})

	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, nil
}

func (s *AuthService) VerifyUser(ctx context.Context, username, password string) (string, *model.User, error) {
	user, err := s.repo.VerifyUser(ctx, username)
	if err != nil {
		return "", nil, fmt.Errorf("failed to verify user: %w", err)
	}

	if !checkPasswordHash(password, user.Password) {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id":  user.ID,
			"username": user.Username,
			"exp":      time.Now().Add(24 * time.Hour).Unix(),
		})

	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, user, nil
}

//

func (s *AuthService) VerifyToken(tokenString string) (*model.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*model.TokenClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

//

func hashPassword(password string) (string, error) {
	fmt.Println("Password to hash:", fmt.Sprintf("%q", password)) // добавил %q чтобы увидеть пробелы
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func checkPasswordHash(password, hash string) bool {
	fmt.Println("Password to check:", password)
	fmt.Println("Hash to compare:", hash)
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println("Password comparison failed:", err)
	}
	return err == nil
}
