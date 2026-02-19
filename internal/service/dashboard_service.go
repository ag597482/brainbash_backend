package service

import (
	"context"
	"sort"
	"time"

	"brainbash_backend/internal/model/entity"
	"brainbash_backend/internal/repository"
)

// DashboardService provides dashboard (top-10 leaderboard) read and update.
type DashboardService struct {
	dashboardRepo *repository.DashboardRepository
	userService   *UserService
}

// NewDashboardService creates a new DashboardService.
func NewDashboardService(dashboardRepo *repository.DashboardRepository, userService *UserService) *DashboardService {
	return &DashboardService{
		dashboardRepo: dashboardRepo,
		userService:   userService,
	}
}

// GetDashboard returns the full dashboard (top 10 per game type). Returns empty dashboard if not found.
func (s *DashboardService) GetDashboard(ctx context.Context) (*entity.Dashboard, error) {
	d, err := s.dashboardRepo.FindByID(ctx, entity.DashboardDocID)
	if err != nil {
		return nil, err
	}
	if d == nil {
		return &entity.Dashboard{}, nil
	}
	return d, nil
}

// MaybeUpdateTop10 adds the given session to the dashboard for the game type if it qualifies for top 10.
// Called after each game result. userID is the authenticated user's ID; sessionScore and timestamp describe the session.
func (s *DashboardService) MaybeUpdateTop10(ctx context.Context, gameType, userID, sessionID string, sessionScore entity.SessionScoreDetail, timestamp time.Time) error {
	user, err := s.userService.FindByUserID(ctx, userID)
	if err != nil || user == nil {
		return nil // skip update if user not found
	}

	entry := entity.DashboardEntry{
		SessionID: sessionID,
		User: entity.DashboardUserSummary{
			ID:    user.UserID.Hex(),
			GaID:  user.GaID,
			Name:  user.Name,
			Email: user.Email,
			Photo: user.Picture,
		},
		SessionScore: sessionScore,
		Timestamp:    timestamp,
	}

	d, err := s.dashboardRepo.FindByID(ctx, entity.DashboardDocID)
	if err != nil {
		return err
	}
	if d == nil {
		d = &entity.Dashboard{}
	}

	list := getDashboardEntriesForGameType(d, gameType)
	list = append(list, entry)
	sort.Slice(list, func(i, j int) bool {
		return list[i].SessionScore.Score > list[j].SessionScore.Score
	})
	if len(list) > entity.DashboardTopN {
		list = list[:entity.DashboardTopN]
	}
	setDashboardEntriesForGameType(d, gameType, list)

	return s.dashboardRepo.Upsert(ctx, d)
}

func getDashboardEntriesForGameType(d *entity.Dashboard, gameType string) []entity.DashboardEntry {
	switch gameType {
	case "processing_speed":
		return d.ProcessingSpeed
	case "working_memory":
		return d.WorkingMemory
	case "logical_reasoning":
		return d.LogicalReasoning
	case "math_reasoning":
		return d.MathReasoning
	case "reflex_time":
		return d.ReflexTime
	case "attention_control":
		return d.AttentionControl
	default:
		return nil
	}
}

func setDashboardEntriesForGameType(d *entity.Dashboard, gameType string, entries []entity.DashboardEntry) {
	switch gameType {
	case "processing_speed":
		d.ProcessingSpeed = entries
	case "working_memory":
		d.WorkingMemory = entries
	case "logical_reasoning":
		d.LogicalReasoning = entries
	case "math_reasoning":
		d.MathReasoning = entries
	case "reflex_time":
		d.ReflexTime = entries
	case "attention_control":
		d.AttentionControl = entries
	}
}
