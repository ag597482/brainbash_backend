package repository

import (
	"context"
	"fmt"
	"time"

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

// DeleteEntriesInDateRange removes dashboard entries whose timestamp falls within [start, end].
// Finds the dashboard doc, filters each game-type array in memory, then upserts.
func (r *DashboardRepository) DeleteEntriesInDateRange(ctx context.Context, start, end time.Time) error {
	d, err := r.FindByID(ctx, entity.DashboardDocID)
	if err != nil {
		return err
	}
	if d == nil {
		return nil
	}
	d.ProcessingSpeed = filterEntriesByDateRange(d.ProcessingSpeed, start, end)
	d.WorkingMemory = filterEntriesByDateRange(d.WorkingMemory, start, end)
	d.LogicalReasoning = filterEntriesByDateRange(d.LogicalReasoning, start, end)
	d.MathReasoning = filterEntriesByDateRange(d.MathReasoning, start, end)
	d.ReflexTime = filterEntriesByDateRange(d.ReflexTime, start, end)
	d.AttentionControl = filterEntriesByDateRange(d.AttentionControl, start, end)
	return r.Upsert(ctx, d)
}

func filterEntriesByDateRange(entries []entity.DashboardEntry, start, end time.Time) []entity.DashboardEntry {
	out := make([]entity.DashboardEntry, 0, len(entries))
	for _, e := range entries {
		t := e.Timestamp
		if t.Before(start) || t.After(end) {
			out = append(out, e)
		}
	}
	return out
}
