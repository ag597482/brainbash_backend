package response

// ScoringResponse is the response body for the scoring API.
type ScoringResponse struct {
	Score     float64 `json:"score"`     // 0–100
	Questions int     `json:"questions"` // total number of questions
	Correct   int     `json:"correct"`   // number correct (or solved in sequential)
	Accuracy  float64 `json:"accuracy"`  // correct / questions (0–1)
	AvgTime   float64 `json:"avgTime"`   // average time per question in seconds
}
