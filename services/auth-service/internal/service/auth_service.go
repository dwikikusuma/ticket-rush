package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/dwikikusuma/ticket-rush/services/auth-service/internal/config"
	"github.com/dwikikusuma/ticket-rush/services/auth-service/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo domain.UserRepo
}

func NewAuthService(repo domain.UserRepo) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password string) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to generate hash password: %s", err)
		return err
	}

	err = s.repo.CreateUser(ctx, email, string(hashPassword))
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		log.Printf("Failed to get user by email: %s", err)
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.GSetPassword()), []byte(password))
	if err != nil {
		log.Printf("invalid username password")
		return "", errors.New("invalid username or password")
	}

	return s.generateToken(user.ID, user.Email)
}

func (s *AuthService) generateToken(userID, email string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"iss":   "ticket-rush-auth",
		"exp":   time.Now().Add(3 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(config.JWTSecrete)
	if err != nil {
		log.Printf("failed to create token: %s", err)
		return "", err
	}

	return signed, nil
}
