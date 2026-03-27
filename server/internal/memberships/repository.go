package memberships

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrDuplicateMembership = errors.New("duplicate membership")

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		collection: db.Collection("memberships"),
	}
}

func (r *Repository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "space_id", Value: 1},
				{Key: "user_id", Value: 1},
			},
			Options: options.Index().
				SetUnique(true).
				SetName("space_user_unique"),
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("create membership indexes: %w", err)
	}

	return nil
}

func (r *Repository) Create(ctx context.Context, membership *Membership) error {
	result, err := r.collection.InsertOne(ctx, membership)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ErrDuplicateMembership
		}
		return fmt.Errorf("insert membership: %w", err)
	}

	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("inserted membership id is not an ObjectID")
	}

	membership.ID = objectID
	return nil
}

func (r *Repository) Exists(ctx context.Context, spaceID primitive.ObjectID, userID primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"space_id": spaceID,
		"user_id":  userID,
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("count memberships: %w", err)
	}

	return count > 0, nil
}
