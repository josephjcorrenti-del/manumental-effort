package channels

import (
	"context"
	"fmt"
	"strings"
	"time"

	"manumental-effort/server/internal/memberships"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repository           *Repository
	membershipRepository *memberships.Repository
}

func NewService(repository *Repository, membershipRepository *memberships.Repository) *Service {
	return &Service{
		repository:           repository,
		membershipRepository: membershipRepository,
	}
}

func (s *Service) CreateChannel(ctx context.Context, userID string, spaceID string, input CreateChannelInput) (*Channel, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}

	if input.Visibility != "public" && input.Visibility != "private" {
		return nil, fmt.Errorf("invalid visibility")
	}

	objectUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id")
	}

	objectSpaceID, err := primitive.ObjectIDFromHex(spaceID)
	if err != nil {
		return nil, fmt.Errorf("invalid space id")
	}

	isMember, err := s.membershipRepository.Exists(ctx, objectSpaceID, objectUserID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, fmt.Errorf("membership required")
	}

	now := time.Now().UTC()

	channel := &Channel{
		SpaceID:     objectSpaceID,
		Name:        name,
		Slug:        generateSlug(name),
		Description: strings.TrimSpace(input.Description),
		Visibility:  input.Visibility,
		CreatedBy:   objectUserID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repository.Create(ctx, channel); err != nil {
		return nil, err
	}

	return channel, nil
}

func (s *Service) ListChannels(ctx context.Context, userID string, spaceID string) ([]Channel, error) {
	objectUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id")
	}

	objectSpaceID, err := primitive.ObjectIDFromHex(spaceID)
	if err != nil {
		return nil, fmt.Errorf("invalid space id")
	}

	isMember, err := s.membershipRepository.Exists(ctx, objectSpaceID, objectUserID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, fmt.Errorf("membership required")
	}

	return s.repository.ListBySpaceID(ctx, objectSpaceID)
}

func generateSlug(name string) string {
	slug := strings.ToLower(strings.TrimSpace(name))
	slug = strings.ReplaceAll(slug, " ", "-")
	return slug
}
