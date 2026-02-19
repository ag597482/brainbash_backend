package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"brainbash_backend/internal/model/entity"
)

const dashboardCollection = "dashboard"

// DashboardRepository handles MongoDB operations for the dashboard (leaderboard) collection.
type DashboardRepository struct {
	collection *mongo.Collection
}

// NewDashboardRepository creates a new DashboardRepository.
func NewDashboardRepository(db *mongo.Database) *DashboardRepository {
	return &DashboardRepository{
		collection: db.Collection(dashboardCollection),
	}
}

// FindByID returns the dashboard document (e.g. leaderboard), or nil if not found.
func (r *DashboardRepository) FindByID(ctx context.Context, id string) (*entity.Dashboard, error) {
	var doc entity.Dashboard
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("find dashboard: %w", err)
	}
	return &doc, nil
}

// Upsert replaces the dashboard document (full document replace).
func (r *DashboardRepository) Upsert(ctx context.Context, d *entity.Dashboard) error {
	d.ID = entity.DashboardDocID
	filter := bson.M{"_id": d.ID}
	opts := options.Replace().SetUpsert(true)
	_, err := r.collection.ReplaceOne(ctx, filter, d, opts)
	if err != nil {
		return fmt.Errorf("upsert dashboard: %w", err)
	}
	return nil
}
