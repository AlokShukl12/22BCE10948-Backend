package services

import (
	"errors"
	"filesharing/auth"
	"filesharing/models"
	"filesharing/repositories"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    *repositories.UserRepository
	redisClient *redis.Client
}

func NewAuthService(userRepo *repositories.UserRepository, redisClient *redis.Client) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		redisClient: redisClient,
	}
}

func (s *AuthService) Register(register *models.UserRegister) (*models.UserResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.FindByEmail(register.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		ID:        uuid.New(),
		Email:     register.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *AuthService) Login(login *models.UserLogin) (string, error) {
	// Find user
	user, err := s.userRepo.FindByEmail(login.Email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		return "", err
	}

	return token, nil
} 