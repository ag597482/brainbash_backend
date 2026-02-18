package scoring

import "brainbash_backend/internal/model/request"

const outcomeCorrect = "correct"

// TimedOutcomeStrategy scores based on time_taken and outcome (correct/incorrect/unsolved).
// Score out of 100 is driven by accuracy; avgTime is average time per question.
type TimedOutcomeStrategy struct{}

func NewTimedOutcomeStrategy() *TimedOutcomeStrategy {
	return &TimedOutcomeStrategy{}
}

func (s *TimedOutcomeStrategy) Calculate(responses []request.QuestionResponse) *ScoreResult {
	n := len(responses)
	if n == 0 {
		return &ScoreResult{
			Score:     0,
			Questions: 0,
			Correct:   0,
			Accuracy:  0,
			AvgTime:   0,
		}
	}

	var correct int
	var totalTime float64
	for _, r := range responses {
		totalTime += r.TimeTaken
		if r.Outcome == outcomeCorrect {
			correct++
		}
	}

	accuracy := 0.0
	if n > 0 {
		accuracy = float64(correct) / float64(n)
	}
	score := accuracy * 100
	avgTime := totalTime / float64(n)

	return &ScoreResult{
		Score:     score,
		Questions: n,
		Correct:   correct,
		Accuracy:  accuracy,
		AvgTime:   avgTime,
	}
}
