package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"brainbash_backend/internal/game"
	"brainbash_backend/internal/model/entity"
	"brainbash_backend/internal/model/request"
	"brainbash_backend/internal/model/response"
	"brainbash_backend/internal/scoring"
	"brainbash_backend/internal/service"
	"brainbash_backend/internal/utils"
)

// ScoreController handles score calculation, game result submission, and user stats.
type ScoreController struct {
	scorer       *scoring.Scorer
	scoreService *service.ScoreService
}

// NewScoreController creates a new ScoreController.
func NewScoreController(scorer *scoring.Scorer, scoreService *service.ScoreService) *ScoreController {
	return &ScoreController{
		scorer:       scorer,
		scoreService: scoreService,
	}
}

// Calculate handles POST /score. Expects strategy and question_responses; returns score, questions, correct, accuracy, avgTime.
func (sc *ScoreController) Calculate(c *gin.Context) {
	var req request.ScoringRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: strategy and question_responses required"})
		return
	}

	result, err := sc.scorer.Calculate(req.Strategy, req.QuestionResponses)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.ScoringResponse{
		Score:     result.Score,
		Questions: result.Questions,
		Correct:   result.Correct,
		Accuracy:  result.Accuracy,
		AvgTime:   result.AvgTime,
	})
}

// GameCalculate handles POST /api/game/guest/result. Same request as /api/game/result (gametype, question_responses),
// but no auth and no DB: only computes score and returns the result.
func (sc *ScoreController) GameCalculate(c *gin.Context) {
	var req request.GameResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: gametype and question_responses are required"})
		return
	}

	gt := game.GameType(req.GameType)
	if err := gt.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := sc.scorer.Calculate(gt.StrategyFor(), req.QuestionResponses)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.ScoringResponse{
		Score:     result.Score,
		Questions: result.Questions,
		Correct:   result.Correct,
		Accuracy:  result.Accuracy,
		AvgTime:   result.AvgTime,
	})
}

// GameResult handles POST /api/game/result. Calculates score, stores session in score collection, returns result.
func (sc *ScoreController) GameResult(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user context"})
		return
	}

	var req request.GameResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gametype and question_responses are required"})
		return
	}

	result, err := sc.scoreService.SubmitGameResult(c.Request.Context(), userID, req)
	if err != nil {
		log.Printf("GameResult SubmitGameResult: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.ScoringResponse{
		Score:     result.Score,
		Questions: result.Questions,
		Correct:   result.Correct,
		Accuracy:  result.Accuracy,
		AvgTime:   result.AvgTime,
	})
}

// UserStats handles GET /api/user/stats. Returns the authenticated user's scores per game type.
func (sc *ScoreController) UserStats(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user context"})
		return
	}

	score, err := sc.scoreService.GetUserStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load stats"})
		return
	}

	// Response: { overall_score, <gametype>: { avg_score, max_score }, ... } — all game types included, 0 when no data
	out := make(map[string]interface{})
	out["overall_score"] = 0.0
	if score != nil {
		out["overall_score"] = score.OverallScore
	}
	for _, kv := range gameTypeKeysAndValues(score) {
		if kv.value != nil {
			out[kv.key] = response.GameTypeStats{
				AvgScore: kv.value.AvgScore,
				MaxScore: kv.value.HighScore,
			}
		} else {
			out[kv.key] = response.GameTypeStats{AvgScore: 0, MaxScore: 0}
		}
	}

	c.JSON(http.StatusOK, out)
}

type gameTypeKeyVal struct {
	key   string
	value *entity.GameTypeScore
}

func gameTypeKeysAndValues(score *entity.Score) []gameTypeKeyVal {
	if score == nil {
		return []gameTypeKeyVal{
			{"processing_speed", nil}, {"working_memory", nil},
			{"logical_reasoning", nil}, {"math_reasoning", nil},
			{"reflex_time", nil}, {"attention_control", nil},
		}
	}
	return []gameTypeKeyVal{
		{"processing_speed", score.ProcessingSpeed},
		{"working_memory", score.WorkingMemory},
		{"logical_reasoning", score.LogicalReasoning},
		{"math_reasoning", score.MathReasoning},
		{"reflex_time", score.ReflexTime},
		{"attention_control", score.AttentionControl},
	}
}
