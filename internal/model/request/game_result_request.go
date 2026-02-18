package request

// GameResultRequest is the request body for POST /api/game/result.
type GameResultRequest struct {
	GameType          string             `json:"gametype" binding:"required"`
	QuestionResponses []QuestionResponse `json:"question_responses" binding:"required"`
}
