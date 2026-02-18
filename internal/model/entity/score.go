package entity

import "time"

// Score is the document stored in the "scores" collection (one per user).
// _id is the user_id. Each game type has avg_score, high_score, and sessions list.
type Score struct {
	ID               string          `bson:"_id,omitempty"` // user_id
	UserID           string          `bson:"user_id"`
	OverallScore     float64         `bson:"overall_score"`
	ProcessingSpeed  *GameTypeScore  `bson:"processing_speed,omitempty"`
	WorkingMemory    *GameTypeScore  `bson:"working_memory,omitempty"`
	LogicalReasoning *GameTypeScore  `bson:"logical_reasoning,omitempty"`
	MathReasoning    *GameTypeScore  `bson:"math_reasoning,omitempty"`
	ReflexTime       *GameTypeScore  `bson:"reflex_time,omitempty"`
	AttentionControl *GameTypeScore  `bson:"attention_control,omitempty"`
}

// GameTypeScore holds per-game-type aggregates and sessions.
type GameTypeScore struct {
	AvgScore  float64   `bson:"avg_score"`
	HighScore float64   `bson:"high_score"`
	Sessions  []Session `bson:"sessions"`
}

// Session is one game session for a game type.
type Session struct {
	SessionID        string            `bson:"session_id"`
	QuestionResponses interface{}      `bson:"question_responses"`
	SessionScore     SessionScoreDetail `bson:"session_score"`
	Timestamp        time.Time         `bson:"timestamp"`
}

// SessionScoreDetail is the score breakdown stored per session.
type SessionScoreDetail struct {
	Score     float64 `bson:"score"`
	Questions int     `bson:"questions"`
	Correct   int     `bson:"correct"`
	Accuracy  float64 `bson:"accuracy"`
	AvgTime   float64 `bson:"avgTime"`
}
