package spaces

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
		collection: db.Collection("spaces"),
	}
}

func (r *Repository) Create(ctx context.Context, space *Space) error {
	result, err := r.collection.InsertOne(ctx, space)
	if err != nil {
		return fmt.Errorf("insert space: %w", err)
	}

	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("inserted space id is not an ObjectID")
	}

	space.ID = objectID
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id primitive.ObjectID) (*Space, error) {
	var space Space

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&space)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("find space by id: %w", err)
	}

	return &space, nil
}
