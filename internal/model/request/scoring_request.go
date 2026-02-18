package request

// ScoringRequest is the request body for the scoring API.
type ScoringRequest struct {
	Strategy          string              `json:"strategy" binding:"required"` // "timed_outcome" or "sequential_time"
	QuestionResponses []QuestionResponse  `json:"question_responses" binding:"required"`
}

// QuestionResponse represents one question's response.
// For timed_outcome: TimeTaken and Outcome (correct/incorrect/unsolved) are used.
// For sequential_time: only TimeTaken is used (each item = one solved question).
type QuestionResponse struct {
	TimeTaken float64 `json:"time_taken"`           // seconds
	Outcome   string  `json:"outcome,omitempty"`   // "correct" | "incorrect" | "unsolved" (strategy 1 only)
}
