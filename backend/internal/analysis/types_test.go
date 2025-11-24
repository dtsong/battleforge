package analysis

import (
	"testing"
	"time"
)

func TestBattleSummaryInitialization(t *testing.T) {
	summary := &BattleSummary{
		ID:        "test-123",
		Format:    "gen9vgc2025reghbo3",
		Timestamp: time.Now(),
		Duration:  300,
		Player1: Player{
			Name:      "Player1",
			TotalLeft: 6,
		},
		Player2: Player{
			Name:      "Player2",
			TotalLeft: 6,
		},
		Winner: "player1",
	}

	if summary.ID != "test-123" {
		t.Errorf("expected ID 'test-123', got %q", summary.ID)
	}

	if summary.Format == "" {
		t.Error("expected format to be set")
	}

	if summary.Player1.Name == "" || summary.Player2.Name == "" {
		t.Error("expected both players to be named")
	}

	if summary.Winner != "player1" && summary.Winner != "player2" && summary.Winner != "draw" {
		t.Errorf("expected winner to be player1, player2, or draw, got %q", summary.Winner)
	}
}

func TestPlayerTeamManagement(t *testing.T) {
	player := &Player{
		Name:      "TestPlayer",
		TotalLeft: 6,
		Team: []Pokémon{
			{ID: "pikachu", Name: "Pikachu", Level: 50},
			{ID: "charizard", Name: "Charizard", Level: 50},
		},
	}

	if len(player.Team) != 2 {
		t.Errorf("expected 2 pokémon in team, got %d", len(player.Team))
	}

	if player.TotalLeft != 6 {
		t.Errorf("expected 6 pokémon total, got %d", player.TotalLeft)
	}

	// Simulate losing one pokémon
	player.Losses++
	expectedRemaining := player.TotalLeft - player.Losses
	if expectedRemaining != 5 {
		t.Errorf("expected 5 remaining after 1 loss, got %d", expectedRemaining)
	}
}

func TestPokémonStats(t *testing.T) {
	poke := &Pokémon{
		ID:        "pikachu",
		Name:      "Pikachu",
		Level:     50,
		CurrentHP: 90,
		MaxHP:     100,
		Status:    "paralysis",
		TeraType:  "Electric",
		Stats: Stats{
			HP:      35,
			Attack:  55,
			Defense: 40,
			SpAtk:   50,
			SpDef:   50,
			Speed:   90,
		},
	}

	if poke.CurrentHP > poke.MaxHP {
		t.Error("current HP should not exceed max HP")
	}

	if poke.Stats.Speed < 0 {
		t.Error("stat values should be non-negative")
	}

	// Calculate HP percentage
	hpPercent := float64(poke.CurrentHP) / float64(poke.MaxHP) * 100
	if hpPercent > 100 {
		t.Error("HP percentage should not exceed 100")
	}
}

func TestMoveProperties(t *testing.T) {
	move := &Move{
		ID:       "thunderbolt",
		Name:     "Thunderbolt",
		Type:     "Electric",
		Power:    90,
		Accuracy: 100,
		PP:       15,
	}

	if move.Power < 0 {
		t.Error("power should be non-negative")
	}

	if move.Accuracy < 0 || move.Accuracy > 100 {
		t.Errorf("accuracy should be 0-100, got %d", move.Accuracy)
	}

	if move.PP < 0 {
		t.Error("PP should be non-negative")
	}

	if move.Type == "" {
		t.Error("move type should be set")
	}
}

func TestTurnStructure(t *testing.T) {
	turn := &Turn{
		TurnNumber: 1,
		Actions: []Action{
			{
				Player:     "player1",
				ActionType: "move",
				Move: &Move{
					ID:   "thunderbolt",
					Name: "Thunderbolt",
					Type: "Electric",
				},
			},
		},
		DamageDealt: map[string]int{
			"player1": 50,
			"player2": 30,
		},
		PositionScore: &PositionScore{
			Player1Score:   75.0,
			Player2Score:   60.0,
			MomentumPlayer: "player1",
		},
	}

	if turn.TurnNumber <= 0 {
		t.Error("turn number should be positive")
	}

	if len(turn.Actions) == 0 {
		t.Error("turn should have actions")
	}

	for player, damage := range turn.DamageDealt {
		if damage < 0 {
			t.Errorf("damage for %q should be non-negative, got %d", player, damage)
		}
	}

	if turn.PositionScore == nil {
		t.Error("turn should have position score")
	}

	if turn.PositionScore.Player1Score < 0 || turn.PositionScore.Player1Score > 100 {
		t.Errorf("position score should be 0-100, got %v", turn.PositionScore.Player1Score)
	}
}

func TestPositionScore(t *testing.T) {
	tests := []struct {
		name         string
		p1Score      float64
		p2Score      float64
		expectValid  bool
	}{
		{
			name:        "valid scores",
			p1Score:     75.0,
			p2Score:     60.0,
			expectValid: true,
		},
		{
			name:        "equal scores",
			p1Score:     50.0,
			p2Score:     50.0,
			expectValid: true,
		},
		{
			name:        "extreme p1 advantage",
			p1Score:     100.0,
			p2Score:     0.0,
			expectValid: true,
		},
		{
			name:        "extreme p2 advantage",
			p1Score:     0.0,
			p2Score:     100.0,
			expectValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := &PositionScore{
				Player1Score:   tt.p1Score,
				Player2Score:   tt.p2Score,
				MomentumPlayer: "player1",
			}

			if score.Player1Score < 0 || score.Player1Score > 100 {
				t.Errorf("player1 score out of range: %v", score.Player1Score)
			}

			if score.Player2Score < 0 || score.Player2Score > 100 {
				t.Errorf("player2 score out of range: %v", score.Player2Score)
			}
		})
	}
}

func TestTurningPoint(t *testing.T) {
	tp := &TurningPoint{
		TurnNumber:   5,
		Score1Before: 45.0,
		Score1After:  65.0,
		Score2Before: 55.0,
		Score2After:  35.0,
		MomentumShift: 20.0,
		Significance: 8,
		Description:  "Player 1 ko'd player 2's lead",
	}

	if tp.TurnNumber <= 0 {
		t.Error("turn number should be positive")
	}

	if tp.Significance < 1 || tp.Significance > 10 {
		t.Errorf("significance should be 1-10, got %d", tp.Significance)
	}

	if tp.Description == "" {
		t.Error("description should be set")
	}

	// Verify score changes make sense
	p1Change := tp.Score1After - tp.Score1Before
	if p1Change != 20.0 {
		t.Errorf("expected p1 change of 20, got %v", p1Change)
	}
}

func TestBattleStats(t *testing.T) {
	stats := &BattleStats{
		TotalTurns:       10,
		Switch:           4,
		CriticalHits:     2,
		SuperEffective:   5,
		NotVeryEffective: 3,
		AvgDamagePerTurn: 25.5,
		AvgHealPerTurn:   5.0,
		MoveFrequency: map[string]int{
			"thunderbolt": 3,
			"flamethrower": 2,
			"protect":     2,
		},
		TypeCoverage: map[string]int{
			"electric": 4,
			"fire":     2,
			"water":    3,
		},
	}

	if stats.TotalTurns <= 0 {
		t.Error("total turns should be positive")
	}

	if stats.Switch < 0 {
		t.Error("switches should be non-negative")
	}

	if stats.CriticalHits < 0 {
		t.Error("critical hits should be non-negative")
	}

	if stats.AvgDamagePerTurn < 0 {
		t.Error("average damage should be non-negative")
	}

	// Verify move frequency
	totalMoves := 0
	for move, count := range stats.MoveFrequency {
		if count <= 0 {
			t.Errorf("move %q should have positive count, got %d", move, count)
		}
		totalMoves += count
	}

	if totalMoves <= 0 {
		t.Error("should have counted moves")
	}
}

func TestPlayerStats(t *testing.T) {
	pStats := &PlayerStats{
		MoveCount:       10,
		SwitchCount:     3,
		DamageDealt:     250,
		DamageTaken:     150,
		HealingDone:     20,
		HealingReceived: 10,
		MovesByType: map[string]int{
			"electric": 5,
			"water":    3,
			"psychic":  2,
		},
		Effectiveness: EffectivenessStats{
			SuperEffective:   4,
			NotVeryEffective: 2,
			Neutral:          4,
		},
	}

	if pStats.MoveCount < 0 {
		t.Error("move count should be non-negative")
	}

	if pStats.DamageDealt < 0 {
		t.Error("damage dealt should be non-negative")
	}

	if pStats.SwitchCount < 0 {
		t.Error("switch count should be non-negative")
	}

	totalEffectiveness := pStats.Effectiveness.SuperEffective + pStats.Effectiveness.NotVeryEffective + pStats.Effectiveness.Neutral
	if totalEffectiveness <= 0 {
		t.Error("should have some effectiveness data")
	}
}

func TestKeyMoment(t *testing.T) {
	moment := &KeyMoment{
		TurnNumber:  3,
		Type:        "KO",
		Description: "Player 2's Pikachu was knocked out",
		Significance: 9,
	}

	if moment.TurnNumber <= 0 {
		t.Error("turn number should be positive")
	}

	if moment.Type == "" {
		t.Error("moment type should be set")
	}

	if moment.Description == "" {
		t.Error("description should be set")
	}

	if moment.Significance < 1 || moment.Significance > 10 {
		t.Errorf("significance should be 1-10, got %d", moment.Significance)
	}

	validTypes := map[string]bool{
		"KO":       true,
		"switch":   true,
		"status":   true,
		"weather":  true,
		"crit":     true,
		"other":    true,
	}

	if !validTypes[moment.Type] {
		t.Errorf("moment type %q not recognized", moment.Type)
	}
}

func TestActionTypes(t *testing.T) {
	tests := []struct {
		name       string
		action     Action
		shouldHave string
	}{
		{
			name: "move action",
			action: Action{
				Player:     "player1",
				ActionType: "move",
				Move: &Move{
					ID:   "tackle",
					Name: "Tackle",
					Type: "Normal",
				},
			},
			shouldHave: "move",
		},
		{
			name: "switch action",
			action: Action{
				Player:     "player2",
				ActionType: "switch",
				SwitchTo:   "Pikachu",
			},
			shouldHave: "switch",
		},
		{
			name: "item action",
			action: Action{
				Player:     "player1",
				ActionType: "item",
				Item:       "Full Restore",
			},
			shouldHave: "item",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.action.ActionType != tt.shouldHave {
				t.Errorf("expected action type %q, got %q", tt.shouldHave, tt.action.ActionType)
			}

			if tt.action.Player == "" {
				t.Error("player should be set")
			}
		})
	}
}
