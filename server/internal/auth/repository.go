package auth

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		collection: db.Collection("user_credentials"),
	}
}

func (r *Repository) Create(ctx context.Context, credential *Credential) error {
	result, err := r.collection.InsertOne(ctx, credential)
	if err != nil {
		return fmt.Errorf("insert credential: %w", err)
	}

	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("inserted credential id is not an ObjectID")
	}

	credential.ID = objectID
	return nil
}

func (r *Repository) GetByEmailNormalized(ctx context.Context, emailNormalized string) (*Credential, error) {
	var credential Credential

	err := r.collection.FindOne(ctx, bson.M{"email_normalized": emailNormalized}).Decode(&credential)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("find credential by email: %w", err)
	}

	return &credential, nil
}
