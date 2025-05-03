package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/taufikmulyawan/ticketing-system/config"
	"github.com/taufikmulyawan/ticketing-system/entity"
	"github.com/taufikmulyawan/ticketing-system/repository"
)

type UserService interface {
	Register(user *entity.User) error
	Login(email, password string) (string, error)
	GetUser(id uint) (*entity.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) Register(user *entity.User) error {
	// Check if email already exists
	existingUser, err := s.userRepo.FindByEmail(user.Email)
	if err == nil && existingUser != nil {
		return errors.New("email already registered")
	}

	// Set default role if not specified
	if user.Role == "" {
		user.Role = entity.RoleUser
	}

	return s.userRepo.Save(user)
}

func (s *userService) Login(email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Compare password
	err = user.ComparePassword(password)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(), // 24 hours expiration
	})

	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *userService) GetUser(id uint) (*entity.User, error) {
	return s.userRepo.FindByID(id)
} 