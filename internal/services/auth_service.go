package services

import (
	"errors"

	"github.com/kazantsev/mentorship-backend/internal/config"
	"github.com/kazantsev/mentorship-backend/internal/models"
	"github.com/kazantsev/mentorship-backend/internal/repositories"
	"github.com/kazantsev/mentorship-backend/internal/utils"
)

// AuthService handles user authentication and registration
type AuthService struct {
	userRepo *repositories.UserRepository
	cfg      *config.Config
}

// NewAuthService returns a new AuthService instance
func NewAuthService(userRepo *repositories.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

// RegisterRequest contains the data required to register a new user
type RegisterRequest struct {
	Login            string
	Password         string
	DisplayName      string
	TelegramUsername string
	Roles            []models.Role
}

// Register creates a new user account with the provided details
func (s *AuthService) Register(req RegisterRequest) (*models.User, error) {
	existing, _, _ := s.userRepo.FindByLogin(req.Login)
	if existing != nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Login:            req.Login,
		PasswordHash:     hashedPassword,
		DisplayName:      req.DisplayName,
		TelegramUsername: req.TelegramUsername,
		IsProfilePrivate: false,
		IsDeleted:        false,
	}

	if err := s.userRepo.Create(user, req.Roles); err != nil {
		return nil, err
	}
	return user, nil
}

// LoginRequest contains credentials for authentication
type LoginRequest struct {
	Login    string
	Password string
}

// LoginResponse contains the result of a successful authentication
type LoginResponse struct {
	Token string
	User  *models.User
	Roles []models.Role
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(req LoginRequest) (*LoginResponse, error) {
	user, roles, err := s.userRepo.FindByLogin(req.Login)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}
	roleStrs := make([]string, len(roles))
	for i, r := range roles {
		roleStrs[i] = string(r)
	}
	token, err := utils.GenerateJWT(user.ID, user.Login, roleStrs, s.cfg.JWT.Secret, s.cfg.JWT.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &LoginResponse{
		Token: token,
		User:  user,
		Roles: roles,
	}, nil
}
