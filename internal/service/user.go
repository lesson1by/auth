package service

import (
	"authProject/internal/models"
	"authProject/internal/store"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type UserService interface {
	ValidateCredentials(username, password string) (bool, error)
	Register(username, password string) error
	GenerateToken(username string) (string, error)
	RefreshToken(token string) (string, error)
}
type userService struct {
	UserStore store.UserStore
	Config    models.Config
}

func NewUserService(config *models.Config) UserService {
	if config == nil {
		return &userService{
			UserStore: store.NewInMemoryStore(),
			Config:    models.Config{},
		}
	}
	return &userService{
		UserStore: store.NewInMemoryStore(),
		Config:    *config,
	}
}

func (s *userService) ValidateCredentials(username, password string) (bool, error) {
	user, err := s.UserStore.Get(username)
	if err != nil {
		return false, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return false, nil
	}
	return true, nil
}

func (s *userService) Register(username, password string) error {
	if username == "" {
		return errors.New("username is required")
	}
	if password == "" {
		return errors.New("password is required")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.UserStore.Create(models.User{
		Username: username,
		Password: string(hash),
	})
}

func (s *userService) GenerateToken(username string) (string, error) {
	secret := []byte(s.Config.JWT.Secret)
	if len(secret) < 32 {
		return "", errors.New("secret must be at least 32 bytes")
	}
	expiration := time.Now().Add(time.Duration(s.Config.JWT.ExpirationMinutes) * time.Minute)
	claims := jwt.RegisteredClaims{
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(expiration),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		ID:        uuid.New().String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func (s *userService) RefreshToken(token string) (string, error) { // узнать как это работает
	secret := []byte(s.Config.JWT.Secret)
	if len(secret) < 32 {
		return "", errors.New("secret must be at least 32 bytes")
	}

	token = strings.TrimSpace(token)
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimSpace(strings.TrimPrefix(token, "Bearer "))
	}

	parsedClaims := &jwt.RegisteredClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, parsedClaims, func(t *jwt.Token) (any, error) {
		return secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return "", err
	}
	if !parsedToken.Valid {
		return "", errors.New("invalid token")
	}
	if parsedClaims.Subject == "" {
		return "", errors.New("token subject is empty")
	}

	newExpiration := time.Now().Add(time.Duration(s.Config.JWT.ExpirationMinutes) * time.Minute)
	newClaims := jwt.RegisteredClaims{
		Subject:   parsedClaims.Subject,
		ExpiresAt: jwt.NewNumericDate(newExpiration),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		ID:        uuid.New().String(),
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	return newToken.SignedString(secret)
}
