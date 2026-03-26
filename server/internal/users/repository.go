package users

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrDuplicateUsername = errors.New("duplicate username")
var ErrDuplicateEmail = errors.New("duplicate email")

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		collection: db.Collection("users"),
	}
}

func (r *Repository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "username_normalized", Value: 1}},
			Options: options.Index().
				SetUnique(true).
				SetName("username_normalized_unique"),
		},
		{
			Keys: bson.D{{Key: "email_normalized", Value: 1}},
			Options: options.Index().
				SetUnique(true).
				SetName("email_normalized_unique"),
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("create user indexes: %w", err)
	}

	return nil
}

func (r *Repository) Create(ctx context.Context, user *User) error {
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			switch {
			case r.collectionHasValue(ctx, "username_normalized", user.UsernameNormalized):
				return ErrDuplicateUsername
			case r.collectionHasValue(ctx, "email_normalized", user.EmailNormalized):
				return ErrDuplicateEmail
			default:
				return fmt.Errorf("insert user duplicate key: %w", err)
			}
		}
		return fmt.Errorf("insert user: %w", err)
	}

	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("inserted id is not an ObjectID")
	}

	user.ID = objectID
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id primitive.ObjectID) (*User, error) {
	var user User

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}

	return &user, nil
}

func (r *Repository) collectionHasValue(ctx context.Context, field string, value string) bool {
	filter := bson.M{field: value}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false
	}

	return count > 0
}
