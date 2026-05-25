package services

import (
	"errors"

	"github.com/kazantsev/mentorship-backend/internal/config"
	"github.com/kazantsev/mentorship-backend/internal/models"
	"github.com/kazantsev/mentorship-backend/internal/repositories"
	"github.com/kazantsev/mentorship-backend/internal/utils"
)

type AuthService struct {
	userRepo *repositories.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo *repositories.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

type RegisterRequest struct {
	Login            string
	Password         string
	DisplayName      string
	TelegramUsername string
	Roles            []models.Role
}

func (s *AuthService) Register(req RegisterRequest) (*models.User, error) {
	// Проверяем, существует ли пользователь
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

type LoginRequest struct {
	Login    string
	Password string
}

type LoginResponse struct {
	Token string
	User  *models.User
	Roles []models.Role
}

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
