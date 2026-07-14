package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type TokenGenerator interface {
	GenerateAccessToken(userID uuid.UUID) (string, error)
}

type AuthService interface {
	Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error)
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
}

type Service struct {
	repo  UserRepository
	token TokenGenerator
}

func NewService(repo UserRepository, tokenGenerator TokenGenerator) *Service {
	return &Service{
		repo:  repo,
		token: tokenGenerator,
	}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &User{
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return &RegisterResponse{
		User: UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("compare password: %w", err)
	}

	accessToken, err := s.token.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	return &LoginResponse{
		AccessToken: accessToken,
		User: UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}
