package spaces

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

func (s *Service) CreateSpace(ctx context.Context, userID string, input CreateSpaceInput) (*Space, error) {
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

	now := time.Now().UTC()

	space := &Space{
		Name:         name,
		Slug:         generateSlug(name),
		Description:  strings.TrimSpace(input.Description),
		Visibility:   input.Visibility,
		Discoverable: input.Discoverable,
		CreatedBy:    objectUserID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repository.Create(ctx, space); err != nil {
		return nil, err
	}

	membership := &memberships.Membership{
		SpaceID:   space.ID,
		UserID:    objectUserID,
		Role:      memberships.RoleOwner,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.membershipRepository.Create(ctx, membership); err != nil {
		return nil, fmt.Errorf("create owner membership: %w", err)
	}

	return space, nil
}

func (s *Service) JoinSpace(ctx context.Context, userID string, spaceID string) error {
	objectUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user id")
	}

	objectSpaceID, err := primitive.ObjectIDFromHex(spaceID)
	if err != nil {
		return fmt.Errorf("invalid space id")
	}

	space, err := s.repository.GetByID(ctx, objectSpaceID)
	if err != nil {
		return err
	}
	if space == nil {
		return fmt.Errorf("space not found")
	}

	if space.Visibility != "public" {
		return fmt.Errorf("space is not joinable")
	}

	now := time.Now().UTC()

	membership := &memberships.Membership{
		SpaceID:   objectSpaceID,
		UserID:    objectUserID,
		Role:      memberships.RoleMember,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.membershipRepository.Create(ctx, membership); err != nil {
		return err
	}

	return nil
}

func generateSlug(name string) string {
	slug := strings.ToLower(strings.TrimSpace(name))
	slug = strings.ReplaceAll(slug, " ", "-")
	return slug
}
