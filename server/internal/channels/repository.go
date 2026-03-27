package channels

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
		collection: db.Collection("channels"),
	}
}

func (r *Repository) Create(ctx context.Context, channel *Channel) error {
	result, err := r.collection.InsertOne(ctx, channel)
	if err != nil {
		return fmt.Errorf("insert channel: %w", err)
	}

	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("inserted channel id is not an ObjectID")
	}

	channel.ID = objectID
	return nil
}

func (r *Repository) ListBySpaceID(ctx context.Context, spaceID primitive.ObjectID) ([]Channel, error) {
	filter := bson.M{"space_id": spaceID}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("find channels by space id: %w", err)
	}
	defer cursor.Close(ctx)

	var channels []Channel
	if err := cursor.All(ctx, &channels); err != nil {
		return nil, fmt.Errorf("decode channels: %w", err)
	}

	return channels, nil
}
