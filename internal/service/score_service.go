package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"brainbash_backend/internal/game"
	"brainbash_backend/internal/model/entity"
	"brainbash_backend/internal/model/request"
	"brainbash_backend/internal/repository"
	"brainbash_backend/internal/scoring"
)

// ScoreService appends sessions and maintains per-game-type and overall scores.
type ScoreService struct {
	scoreRepo *repository.ScoreRepository
	scorer    *scoring.Scorer
}

// NewScoreService creates a new ScoreService.
func NewScoreService(scoreRepo *repository.ScoreRepository, scorer *scoring.Scorer) *ScoreService {
	return &ScoreService{scoreRepo: scoreRepo, scorer: scorer}
}

// SubmitGameResult validates gametype, calculates score, persists the session, and returns the score result.
func (s *ScoreService) SubmitGameResult(ctx context.Context, userID string, req request.GameResultRequest) (*scoring.ScoreResult, error) {
	gt := game.GameType(req.GameType)
	if err := gt.Validate(); err != nil {
		return nil, err
	}

	strategy := gt.StrategyFor()
	result, err := s.scorer.Calculate(strategy, req.QuestionResponses)
	if err != nil {
		return nil, err
	}

	if err := s.AppendSession(ctx, userID, req.GameType, req.QuestionResponses, result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetUserStats returns the user's score document from the scores collection, or nil if not found.
func (s *ScoreService) GetUserStats(ctx context.Context, userID string) (*entity.Score, error) {
	return s.scoreRepo.FindByUserID(ctx, userID)
}

// AppendSession adds a session for the user and game type, then recomputes avg_score, high_score, and overall_score.
func (s *ScoreService) AppendSession(ctx context.Context, userID, gameType string, questionResponses interface{}, result *scoring.ScoreResult) error {
	score, err := s.scoreRepo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if score == nil {
		score = &entity.Score{UserID: userID, OverallScore: result.Score}
	}

	session := entity.Session{
		SessionID:        bson.NewObjectID().Hex(),
		QuestionResponses: questionResponses,
		SessionScore: entity.SessionScoreDetail{
			Score:     result.Score,
			Questions: result.Questions,
			Correct:   result.Correct,
			Accuracy:  result.Accuracy,
			AvgTime:   result.AvgTime,
		},
		Timestamp: time.Now().UTC(),
	}

	gt := getGameTypeScore(score, gameType)
	if gt == nil {
		gt = &entity.GameTypeScore{Sessions: []entity.Session{}}
		setGameTypeScore(score, gameType, gt)
	}
	gt.Sessions = append(gt.Sessions, session)

	// Recompute avg and high for this game type (using session_score.score)
	var sum float64
	high := 0.0
	for _, se := range gt.Sessions {
		sum += se.SessionScore.Score
		if se.SessionScore.Score > high {
			high = se.SessionScore.Score
		}
	}
	gt.AvgScore = sum / float64(len(gt.Sessions))
	gt.HighScore = high

	// Recompute overall_score (average of all session scores across all game types)
	var totalSum float64
	var totalCount int
	for _, g := range allGameTypeScores(score) {
		if g == nil {
			continue
		}
		for _, se := range g.Sessions {
			totalSum += se.SessionScore.Score
			totalCount++
		}
	}
	if totalCount > 0 {
		score.OverallScore = totalSum / float64(totalCount)
	} else {
		score.OverallScore = result.Score
	}

	return s.scoreRepo.Upsert(ctx, score)
}

func getGameTypeScore(score *entity.Score, gameType string) *entity.GameTypeScore {
	switch gameType {
	case "processing_speed":
		return score.ProcessingSpeed
	case "working_memory":
		return score.WorkingMemory
	case "logical_reasoning":
		return score.LogicalReasoning
	case "math_reasoning":
		return score.MathReasoning
	case "reflex_time":
		return score.ReflexTime
	case "attention_control":
		return score.AttentionControl
	default:
		return nil
	}
}

func setGameTypeScore(score *entity.Score, gameType string, gt *entity.GameTypeScore) {
	switch gameType {
	case "processing_speed":
		score.ProcessingSpeed = gt
	case "working_memory":
		score.WorkingMemory = gt
	case "logical_reasoning":
		score.LogicalReasoning = gt
	case "math_reasoning":
		score.MathReasoning = gt
	case "reflex_time":
		score.ReflexTime = gt
	case "attention_control":
		score.AttentionControl = gt
	}
}

func allGameTypeScores(score *entity.Score) []*entity.GameTypeScore {
	return []*entity.GameTypeScore{
		score.ProcessingSpeed, score.WorkingMemory, score.LogicalReasoning,
		score.MathReasoning, score.ReflexTime, score.AttentionControl,
	}
}
