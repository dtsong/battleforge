package db

import "time"

// Battle represents a stored battle record.
type Battle struct {
	ID         string
	Format     string
	Timestamp  time.Time
	DurationSec int
	Winner     string // "player1", "player2", or "draw"
	Player1ID  string
	Player2ID  string
	BattleLog  string
	IsPrivate  bool
	Analysis   *BattleAnalysis
	KeyMoments []*KeyMoment
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// BattleAnalysis stores computed statistics for a battle.
type BattleAnalysis struct {
	BattleID                string
	TotalTurns              int
	AvgDamagePerTurn        float64
	AvgHealPerTurn          float64
	MovesUsedCount          int
	SwitchesCount           int
	SuperEffectiveMoves     int
	NotVeryEffectiveMoves   int
	CriticalHits            int
	Player1DamageDealt      int
	Player1DamageTaken      int
	Player1HealingDone      int
	Player2DamageDealt      int
	Player2DamageTaken      int
	Player2HealingDone      int
	CreatedAt               time.Time
}

// KeyMoment represents a significant moment in a battle.
type KeyMoment struct {
	BattleID    string
	TurnNumber  int
	MomentType  string // "switch", "ko", "status", "weather", "critical", "other"
	Description string
	Significance int
	CreatedAt   time.Time
}

// BattleFilter is used for filtering battles in queries.
type BattleFilter struct {
	Format    string
	IsPrivate *bool
}
