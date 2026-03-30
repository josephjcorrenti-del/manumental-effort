package messages

import (
	"context"
	"errors"
	"strings"
	"time"

	"manumental-effort/server/internal/channels"
	"manumental-effort/server/internal/memberships"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrMessageBodyRequired = errors.New("message body is required")
	ErrMessageBodyTooLong  = errors.New("message body is too long")
	ErrMembershipRequired  = errors.New("membership required")
)

const DefaultMessageListLimit = 50
const MaxMessageBodyLength = 2000

type MessageBroadcaster interface {
	BroadcastMessageCreated(message *Message)
}

type MessageRepository interface {
	EnsureIndexes(ctx context.Context) error
	Create(ctx context.Context, message *Message) error
	ListByChannelBefore(
		ctx context.Context,
		channelID primitive.ObjectID,
		beforeID *primitive.ObjectID,
		limit int,
	) ([]Message, error)
}

//type Service struct {
//	messageRepository    MessageRepository
//	channelRepository    *channels.Repository
//	membershipRepository *memberships.Repository
//}

type Service struct {
	messageRepository    MessageRepository
	channelRepository    *channels.Repository
	membershipRepository *memberships.Repository
	broadcaster          MessageBroadcaster
}

//func NewService(
//	messageRepository MessageRepository,
//	channelRepository *channels.Repository,
//	membershipRepository *memberships.Repository,
//) *Service {
//	return &Service{
//		messageRepository:    messageRepository,
//		channelRepository:    channelRepository,
//		membershipRepository: membershipRepository,
//	}
//}

func NewService(
	messageRepository MessageRepository,
	channelRepository *channels.Repository,
	membershipRepository *memberships.Repository,
	broadcaster MessageBroadcaster,
) *Service {
	return &Service{
		messageRepository:    messageRepository,
		channelRepository:    channelRepository,
		membershipRepository: membershipRepository,
		broadcaster:          broadcaster,
	}
}

func (s *Service) CreateMessage(
	ctx context.Context,
	channelID primitive.ObjectID,
	userID primitive.ObjectID,
	body string,
) (*Message, error) {
	body = strings.TrimSpace(body)
	if body == "" {
		return nil, ErrMessageBodyRequired
	}
	if len(body) > MaxMessageBodyLength {
		return nil, ErrMessageBodyTooLong
	}

	channel, err := s.channelRepository.GetByID(ctx, channelID)
	if err != nil {
		return nil, err
	}

	isMember, err := s.membershipRepository.Exists(ctx, channel.SpaceID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrMembershipRequired
	}

	now := time.Now().UTC()

	message := &Message{
		ChannelID: channelID,
		UserID:    userID,
		Body:      body,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.messageRepository.Create(ctx, message); err != nil {
		return nil, err
	}

	if s.broadcaster != nil {
		s.broadcaster.BroadcastMessageCreated(message)
	}

	return message, nil
}

func (s *Service) ListMessages(
	ctx context.Context,
	channelID primitive.ObjectID,
	userID primitive.ObjectID,
	beforeID *primitive.ObjectID,
	limit int,
) (*ListMessagesResult, error) {
	if limit <= 0 {
		limit = DefaultMessageListLimit
	}
	if limit > DefaultMessageListLimit {
		limit = DefaultMessageListLimit
	}

	channel, err := s.channelRepository.GetByID(ctx, channelID)
	if err != nil {
		return nil, err
	}

	isMember, err := s.membershipRepository.Exists(ctx, channel.SpaceID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrMembershipRequired
	}

	rows, err := s.messageRepository.ListByChannelBefore(ctx, channelID, beforeID, limit+1)
	if err != nil {
		return nil, err
	}

	var nextCursor *string
	if len(rows) > limit {
		rows = rows[:limit]
		cursor := rows[len(rows)-1].ID.Hex()
		nextCursor = &cursor
	}

	reverseMessages(rows)

	return &ListMessagesResult{
		Items:      rows,
		NextCursor: nextCursor,
	}, nil
}

func reverseMessages(items []Message) {
	for left, right := 0, len(items)-1; left < right; left, right = left+1, right-1 {
		items[left], items[right] = items[right], items[left]
	}
}
