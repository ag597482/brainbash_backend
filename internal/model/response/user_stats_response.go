package response

// GameTypeStats is per-game-type stats in GET /api/user/stats.
type GameTypeStats struct {
	AvgScore float64 `json:"avg_score"`
	MaxScore float64 `json:"max_score"`
}
