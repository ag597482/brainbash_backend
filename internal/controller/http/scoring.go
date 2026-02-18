package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"brainbash_backend/internal/game"
	"brainbash_backend/internal/model/request"
	"brainbash_backend/internal/model/response"
	"brainbash_backend/internal/scoring"
)

// ScoringController handles score calculation requests.
type ScoringController struct {
	scorer *scoring.Scorer
}

// NewScoringController creates a new ScoringController.
func NewScoringController(scorer *scoring.Scorer) *ScoringController {
	return &ScoringController{scorer: scorer}
}

// Calculate handles POST /score. Expects strategy and question_responses; returns score, questions, correct, accuracy, avgTime.
func (sc *ScoringController) Calculate(c *gin.Context) {
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

// GameResult handles POST /api/game/result. Accepts only allowed gametype and question_responses; maps gametype to scoring strategy internally.
func (sc *ScoringController) GameResult(c *gin.Context) {
	var req request.GameResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gametype and question_responses are required"})
		return
	}

	gt := game.GameType(req.GameType)
	if err := gt.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	strategy := gt.StrategyFor()
	result, err := sc.scorer.Calculate(strategy, req.QuestionResponses)
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
