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

func (s *Service) ListSpaces(ctx context.Context, userID string) ([]Space, error) {
	objectUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id")
	}

	memberships, err := s.membershipRepository.ListByUserID(ctx, objectUserID)
	if err != nil {
		return nil, err
	}

	spaceIDs := make([]primitive.ObjectID, 0, len(memberships))
	for _, m := range memberships {
		spaceIDs = append(spaceIDs, m.SpaceID)
	}

	spaces, err := s.repository.ListByIDs(ctx, spaceIDs)
	if err != nil {
		return nil, err
	}

	return spaces, nil
}

func (r *Repository) ListByIDs(ctx context.Context, ids []primitive.ObjectID) ([]Space, error) {
	if len(ids) == 0 {
		return []Space{}, nil
	}

	filter := bson.M{
		"_id": bson.M{"$in": ids},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("find spaces by ids: %w", err)
	}
	defer cursor.Close(ctx)

	var spaces []Space
	if err := cursor.All(ctx, &spaces); err != nil {
		return nil, fmt.Errorf("decode spaces: %w", err)
	}

	return spaces, nil
}
