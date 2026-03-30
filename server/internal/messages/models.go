package messages

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ChannelID primitive.ObjectID `bson:"channel_id" json:"channel_id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Body      string             `bson:"body" json:"body"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type CreateMessageInput struct {
	Body string `json:"body"`
}

type ListMessagesResult struct {
	Items      []Message `json:"items"`
	NextCursor *string   `json:"next_cursor"`
}
