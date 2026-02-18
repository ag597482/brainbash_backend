package scoring

import (
	"fmt"

	"brainbash_backend/internal/model/request"
)

// Scorer selects a strategy by name and computes the score.
type Scorer struct {
	strategies map[string]Strategy
}

// NewScorer builds a Scorer with all strategies registered.
func NewScorer() *Scorer {
	return &Scorer{
		strategies: map[string]Strategy{
			StrategyTimedOutcome:   NewTimedOutcomeStrategy(),
			StrategySequentialTime: NewSequentialTimeStrategy(),
		},
	}
}

// Calculate returns the score result for the given strategy and responses.
// Returns error if strategy is unknown.
func (sc *Scorer) Calculate(strategyName string, responses []request.QuestionResponse) (*ScoreResult, error) {
	s, ok := sc.strategies[strategyName]
	if !ok {
		return nil, fmt.Errorf("unknown strategy: %s", strategyName)
	}
	return s.Calculate(responses), nil
}
