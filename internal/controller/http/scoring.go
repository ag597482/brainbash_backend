package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"brainbash_backend/internal/model/request"
	"brainbash_backend/internal/model/response"
	"brainbash_backend/internal/scoring"
	"brainbash_backend/internal/service"
	"brainbash_backend/internal/utils"
)

// ScoringController handles score calculation requests.
type ScoringController struct {
	scorer       *scoring.Scorer
	scoreService *service.ScoreService
}

// NewScoringController creates a new ScoringController.
func NewScoringController(scorer *scoring.Scorer, scoreService *service.ScoreService) *ScoringController {
	return &ScoringController{
		scorer:       scorer,
		scoreService: scoreService,
	}
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

// GameResult handles POST /api/game/result. Calculates score, stores session in score collection, returns result.
func (sc *ScoringController) GameResult(c *gin.Context) {
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
