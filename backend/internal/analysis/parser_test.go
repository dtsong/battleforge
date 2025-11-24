package analysis

import (
	"strings"
	"testing"
)

func TestParseShowdownLogBasicValid(t *testing.T) {
	log := sampleBattleLog()
	summary, err := ParseShowdownLog(log)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if summary == nil {
		t.Fatal("expected summary, got nil")
	}

	if summary.ID == "" {
		t.Error("expected battle ID to be set")
	}

	if summary.Format == "" {
		t.Error("expected format to be set")
	}

	if summary.Player1.Name == "" || summary.Player2.Name == "" {
		t.Error("expected player names to be set")
	}

	if len(summary.Turns) == 0 {
		t.Error("expected turns to be parsed")
	}
}

func TestParseShowdownLogPlayerNames(t *testing.T) {
	log := sampleBattleLog()
	summary, _ := ParseShowdownLog(log)

	if summary.Player1.Name != "Player1" {
		t.Errorf("expected player1 name 'Player1', got %q", summary.Player1.Name)
	}

	if summary.Player2.Name != "Player2" {
		t.Errorf("expected player2 name 'Player2', got %q", summary.Player2.Name)
	}
}

func TestParseShowdownLogFormat(t *testing.T) {
	log := sampleBattleLog()
	summary, _ := ParseShowdownLog(log)

	if !strings.Contains(summary.Format, "VGC 2025") {
		t.Errorf("expected format to contain 'VGC 2025', got %q", summary.Format)
	}
}

func TestParseShowdownLogTurns(t *testing.T) {
	log := sampleBattleLog()
	summary, _ := ParseShowdownLog(log)

	if len(summary.Turns) == 0 {
		t.Fatal("expected at least one turn")
	}

	// Check turn numbers are sequential
	for i, turn := range summary.Turns {
		expectedTurn := i + 1
		if turn.TurnNumber != expectedTurn {
			t.Errorf("turn %d: expected turn number %d, got %d", i, expectedTurn, turn.TurnNumber)
		}
	}
}

func TestParseShowdownLogActions(t *testing.T) {
	log := sampleBattleLog()
	summary, _ := ParseShowdownLog(log)

	if len(summary.Turns) == 0 {
		t.Fatal("expected turns")
	}

	firstTurn := summary.Turns[0]
	if len(firstTurn.Actions) == 0 {
		t.Fatal("expected actions in first turn")
	}

	for _, action := range firstTurn.Actions {
		if action.Player != "player1" && action.Player != "player2" {
			t.Errorf("expected player to be player1 or player2, got %q", action.Player)
		}

		if action.ActionType == "" {
			t.Error("expected action type to be set")
		}
	}
}

func TestParseShowdownLogWinner(t *testing.T) {
	log := sampleBattleLog()
	summary, _ := ParseShowdownLog(log)

	if summary.Winner == "" {
		t.Error("expected winner to be set")
	}

	if summary.Winner != "player1" && summary.Winner != "player2" {
		t.Errorf("expected winner to be player1 or player2, got %q", summary.Winner)
	}
}

func TestParseShowdownLogStats(t *testing.T) {
	log := sampleBattleLog()
	summary, _ := ParseShowdownLog(log)

	if summary.Stats.TotalTurns != len(summary.Turns) {
		t.Errorf("expected total turns %d, got %d", len(summary.Turns), summary.Stats.TotalTurns)
	}

	if summary.Stats.MoveFrequency == nil {
		t.Error("expected move frequency map")
	}

	// Should have counted at least one move
	if len(summary.Stats.MoveFrequency) == 0 {
		t.Error("expected at least one move in frequency")
	}
}

func TestParseShowdownLogMinimalLog(t *testing.T) {
	log := minimalBattleLog()
	summary, err := ParseShowdownLog(log)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if summary == nil {
		t.Fatal("expected summary")
	}

	if summary.Player1.Name == "" || summary.Player2.Name == "" {
		t.Error("expected player names")
	}
}

func TestParseShowdownLogEmptyLog(t *testing.T) {
	log := ""
	summary, _ := ParseShowdownLog(log)

	// Empty log should not error, but return minimal summary
	if summary == nil {
		t.Fatal("expected summary even for empty log")
	}

	// Should have no turns
	if len(summary.Turns) > 0 {
		t.Error("expected no turns for empty log")
	}
}

func TestParseShowdownLogMalformedLog(t *testing.T) {
	log := malformedBattleLog()
	summary, _ := ParseShowdownLog(log)

	// Parser should be resilient
	if summary == nil {
		t.Fatal("expected summary even for malformed log")
	}

	// Should handle gracefully with minimal data
	if len(summary.Turns) > 0 && summary.Player1.Name == "" {
		t.Error("should still parse some structure")
	}
}

func TestParseShowdownLogMoveParsing(t *testing.T) {
	log := sampleBattleLog()
	summary, _ := ParseShowdownLog(log)

	// Find a move action
	foundMove := false
	for _, turn := range summary.Turns {
		for _, action := range turn.Actions {
			if action.ActionType == "move" && action.Move != nil {
				foundMove = true
				if action.Move.ID == "" {
					t.Error("expected move ID to be set")
				}
				if action.Move.Name == "" {
					t.Error("expected move name to be set")
				}
			}
		}
	}

	if !foundMove {
		t.Error("expected at least one move in the battle")
	}
}

func TestParseShowdownLogSwitchParsing(t *testing.T) {
	logWithSwitch := `|j|☆Player1
|j|☆Player2
|player|p1|Player1|test|1500
|player|p2|Player2|test|1500
|tier|[Gen 9] VGC 2025 Reg H (Bo3)
|start
|switch|p1a: Pikachu|Pikachu, L50, M|100/100
|turn|1
|move|p1a: Pikachu|Tackle|p2a: Charizard
|upkeep
|switch|p1a: Charizard|Charizard, L50, M|100/100
|turn|2
|move|p1a: Charizard|Flamethrower|p2a: Charizard
|upkeep
|win|Player1`

	summary, _ := ParseShowdownLog(logWithSwitch)

	// Find a switch action
	foundSwitch := false
	for _, turn := range summary.Turns {
		for _, action := range turn.Actions {
			if action.ActionType == "switch" {
				foundSwitch = true
				if action.SwitchTo == "" {
					t.Error("expected switch target to be set")
				}
			}
		}
	}

	if !foundSwitch {
		t.Error("expected at least one switch in the battle")
	}
}

func TestParseShowdownLogKeyMoments(t *testing.T) {
	logWithFaint := `|j|☆Player1
|j|☆Player2
|player|p1|Player1|test|1500
|player|p2|Player2|test|1500
|tier|[Gen 9] VGC 2025 Reg H (Bo3)
|start
|turn|1
|move|p1a: Pikachu|Thunderbolt|p2a: Blastoise
|-damage|p2a: Blastoise|0 fnt
|faint|p2a: Blastoise
|upkeep
|win|Player1`

	summary, _ := ParseShowdownLog(logWithFaint)

	if len(summary.KeyMoments) == 0 {
		t.Error("expected key moments to be recorded")
	}

	// Check for KO moment
	hasKO := false
	for _, moment := range summary.KeyMoments {
		if moment.Type == "KO" {
			hasKO = true
		}
	}

	if !hasKO {
		t.Error("expected KO to be recorded as key moment")
	}
}

func TestParseShowdownLogPlayerLosses(t *testing.T) {
	logWithFaint := `|j|☆Player1
|j|☆Player2
|player|p1|Player1|test|1500
|player|p2|Player2|test|1500
|tier|[Gen 9] VGC 2025 Reg H (Bo3)
|poke|p1|Pikachu, L50, M|
|poke|p1|Charizard, L50, M|
|poke|p2|Blastoise, L50, M|
|poke|p2|Dragonite, L50, M|
|teamsize|p1|2
|teamsize|p2|2
|start
|turn|1
|move|p1a: Pikachu|Thunderbolt|p2a: Blastoise
|-damage|p2a: Blastoise|0 fnt
|faint|p2a: Blastoise
|upkeep
|turn|2
|switch|p2a: Dragonite|Dragonite, L50, M|100/100
|move|p1a: Pikachu|Thunder Wave|p2a: Dragonite
|upkeep
|win|Player1`

	summary, _ := ParseShowdownLog(logWithFaint)

	if summary.Player2.Losses != 1 {
		t.Errorf("expected player2 to have 1 loss, got %d", summary.Player2.Losses)
	}
}

func TestParseShowdownLogUUIDUniqueness(t *testing.T) {
	log1 := sampleBattleLog()
	log2 := sampleBattleLog()

	summary1, _ := ParseShowdownLog(log1)
	summary2, _ := ParseShowdownLog(log2)

	if summary1.ID == summary2.ID {
		t.Error("expected different UUIDs for different parses")
	}
}

// Test fixtures

func sampleBattleLog() string {
	return `|j|☆Player1
|j|☆Player2
|html|<table width="100%"><tr><td align="left">Player1</td><td align="right">Player2</td></tr></table>
|t:|1763188046
|gametype|doubles
|player|p1|Player1|giovanni|1487
|player|p2|Player2|steven|1398
|gen|9
|tier|[Gen 9] VGC 2025 Reg H (Bo3)
|rated|
|rule|Species Clause: Limit one of each Pokémon
|rule|Item Clause: Limit 1 of each item
|clearpoke
|poke|p1|Pikachu, L50, M|
|poke|p1|Charizard, L50, M|
|poke|p2|Blastoise, L50, M|
|poke|p2|Dragonite, L50, M|
|teampreview|2
|teamsize|p1|2
|teamsize|p2|2
|start
|switch|p1a: Pikachu|Pikachu, L50, M|100/100
|switch|p2a: Blastoise|Blastoise, L50, M|100/100
|turn|1
|move|p1a: Pikachu|Thunderbolt|p2a: Blastoise
|-supereffective|p2a: Blastoise
|-damage|p2a: Blastoise|65/100
|move|p2a: Blastoise|Hydro Pump|p1a: Pikachu
|-supereffective|p1a: Pikachu
|-damage|p1a: Pikachu|30/100
|upkeep
|turn|2
|move|p1a: Pikachu|Thunder Wave|p2a: Blastoise
|-damage|p2a: Blastoise|60/100
|move|p2a: Blastoise|Protect|p2a: Blastoise
|-singleturn|p2a: Blastoise|Protect
|upkeep
|turn|3
|switch|p1a: Charizard|Charizard, L50, M|100/100
|move|p2a: Blastoise|Ice Beam|p1a: Charizard
|-supereffective|p1a: Charizard
|-damage|p1a: Charizard|40/100
|upkeep
|turn|4
|move|p1a: Charizard|Flamethrower|p2a: Blastoise
|-resisted|p2a: Blastoise
|-damage|p2a: Blastoise|30/100
|move|p2a: Blastoise|Waterfall|p1a: Charizard
|-supereffective|p1a: Charizard
|-damage|p1a: Charizard|0 fnt
|faint|p1a: Charizard
|upkeep
|
|switch|p1a: Pikachu|Pikachu, L50, M|30/100
|turn|5
|move|p1a: Pikachu|Quick Attack|p2a: Blastoise
|-damage|p2a: Blastoise|20/100
|move|p2a: Blastoise|Waterfall|p1a: Pikachu
|-supereffective|p1a: Pikachu
|-damage|p1a: Pikachu|0 fnt
|faint|p1a: Pikachu
|upkeep
|
|win|Player2`
}

func minimalBattleLog() string {
	return `|j|☆Player1
|j|☆Player2
|player|p1|Player1|test|1500
|player|p2|Player2|test|1500
|tier|[Gen 9] VGC 2025 Reg H (Bo3)
|start
|turn|1
|move|p1a: Test|Tackle|p2a: Test
|upkeep
|win|Player1`
}

func malformedBattleLog() string {
	return `this is not a valid showdown log
it has no pipe delimiters
and no proper structure`
}

// Edge case tests for comprehensive coverage

func TestParseShowdownLogDamageTracking(t *testing.T) {
	log := sampleBattleLog()
	summary, _ := ParseShowdownLog(log)

	// Verify damage is tracked
	p1Damage := summary.Stats.Player1Stats.DamageDealt
	p2Damage := summary.Stats.Player2Stats.DamageDealt
	if p1Damage <= 0 && p2Damage <= 0 {
		t.Error("expected damage to be tracked for at least one player")
	}
}

func TestParseShowdownLogEffectiveness(t *testing.T) {
	log := sampleBattleLog()
	summary, _ := ParseShowdownLog(log)

	if summary.Stats.SuperEffective < 0 {
		t.Error("expected super effective moves count to be non-negative")
	}

	if summary.Stats.NotVeryEffective < 0 {
		t.Error("expected not very effective moves count to be non-negative")
	}
}

func TestParseShowdownLogTurningPoints(t *testing.T) {
	log := sampleBattleLog()
	summary, _ := ParseShowdownLog(log)

	// Turning points may or may not exist depending on battle dynamics
	// Just verify the structure is valid
	for _, tp := range summary.Stats.TurningPoints {
		if tp.TurnNumber <= 0 {
			t.Error("expected turn number to be positive")
		}

		if tp.Score1Before < 0 || tp.Score1Before > 100 {
			t.Errorf("expected score1Before to be 0-100, got %v", tp.Score1Before)
		}
	}
}

func TestParseShowdownLogPositionScores(t *testing.T) {
	log := sampleBattleLog()
	summary, _ := ParseShowdownLog(log)

	// Verify position scores for each turn
	for i, turn := range summary.Turns {
		if turn.PositionScore == nil {
			t.Errorf("turn %d: expected position score", i+1)
			continue
		}

		// Scores should be between 0 and 100
		if turn.PositionScore.Player1Score < 0 || turn.PositionScore.Player1Score > 100 {
			t.Errorf("turn %d: player1 score out of range: %v", i+1, turn.PositionScore.Player1Score)
		}

		if turn.PositionScore.Player2Score < 0 || turn.PositionScore.Player2Score > 100 {
			t.Errorf("turn %d: player2 score out of range: %v", i+1, turn.PositionScore.Player2Score)
		}
	}
}

func TestParseShowdownLogConsecutiveSwitches(t *testing.T) {
	logWithSwitches := `|j|☆Player1
|j|☆Player2
|player|p1|Player1|test|1500
|player|p2|Player2|test|1500
|tier|[Gen 9] VGC 2025 Reg H (Bo3)
|poke|p1|Poke1, L50|
|poke|p1|Poke2, L50|
|poke|p1|Poke3, L50|
|poke|p2|Poke1, L50|
|poke|p2|Poke2, L50|
|teamsize|p1|3
|teamsize|p2|2
|start
|turn|1
|switch|p1a: Poke1|Poke1, L50|100/100
|switch|p2a: Poke1|Poke1, L50|100/100
|move|p1a: Poke1|Tackle|p2a: Poke1
|move|p2a: Poke1|Tackle|p1a: Poke1
|upkeep
|turn|2
|switch|p1a: Poke2|Poke2, L50|100/100
|switch|p2a: Poke2|Poke2, L50|100/100
|move|p1a: Poke2|Tackle|p2a: Poke2
|move|p2a: Poke2|Tackle|p1a: Poke2
|upkeep
|win|Player1`

	summary, _ := ParseShowdownLog(logWithSwitches)

	totalSwitches := summary.Stats.Player1Stats.SwitchCount + summary.Stats.Player2Stats.SwitchCount
	if totalSwitches != 4 {
		t.Errorf("expected 4 switches total, got %d", totalSwitches)
	}
}

func TestParseShowdownLogPartialDamage(t *testing.T) {
	logPartialDamage := `|j|☆Player1
|j|☆Player2
|player|p1|Player1|test|1500
|player|p2|Player2|test|1500
|tier|[Gen 9] VGC 2025 Reg H (Bo3)
|poke|p1|Pikachu, L50|
|poke|p2|Charizard, L50|
|teamsize|p1|1
|teamsize|p2|1
|start
|turn|1
|move|p1a: Pikachu|Thunder|p2a: Charizard
|-damage|p2a: Charizard|50/100
|move|p2a: Charizard|Flare Blitz|p1a: Pikachu
|-damage|p1a: Pikachu|25/100
|upkeep
|win|Player1`

	summary, _ := ParseShowdownLog(logPartialDamage)

	// Both players should have taken damage
	p1Taken := summary.Stats.Player1Stats.DamageTaken
	p2Taken := summary.Stats.Player2Stats.DamageTaken
	if p1Taken == 0 || p2Taken == 0 {
		t.Error("expected both players to have taken damage")
	}
}

func TestParseShowdownLogNoMoves(t *testing.T) {
	logNoMoves := `|j|☆Player1
|j|☆Player2
|player|p1|Player1|test|1500
|player|p2|Player2|test|1500
|tier|[Gen 9] VGC 2025 Reg H (Bo3)
|start
|turn|1
|switch|p1a: Pikachu|Pikachu, L50|100/100
|switch|p2a: Charizard|Charizard, L50|100/100
|upkeep
|win|Player1`

	summary, _ := ParseShowdownLog(logNoMoves)

	if summary == nil {
		t.Fatal("expected summary")
	}

	// Should still have turns even without moves
	if len(summary.Turns) == 0 {
		t.Error("expected at least one turn")
	}
}

func TestParseShowdownLogCriticalHits(t *testing.T) {
	logCritical := `|j|☆Player1
|j|☆Player2
|player|p1|Player1|test|1500
|player|p2|Player2|test|1500
|tier|[Gen 9] VGC 2025 Reg H (Bo3)
|poke|p1|Pikachu, L50|
|poke|p2|Charizard, L50|
|teamsize|p1|1
|teamsize|p2|1
|start
|turn|1
|move|p1a: Pikachu|Thunder|p2a: Charizard
|-crit|p2a: Charizard
|-damage|p2a: Charizard|40/100
|move|p2a: Charizard|Flare Blitz|p1a: Pikachu
|upkeep
|win|Player1`

	summary, _ := ParseShowdownLog(logCritical)

	if summary.Stats.CriticalHits == 0 {
		t.Error("expected critical hit to be counted")
	}
}

func TestParseShowdownLogMoveFrequency(t *testing.T) {
	log := sampleBattleLog()
	summary, _ := ParseShowdownLog(log)

	if len(summary.Stats.MoveFrequency) == 0 {
		t.Fatal("expected move frequency map to be populated")
	}

	// All move counts should be positive
	for move, count := range summary.Stats.MoveFrequency {
		if count <= 0 {
			t.Errorf("move %q has non-positive count: %d", move, count)
		}
	}
}

func TestParseShowdownLogPlayerTeamTracking(t *testing.T) {
	log := sampleBattleLog()
	summary, _ := ParseShowdownLog(log)

	if len(summary.Player1.Team) == 0 {
		t.Error("expected player1 team")
	}

	if len(summary.Player2.Team) == 0 {
		t.Error("expected player2 team")
	}

	// Verify team sizes match
	if summary.Player1.TotalLeft != len(summary.Player1.Team) {
		t.Errorf("player1 total left %d doesn't match team size %d", summary.Player1.TotalLeft, len(summary.Player1.Team))
	}
}

func TestParseShowdownLogMultipleFaints(t *testing.T) {
	logMultipleFaints := `|j|☆Player1
|j|☆Player2
|player|p1|Player1|test|1500
|player|p2|Player2|test|1500
|tier|[Gen 9] VGC 2025 Reg H (Bo3)
|poke|p1|Poke1, L50|
|poke|p1|Poke2, L50|
|poke|p2|Poke1, L50|
|poke|p2|Poke2, L50|
|poke|p2|Poke3, L50|
|teamsize|p1|2
|teamsize|p2|3
|start
|turn|1
|move|p1a: Poke1|Tackle|p2a: Poke1
|-damage|p2a: Poke1|0 fnt
|faint|p2a: Poke1
|move|p2a: Poke2|Tackle|p1a: Poke1
|upkeep
|turn|2
|switch|p2a: Poke3|Poke3, L50|100/100
|move|p1a: Poke1|Tackle|p2a: Poke3
|-damage|p2a: Poke3|0 fnt
|faint|p2a: Poke3
|move|p2a: Poke2|Tackle|p1a: Poke1
|-damage|p1a: Poke1|0 fnt
|faint|p1a: Poke1
|upkeep
|turn|3
|switch|p1a: Poke2|Poke2, L50|100/100
|move|p1a: Poke2|Tackle|p2a: Poke2
|move|p2a: Poke2|Tackle|p1a: Poke2
|upkeep
|win|Player1`

	summary, _ := ParseShowdownLog(logMultipleFaints)

	if summary.Player2.Losses != 2 {
		t.Errorf("expected player2 to have 2 losses, got %d", summary.Player2.Losses)
	}

	if summary.Player1.Losses != 1 {
		t.Errorf("expected player1 to have 1 loss, got %d", summary.Player1.Losses)
	}
}
