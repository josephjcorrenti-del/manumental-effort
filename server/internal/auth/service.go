package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Service struct {
	repository   *Repository
	tokenManager *TokenManager
}

func NewService(repository *Repository, tokenManager *TokenManager) *Service {
	return &Service{
		repository:   repository,
		tokenManager: tokenManager,
	}
}

func (s *Service) Login(ctx context.Context, input LoginInput) (*LoginResult, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	password := strings.TrimSpace(input.Password)

	if email == "" || password == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	credential, err := s.repository.GetByEmailNormalized(ctx, email)
	if err != nil {
		return nil, err
	}

	if credential == nil {
		return nil, ErrInvalidCredentials
	}

	if err := CheckPassword(password, credential.PasswordHash); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.tokenManager.CreateToken(credential.UserID.Hex())
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		Token: token,
	}, nil
}
