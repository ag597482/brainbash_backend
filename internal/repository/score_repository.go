package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"brainbash_backend/internal/model/entity"
)

const scoreCollection = "scores"

// ScoreRepository handles MongoDB operations for the score collection.
type ScoreRepository struct {
	collection *mongo.Collection
}

// NewScoreRepository creates a new ScoreRepository.
func NewScoreRepository(db *mongo.Database) *ScoreRepository {
	return &ScoreRepository{
		collection: db.Collection(scoreCollection),
	}
}

// FindByUserID returns the score document for the user, or nil if not found.
func (r *ScoreRepository) FindByUserID(ctx context.Context, userID string) (*entity.Score, error) {
	var doc entity.Score
	err := r.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("find score by user_id: %w", err)
	}
	return &doc, nil
}

// Upsert replaces the score document for the user (one doc per user, keyed by user_id).
func (r *ScoreRepository) Upsert(ctx context.Context, score *entity.Score) error {
	score.ID = score.UserID
	filter := bson.M{"user_id": score.UserID}
	opts := options.Replace().SetUpsert(true)
	_, err := r.collection.ReplaceOne(ctx, filter, score, opts)
	if err != nil {
		return fmt.Errorf("upsert score: %w", err)
	}
	return nil
}

// FindAll returns all score documents (for cleanup by date range).
func (r *ScoreRepository) FindAll(ctx context.Context) ([]*entity.Score, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("find all scores: %w", err)
	}
	defer cursor.Close(ctx)

	var out []*entity.Score
	if err := cursor.All(ctx, &out); err != nil {
		return nil, fmt.Errorf("decode scores: %w", err)
	}
	return out, nil
}
