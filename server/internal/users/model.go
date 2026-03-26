package users

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username           string             `bson:"username" json:"username"`
	UsernameNormalized string             `bson:"username_normalized" json:"-"`
	DisplayName        string             `bson:"display_name" json:"display_name"`
	Email              string             `bson:"email" json:"email"`
	EmailNormalized    string             `bson:"email_normalized" json:"-"`
	IsActive           bool               `bson:"is_active" json:"is_active"`
	CreatedAt          time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt          time.Time          `bson:"updated_at" json:"updated_at"`
}

type CreateUserInput struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}
