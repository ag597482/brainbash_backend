package scoring

import "brainbash_backend/internal/model/request"

// ScoreResult holds the result of a scoring calculation.
type ScoreResult struct {
	Score     float64
	Questions int
	Correct   int
	Accuracy  float64
	AvgTime   float64
}

// Strategy defines how to compute a score from question responses.
type Strategy interface {
	Calculate(responses []request.QuestionResponse) *ScoreResult
}

const (
	StrategyTimedOutcome   = "timed_outcome"   // questions have outcome: correct/incorrect/unsolved
	StrategySequentialTime = "sequential_time" // next question only after previous solved; only time_taken
)
