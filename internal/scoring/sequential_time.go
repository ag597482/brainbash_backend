package scoring

import "brainbash_backend/internal/model/request"

// SequentialTimeStrategy scores when the next question comes only after the previous is solved.
// Each item in question_responses is one solved question; score is out of 100 based on count.
// All responses count as correct; avgTime is average time per question.
type SequentialTimeStrategy struct{}

func NewSequentialTimeStrategy() *SequentialTimeStrategy {
	return &SequentialTimeStrategy{}
}

func (s *SequentialTimeStrategy) Calculate(responses []request.QuestionResponse) *ScoreResult {
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

	var totalTime float64
	for _, r := range responses {
		totalTime += r.TimeTaken
	}

	avgTime := totalTime / float64(n)
	// All answered = 100% accuracy; score out of 100 based on number of questions solved.
	// Here we treat "score" as 100 when at least one question is solved (scale as you need).
	score := 100.0
	if n == 0 {
		score = 0
	}

	return &ScoreResult{
		Score:     score,
		Questions: n,
		Correct:   n,
		Accuracy:  1.0,
		AvgTime:   avgTime,
	}
}
