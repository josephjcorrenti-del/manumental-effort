package messages

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		collection: db.Collection("messages"),
	}
}

func (r *Repository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "channel_id", Value: 1},
				{Key: "_id", Value: -1},
			},
			Options: options.Index().SetName("channel_id_message_id_desc"),
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("create message indexes: %w", err)
	}

	return nil
}

func (r *Repository) Create(ctx context.Context, message *Message) error {
	result, err := r.collection.InsertOne(ctx, message)
	if err != nil {
		return fmt.Errorf("insert message: %w", err)
	}

	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("inserted message id is not an ObjectID")
	}

	message.ID = objectID
	return nil
}

func (r *Repository) ListByChannelBefore(
	ctx context.Context,
	channelID primitive.ObjectID,
	beforeID *primitive.ObjectID,
	limit int,
) ([]Message, error) {
	filter := bson.M{
		"channel_id": channelID,
	}

	if beforeID != nil {
		filter["_id"] = bson.M{"$lt": *beforeID}
	}

	findOptions := options.Find().
		SetSort(bson.D{{Key: "_id", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("find messages by channel: %w", err)
	}
	defer cursor.Close(ctx)

	var messages []Message
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, fmt.Errorf("decode messages: %w", err)
	}

	return messages, nil
}
