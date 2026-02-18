package game

import (
	"fmt"

	"brainbash_backend/internal/scoring"
)

// GameType represents the allowed game types for /api/game/result.
type GameType string

const (
	ProcessingSpeed   GameType = "processing_speed"   // Processing Speed
	WorkingMemory     GameType = "working_memory"     // Working Memory
	LogicalReasoning  GameType = "logical_reasoning"  // Logical Reasoning
	MathReasoning     GameType = "math_reasoning"     // Math Reasoning
	ReflexTime        GameType = "reflex_time"       // Reflex Time
	AttentionControl  GameType = "attention_control"  // Attention Control
)

// AllGameTypes is the list of valid game types (for validation).
var AllGameTypes = []GameType{
	ProcessingSpeed, WorkingMemory, LogicalReasoning,
	MathReasoning, ReflexTime, AttentionControl,
}

var validGameTypes = map[GameType]struct{}{
	ProcessingSpeed:  {},
	WorkingMemory:    {},
	LogicalReasoning: {},
	MathReasoning:    {},
	ReflexTime:       {},
	AttentionControl: {},
}

// IsValid returns true if g is an allowed game type.
func (g GameType) IsValid() bool {
	_, ok := validGameTypes[g]
	return ok
}

// StrategyFor returns the scoring strategy to use for this game type.
// reflex_time -> sequential_time; all others -> timed_outcome.
func (g GameType) StrategyFor() string {
	if g == ReflexTime {
		return scoring.StrategySequentialTime
	}
	return scoring.StrategyTimedOutcome
}

// Validate returns an error if the game type is not allowed.
func (g GameType) Validate() error {
	if g.IsValid() {
		return nil
	}
	return fmt.Errorf("invalid gametype: %q (allowed: processing_speed, working_memory, logical_reasoning, math_reasoning, reflex_time, attention_control)", g)
}
