package auth

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Credential struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          primitive.ObjectID `bson:"user_id" json:"user_id"`
	EmailNormalized string             `bson:"email_normalized" json:"-"`
	PasswordHash    string             `bson:"password_hash" json:"-"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResult struct {
	Token string `json:"token"`
}
