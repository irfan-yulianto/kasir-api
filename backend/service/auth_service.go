package service

import (
	"errors"
	"time"

	"kasir-api/model"
	"kasir-api/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists         = errors.New("username already exists")
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrInvalidRole        = errors.New("role must be 'admin' or 'kasir'")
)

type AuthService interface {
	Register(req model.RegisterRequest) (*model.AuthResponse, error)
	Login(req model.LoginRequest) (*model.AuthResponse, error)
}

type authService struct {
	repo      repository.UserRepository
	jwtSecret []byte
	jwtExpiry time.Duration
}

func NewAuthService(repo repository.UserRepository, jwtSecret string, jwtExpiryHours int) AuthService {
	return &authService{
		repo:      repo,
		jwtSecret: []byte(jwtSecret),
		jwtExpiry: time.Duration(jwtExpiryHours) * time.Hour,
	}
}

func (s *authService) Register(req model.RegisterRequest) (*model.AuthResponse, error) {
	if req.Role != "admin" && req.Role != "kasir" {
		return nil, ErrInvalidRole
	}

	existing, err := s.repo.GetByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUserExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Role:     req.Role,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{Token: token, User: *user}, nil
}

func (s *authService) Login(req model.LoginRequest) (*model.AuthResponse, error) {
	user, err := s.repo.GetByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{Token: token, User: *user}, nil
}

func (s *authService) generateToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(s.jwtExpiry).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
