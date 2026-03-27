package memberships

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	RoleOwner  = "owner"
	RoleAdmin  = "admin"
	RoleMember = "member"
)

type Membership struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SpaceID   primitive.ObjectID `bson:"space_id" json:"space_id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Role      string             `bson:"role" json:"role"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
