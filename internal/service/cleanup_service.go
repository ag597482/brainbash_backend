package service

import (
	"context"
	"time"

	"brainbash_backend/internal/model/entity"
	"brainbash_backend/internal/repository"
)

// CleanupService removes sessions and dashboard entries within a date range.
type CleanupService struct {
	scoreRepo     *repository.ScoreRepository
	dashboardRepo *repository.DashboardRepository
}

// NewCleanupService creates a new CleanupService.
func NewCleanupService(scoreRepo *repository.ScoreRepository, dashboardRepo *repository.DashboardRepository) *CleanupService {
	return &CleanupService{
		scoreRepo:     scoreRepo,
		dashboardRepo: dashboardRepo,
	}
}

// CleanupByDateRange removes from scores (sessions) and dashboard (entries) all data
// whose timestamp falls within [start, end] (inclusive). For each affected user score,
// per-game avg_score and high_score and overall_score are recomputed from remaining
// sessions and the updated document is persisted.
func (s *CleanupService) CleanupByDateRange(ctx context.Context, start, end time.Time) (scoresUpdated int, err error) {
	scores, err := s.scoreRepo.FindAll(ctx)
	if err != nil {
		return 0, err
	}

	for _, score := range scores {
		changed := s.removeSessionsInDateRange(score, start, end)
		if changed {
			// Persist score with recomputed avg_score, high_score (per game) and overall_score
			if err := s.scoreRepo.Upsert(ctx, score); err != nil {
				return scoresUpdated, err
			}
			scoresUpdated++
		}
	}

	if err := s.dashboardRepo.DeleteEntriesInDateRange(ctx, start, end); err != nil {
		return scoresUpdated, err
	}

	return scoresUpdated, nil
}

// removeSessionsInDateRange filters out sessions in [start, end] from each game type,
// then recomputes avg_score and high_score for each game type and overall_score from
// the remaining sessions. The score struct is updated in place; caller must Upsert to persist.
// Returns true if any session was removed.
func (s *CleanupService) removeSessionsInDateRange(score *entity.Score, start, end time.Time) bool {
	anyChanged := false
	for _, gt := range []*entity.GameTypeScore{
		score.ProcessingSpeed, score.WorkingMemory, score.LogicalReasoning,
		score.MathReasoning, score.ReflexTime, score.AttentionControl,
	} {
		if gt == nil {
			continue
		}
		kept := filterSessionsByDateRange(gt.Sessions, start, end)
		if len(kept) != len(gt.Sessions) {
			anyChanged = true
		}
		gt.Sessions = kept
		// Recompute avg_score and high_score from remaining sessions
		var sum float64
		high := 0.0
		for _, se := range gt.Sessions {
			sum += se.SessionScore.Score
			if se.SessionScore.Score > high {
				high = se.SessionScore.Score
			}
		}
		if len(gt.Sessions) > 0 {
			gt.AvgScore = sum / float64(len(gt.Sessions))
			gt.HighScore = high
		} else {
			gt.AvgScore = 0
			gt.HighScore = 0
		}
	}

	// Recompute overall_score (average of all remaining session scores across all game types)
	var totalSum float64
	var totalCount int
	for _, gt := range []*entity.GameTypeScore{
		score.ProcessingSpeed, score.WorkingMemory, score.LogicalReasoning,
		score.MathReasoning, score.ReflexTime, score.AttentionControl,
	} {
		if gt == nil {
			continue
		}
		for _, se := range gt.Sessions {
			totalSum += se.SessionScore.Score
			totalCount++
		}
	}
	if totalCount > 0 {
		score.OverallScore = totalSum / float64(totalCount)
	} else {
		score.OverallScore = 0
	}

	return anyChanged
}

func filterSessionsByDateRange(sessions []entity.Session, start, end time.Time) []entity.Session {
	out := make([]entity.Session, 0, len(sessions))
	for _, se := range sessions {
		t := se.Timestamp
		if t.Before(start) || t.After(end) {
			out = append(out, se)
		}
	}
	return out
}
