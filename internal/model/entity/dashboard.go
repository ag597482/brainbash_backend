package entity

import "time"

// Dashboard is the document stored in the "dashboard" collection.
// A single document (e.g. _id: "leaderboard") holds top 10 entries per game type.
type Dashboard struct {
	ID                string                     `bson:"_id,omitempty"`
	ProcessingSpeed   []DashboardEntry           `bson:"processing_speed,omitempty"`
	WorkingMemory     []DashboardEntry           `bson:"working_memory,omitempty"`
	LogicalReasoning  []DashboardEntry           `bson:"logical_reasoning,omitempty"`
	MathReasoning     []DashboardEntry           `bson:"math_reasoning,omitempty"`
	ReflexTime        []DashboardEntry           `bson:"reflex_time,omitempty"`
	AttentionControl  []DashboardEntry           `bson:"attention_control,omitempty"`
}

// DashboardEntry is one top-score entry for a game type (session + user summary + score).
type DashboardEntry struct {
	SessionID    string              `bson:"session_id"    json:"session_id"`
	User         DashboardUserSummary `bson:"user"          json:"user"`
	SessionScore SessionScoreDetail  `bson:"session_score" json:"session_score"`
	Timestamp    time.Time           `bson:"timestamp"     json:"timestamp"`
}

// DashboardUserSummary is the user info embedded in a dashboard entry.
type DashboardUserSummary struct {
	ID     string `bson:"_id"     json:"_id"`
	GaID   string `bson:"gaid"    json:"gaid"`
	Name   string `bson:"name"    json:"name"`
	Email  string `bson:"email"   json:"email"`
	Photo  string `bson:"photo"   json:"photo"`
}

const DashboardDocID = "leaderboard"
const DashboardTopN = 10
