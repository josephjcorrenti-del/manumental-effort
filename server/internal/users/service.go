package users

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateUser(ctx context.Context, input CreateUserInput) (*User, error) {
	if err := validateCreateUserInput(input); err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	user := &User{
		Username:           strings.TrimSpace(input.Username),
		UsernameNormalized: normalizeUsername(input.Username),
		DisplayName:        strings.TrimSpace(input.DisplayName),
		Email:              strings.TrimSpace(input.Email),
		EmailNormalized:    normalizeEmail(input.Email),
		IsActive:           true,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := s.repository.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetUserByID(ctx context.Context, id string) (*User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user id")
	}

	return s.repository.GetByID(ctx, objectID)
}
