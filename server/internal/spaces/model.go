package spaces

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Space struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string             `bson:"name" json:"name"`
	Slug         string             `bson:"slug" json:"slug"`
	Description  string             `bson:"description" json:"description"`
	Visibility   string             `bson:"visibility" json:"visibility"`
	Discoverable bool               `bson:"discoverable" json:"discoverable"`
	CreatedBy    primitive.ObjectID `bson:"created_by" json:"created_by"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

type CreateSpaceInput struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Visibility   string `json:"visibility"`
	Discoverable bool   `json:"discoverable"`
}
