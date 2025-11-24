package analysis

import (
	"math"
	"strings"
	"testing"
)

// ===== Edge Case Tests =====

func TestParseShowdownLogNoPlayer(t *testing.T) {
	log := `|tier|[Gen 9] VGC 2025 Reg H (Bo3)
|start
|turn|1
|move|p1a: Pikachu|Tackle|p2a: Charizard
|upkeep
|win|Player1`

	summary, err := ParseShowdownLog(log)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Should handle missing player declarations gracefully
	if summary.Player1.Name == "" {
		t.Error("expected player1 to have default name")
	}
}

func TestParseShowdownLogNoWinner(t *testing.T) {
	log := `|j|☆Player1
|j|☆Player2
|player|p1|Player1|test|1500
|player|p2|Player2|test|1500
|tier|[Gen 9] VGC 2025 Reg H (Bo3)
|start
|turn|1
|move|p1a: Pikachu|Tackle|p2a: Charizard
|upkeep`

	summary, _ := ParseShowdownLog(log)

	// Should not crash if no winner found
	if summary == nil {
		t.Fatal("expected summary")
	}
}

func TestParseShowdownLogMalformedLines(t *testing.T) {
	log := `|j|☆Player1
|j|☆Player2
|player|p1|Player1|test|1500
|player|p2|Player2|test|1500
|tier|[Gen 9] VGC 2025 Reg H (Bo3)
|this line is missing pipes
|start
|turn|1
|move|p1a: Pikachu|Tackle|p2a: Charizard
|malformed|without|enough|pipes|||
|upkeep
|win|Player1`

	summary, err := ParseShowdownLog(log)

	// Should be resilient and not crash
	if err != nil && !strings.Contains(err.Error(), "parse") {
		t.Fatalf("expected nil error or parse error, got %v", err)
	}

	if summary == nil {
		t.Fatal("expected summary even with malformed lines")
	}
}

func TestParseShowdownLogInvalidHP(t *testing.T) {
	log := `|j|☆Player1
|j|☆Player2
|player|p1|Player1|test|1500
|player|p2|Player2|test|1500
|tier|[Gen 9] VGC 2025 Reg H (Bo3)
|start
|turn|1
|move|p1a: Pikachu|Tackle|p2a: Charizard
|-damage|p2a: Charizard|invalid_hp
|upkeep
|win|Player1`

	summary, _ := ParseShowdownLog(log)

	// Should handle invalid HP gracefully
	if summary == nil {
		t.Fatal("expected summary")
	}
}

func TestParseShowdownLogNegativeHP(t *testing.T) {
	log := `|j|☆Player1
|j|☆Player2
|player|p1|Player1|test|1500
|player|p2|Player2|test|1500
|tier|[Gen 9] VGC 2025 Reg H (Bo3)
|start
|turn|1
|move|p1a: Pikachu|Tackle|p2a: Charizard
|-damage|p2a: Charizard|-50/100
|upkeep
|win|Player1`

	summary, _ := ParseShowdownLog(log)

	if summary == nil {
		t.Fatal("expected summary")
	}

	// HP should not be negative
	for _, turn := range summary.Turns {
		if turn.StateAfter.Player1Active != nil && turn.StateAfter.Player1Active.CurrentHP < 0 {
			t.Error("expected HP to not be negative")
		}
		if turn.StateAfter.Player2Active != nil && turn.StateAfter.Player2Active.CurrentHP < 0 {
			t.Error("expected HP to not be negative")
		}
	}
}

func TestParseShowdownLogHPOver100Percent(t *testing.T) {
	log := `|j|☆Player1
|j|☆Player2
|player|p1|Player1|test|1500
|player|p2|Player2|test|1500
|tier|[Gen 9] VGC 2025 Reg H (Bo3)
|start
|turn|1
|move|p1a: Pikachu|Recover|p1a: Pikachu
|-heal|p1a: Pikachu|200/100
|upkeep
|win|Player1`

	summary, _ := ParseShowdownLog(log)

	if summary == nil {
		t.Fatal("expected summary")
	}

	// HP should cap at max
	for _, turn := range summary.Turns {
		if turn.StateAfter.Player1Active != nil {
			if turn.StateAfter.Player1Active.CurrentHP > turn.StateAfter.Player1Active.MaxHP {
				t.Errorf("expected HP <= maxHP, got %d/%d", turn.StateAfter.Player1Active.CurrentHP, turn.StateAfter.Player1Active.MaxHP)
			}
		}
	}
}

func TestParseShowdownLogDuplicateTurns(t *testing.T) {
	log := `|j|☆Player1
|j|☆Player2
|player|p1|Player1|test|1500
|player|p2|Player2|test|1500
|tier|[Gen 9] VGC 2025 Reg H (Bo3)
|start
|turn|1
|move|p1a: Pikachu|Tackle|p2a: Charizard
|upkeep
|turn|1
|move|p1a: Pikachu|Tackle|p2a: Charizard
|upkeep
|turn|2
|move|p2a: Charizard|Flare Blitz|p1a: Pikachu
|upkeep
|win|Player1`

	summary, _ := ParseShowdownLog(log)

	if summary == nil {
		t.Fatal("expected summary")
	}

	// Should handle duplicate turn numbers gracefully
	if len(summary.Turns) == 0 {
		t.Error("expected at least one turn")
	}
}

func TestParseShowdownLogOutOfOrderTurns(t *testing.T) {
	log := `|j|☆Player1
|j|☆Player2
|player|p1|Player1|test|1500
|player|p2|Player2|test|1500
|tier|[Gen 9] VGC 2025 Reg H (Bo3)
|start
|turn|3
|move|p1a: Pikachu|Tackle|p2a: Charizard
|upkeep
|turn|1
|move|p1a: Pikachu|Tackle|p2a: Charizard
|upkeep
|turn|2
|move|p2a: Charizard|Flare Blitz|p1a: Pikachu
|upkeep
|win|Player1`

	summary, _ := ParseShowdownLog(log)

	if summary == nil {
		t.Fatal("expected summary")
	}

	// Should handle out-of-order turns (parser may reorder or keep as-is)
	if len(summary.Turns) == 0 {
		t.Error("expected at least one turn")
	}
}

func TestParseShowdownLogVeryLongBattle(t *testing.T) {
	// Generate a battle with 100 turns
	logLines := []string{
		"|j|☆Player1",
		"|j|☆Player2",
		"|player|p1|Player1|test|1500",
		"|player|p2|Player2|test|1500",
		"|tier|[Gen 9] VGC 2025 Reg H (Bo3)",
		"|start",
	}

	for i := 1; i <= 100; i++ {
		logLines = append(logLines, "|turn|"+string(rune(i)))
		logLines = append(logLines, "|move|p1a: Pikachu|Tackle|p2a: Charizard")
		logLines = append(logLines, "|upkeep")
	}
	logLines = append(logLines, "|win|Player1")

	log := strings.Join(logLines, "\n")
	summary, err := ParseShowdownLog(log)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if summary == nil {
		t.Fatal("expected summary")
	}

	// Should handle long battles efficiently
	if len(summary.Turns) < 50 {
		t.Errorf("expected many turns, got %d", len(summary.Turns))
	}
}

func TestParseShowdownLogSpecialCharactersInNames(t *testing.T) {
	log := `|j|☆Player1
|j|☆Player2
|player|p1|Player♔One|test|1500
|player|p2|Plâyer★Twö|test|1500
|tier|[Gen 9] VGC 2025 Reg H (Bo3)
|start
|turn|1
|move|p1a: Nidoking♔|Tackle|p2a: Charizard★
|upkeep
|win|Player♔One`

	summary, _ := ParseShowdownLog(log)

	if summary == nil {
		t.Fatal("expected summary")
	}

	// Should handle special characters
	if summary.Player1.Name == "" {
		t.Error("expected player1 name")
	}
}

func TestParseShowdownLogWhitespaceVariations(t *testing.T) {
	log := `|j|☆Player1
|j|☆Player2
|player|p1|  Player1  |test|1500
|player|p2|Player2|test|1500
|tier|[Gen 9] VGC 2025 Reg H (Bo3)
|start
|turn|1
|move|p1a:  Pikachu  |Tackle  |p2a: Charizard
|upkeep
|win|Player1`

	summary, _ := ParseShowdownLog(log)

	if summary == nil {
		t.Fatal("expected summary")
	}

	// Should handle whitespace variations
	if summary.Player1.Name == "" {
		t.Error("expected player1 name despite whitespace")
	}
}

// ===== Position Scoring Tests =====

func TestPositionScoringBasic(t *testing.T) {
	log := sampleBattleLog()
	summary, _ := ParseShowdownLog(log)

	// All turns should have position scores
	for i, turn := range summary.Turns {
		if turn.PositionScore == nil {
			t.Errorf("turn %d: expected position score", i+1)
			continue
		}

		// Scores should be between 0 and 100
		if turn.PositionScore.Player1Score < 0 || turn.PositionScore.Player1Score > 100 {
			t.Errorf("turn %d: player1 score out of range: %f", i+1, turn.PositionScore.Player1Score)
		}

		if turn.PositionScore.Player2Score < 0 || turn.PositionScore.Player2Score > 100 {
			t.Errorf("turn %d: player2 score out of range: %f", i+1, turn.PositionScore.Player2Score)
		}
	}
}

func TestPositionScoringMomentumPlayer(t *testing.T) {
	log := sampleBattleLog()
	summary, _ := ParseShowdownLog(log)

	for i, turn := range summary.Turns {
		if turn.PositionScore == nil {
			continue
		}

		momentum := turn.PositionScore.MomentumPlayer
		if momentum != "player1" && momentum != "player2" && momentum != "neutral" {
			t.Errorf("turn %d: invalid momentum player: %q", i+1, momentum)
		}
	}
}

func TestPositionScoringProgression(t *testing.T) {
	log := logWithPositionProgression()
	summary, _ := ParseShowdownLog(log)

	if len(summary.Turns) < 2 {
		t.Skip("test requires at least 2 turns")
	}

	// First turn should have some initial score
	if summary.Turns[0].PositionScore == nil {
		t.Fatal("expected position score on first turn")
	}

	if summary.Turns[0].PositionScore.Player1Score == 0 && summary.Turns[0].PositionScore.Player2Score == 0 {
		t.Error("expected non-zero position scores")
	}
}

// ===== Turning Point Tests =====

func TestTurningPointsBasic(t *testing.T) {
	log := sampleBattleLog()
	summary, _ := ParseShowdownLog(log)

	// Should have turning points
	if len(summary.Stats.TurningPoints) == 0 {
		t.Error("expected turning points in battle summary")
	}

	for _, tp := range summary.Stats.TurningPoints {
		if tp.TurnNumber <= 0 {
			t.Errorf("expected valid turn number, got %d", tp.TurnNumber)
		}

		if math.Abs(tp.MomentumShift) < 15 {
			t.Errorf("expected significant momentum shift (>15), got %f", tp.MomentumShift)
		}

		if tp.Significance < 1 || tp.Significance > 10 {
			t.Errorf("expected significance between 1-10, got %d", tp.Significance)
		}

		if tp.Description == "" {
			t.Error("expected turning point description")
		}
	}
}

func TestTurningPointsSignificance(t *testing.T) {
	log := logWithSignificantMomentum()
	summary, _ := ParseShowdownLog(log)

	if len(summary.Stats.TurningPoints) == 0 {
		t.Skip("no turning points in test log")
	}

	// Turning points should be sorted by turn number
	for i := 0; i < len(summary.Stats.TurningPoints)-1; i++ {
		if summary.Stats.TurningPoints[i].TurnNumber > summary.Stats.TurningPoints[i+1].TurnNumber {
			t.Error("turning points should be sorted by turn number")
		}
	}
}

// ===== Damage and Healing Tests =====

func TestDamageTracking(t *testing.T) {
	log := sampleBattleLog()
	summary, _ := ParseShowdownLog(log)

	totalP1Damage := 0
	totalP2Damage := 0

	for _, turn := range summary.Turns {
		if turn.DamageDealt["player1"] > 0 {
			totalP1Damage += turn.DamageDealt["player1"]
		}
		if turn.DamageDealt["player2"] > 0 {
			totalP2Damage += turn.DamageDealt["player2"]
		}
	}

	// Should have tracked damage
	if totalP1Damage == 0 && totalP2Damage == 0 {
		t.Error("expected damage to be tracked")
	}
}

func TestHealingTracking(t *testing.T) {
	logWithHealing := `|j|☆Player1
|j|☆Player2
|player|p1|Player1|test|1500
|player|p2|Player2|test|1500
|tier|[Gen 9] VGC 2025 Reg H (Bo3)
|poke|p1|Pikachu, L50, M|
|poke|p2|Charizard, L50, M|
|teamsize|p1|1
|teamsize|p2|1
|start
|switch|p1a: Pikachu|Pikachu, L50, M|100/100
|switch|p2a: Charizard|Charizard, L50, M|100/100
|turn|1
|move|p1a: Pikachu|Tackle|p2a: Charizard
|-damage|p2a: Charizard|50/100
|move|p2a: Charizard|Recover|p2a: Charizard
|-heal|p2a: Charizard|100/100
|upkeep
|turn|2
|move|p1a: Pikachu|Thunderbolt|p2a: Charizard
|-damage|p2a: Charizard|0 fnt
|faint|p2a: Charizard
|upkeep
|win|Player1`

	summary, _ := ParseShowdownLog(logWithHealing)

	if summary == nil {
		t.Fatal("expected summary")
	}

	totalHealing := 0
	for _, turn := range summary.Turns {
		for _, healed := range turn.HealingDone {
			totalHealing += healed
		}
	}

	if totalHealing == 0 {
		t.Error("expected healing to be tracked")
	}
}

// ===== Status Conditions Tests =====

func TestStatusConditionParsing(t *testing.T) {
	logWithStatus := `|j|☆Player1
|j|☆Player2
|player|p1|Player1|test|1500
|player|p2|Player2|test|1500
|tier|[Gen 9] VGC 2025 Reg H (Bo3)
|poke|p1|Pikachu, L50, M|
|poke|p2|Charizard, L50, M|
|teamsize|p1|1
|teamsize|p2|1
|start
|switch|p1a: Pikachu|Pikachu, L50, M|100/100
|switch|p2a: Charizard|Charizard, L50, M|100/100
|turn|1
|move|p1a: Pikachu|Thunder Wave|p2a: Charizard
|-status|p2a: Charizard|par
|upkeep
|turn|2
|move|p2a: Charizard|Flare Blitz|p1a: Pikachu
|-damage|p1a: Pikachu|50/100
|-status|p1a: Pikachu|brn
|upkeep
|win|Player1`

	summary, _ := ParseShowdownLog(logWithStatus)

	if summary == nil {
		t.Fatal("expected summary")
	}

	// Should track status conditions
	hasStatus := false
	for _, turn := range summary.Turns {
		if turn.StateAfter.Player1Active != nil && turn.StateAfter.Player1Active.Status != "" {
			hasStatus = true
		}
		if turn.StateAfter.Player2Active != nil && turn.StateAfter.Player2Active.Status != "" {
			hasStatus = true
		}
	}

	if !hasStatus {
		t.Error("expected status conditions to be tracked")
	}
}

// ===== Move Parsing Tests =====

func TestMoveWithPower(t *testing.T) {
	log := sampleBattleLog()
	summary, _ := ParseShowdownLog(log)

	foundMovesWithPower := false
	for _, turn := range summary.Turns {
		for _, action := range turn.Actions {
			if action.Move != nil && action.Move.Power > 0 {
				foundMovesWithPower = true
				break
			}
		}
	}

	if !foundMovesWithPower {
		t.Error("expected moves with power")
	}
}

func TestMoveFrequencyStats(t *testing.T) {
	log := sampleBattleLog()
	summary, _ := ParseShowdownLog(log)

	if len(summary.Stats.MoveFrequency) == 0 {
		t.Error("expected move frequency stats")
	}

	totalMoves := 0
	for _, count := range summary.Stats.MoveFrequency {
		if count < 0 {
			t.Error("expected non-negative move counts")
		}
		totalMoves += count
	}

	if totalMoves == 0 {
		t.Error("expected at least one move")
	}
}

// ===== Switch Parsing Tests =====

func TestSwitchTracking(t *testing.T) {
	logWithSwitch := `|j|☆Player1
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
|switch|p1a: Pikachu|Pikachu, L50, M|100/100
|switch|p2a: Blastoise|Blastoise, L50, M|100/100
|turn|1
|move|p1a: Pikachu|Tackle|p2a: Blastoise
|upkeep
|switch|p1a: Charizard|Charizard, L50, M|100/100
|turn|2
|move|p1a: Charizard|Flare Blitz|p2a: Blastoise
|upkeep
|switch|p2a: Dragonite|Dragonite, L50, M|100/100
|turn|3
|move|p1a: Charizard|Earthquake|p2a: Dragonite
|upkeep
|win|Player1`

	summary, _ := ParseShowdownLog(logWithSwitch)

	if summary == nil {
		t.Fatal("expected summary")
	}

	switchCount := 0
	for _, turn := range summary.Turns {
		for _, action := range turn.Actions {
			if action.ActionType == "switch" {
				switchCount++
			}
		}
	}

	if switchCount == 0 {
		t.Error("expected switches to be tracked")
	}
}

// ===== Key Moments Tests =====

func TestKeyMomentsKO(t *testing.T) {
	logWithFaint := `|j|☆Player1
|j|☆Player2
|player|p1|Player1|test|1500
|player|p2|Player2|test|1500
|tier|[Gen 9] VGC 2025 Reg H (Bo3)
|poke|p1|Pikachu, L50, M|
|poke|p2|Charizard, L50, M|
|teamsize|p1|1
|teamsize|p2|1
|start
|switch|p1a: Pikachu|Pikachu, L50, M|100/100
|switch|p2a: Charizard|Charizard, L50, M|100/100
|turn|1
|move|p1a: Pikachu|Thunderbolt|p2a: Charizard
|-damage|p2a: Charizard|0 fnt
|faint|p2a: Charizard
|upkeep
|win|Player1`

	summary, _ := ParseShowdownLog(logWithFaint)

	hasKO := false
	for _, moment := range summary.KeyMoments {
		if moment.Type == "KO" {
			hasKO = true
			break
		}
	}

	if !hasKO {
		t.Error("expected KO to be recorded as key moment")
	}
}

// ===== Stats Validation =====

func TestStatsConsistency(t *testing.T) {
	log := sampleBattleLog()
	summary, _ := ParseShowdownLog(log)

	// Total turns should match number of parsed turns
	if summary.Stats.TotalTurns != len(summary.Turns) {
		t.Errorf("expected %d total turns, stats show %d", len(summary.Turns), summary.Stats.TotalTurns)
	}

	// Player stats should be populated
	if summary.Stats.Player1Stats.MoveCount == 0 && len(summary.Turns) > 0 {
		// At least one move expected if there are turns
		t.Error("expected player1 move count > 0")
	}
}

// ===== UUID Uniqueness =====

func TestUUIDUniqueness(t *testing.T) {
	log := sampleBattleLog()

	ids := make(map[string]bool)
	for i := 0; i < 10; i++ {
		summary, _ := ParseShowdownLog(log)
		if ids[summary.ID] {
			t.Error("expected unique UUIDs for each parse")
		}
		ids[summary.ID] = true
	}
}

// ===== Test Fixtures =====

func logWithPositionProgression() string {
	return `|j|☆Player1
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
|switch|p1a: Pikachu|Pikachu, L50, M|100/100
|switch|p2a: Blastoise|Blastoise, L50, M|100/100
|turn|1
|move|p1a: Pikachu|Thunderbolt|p2a: Blastoise
|-supereffective|p2a: Blastoise
|-damage|p2a: Blastoise|30/100
|move|p2a: Blastoise|Hydro Pump|p1a: Pikachu
|-damage|p1a: Pikachu|50/100
|upkeep
|turn|2
|move|p1a: Pikachu|Thunder Wave|p2a: Blastoise
|-damage|p2a: Blastoise|20/100
|move|p2a: Blastoise|Ice Beam|p1a: Pikachu
|-supereffective|p1a: Pikachu
|-damage|p1a: Pikachu|0 fnt
|faint|p1a: Pikachu
|upkeep
|turn|3
|switch|p1a: Charizard|Charizard, L50, M|100/100
|move|p2a: Blastoise|Waterfall|p1a: Charizard
|-supereffective|p1a: Charizard
|-damage|p1a: Charizard|40/100
|upkeep
|turn|4
|move|p1a: Charizard|Flare Blitz|p2a: Blastoise
|-damage|p2a: Blastoise|0 fnt
|faint|p2a: Blastoise
|upkeep
|win|Player1`
}

func logWithSignificantMomentum() string {
	return `|j|☆Player1
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
|switch|p1a: Pikachu|Pikachu, L50, M|100/100
|switch|p2a: Blastoise|Blastoise, L50, M|100/100
|turn|1
|move|p1a: Pikachu|Thunderbolt|p2a: Blastoise
|-supereffective|p2a: Blastoise
|-damage|p2a: Blastoise|10/100
|move|p2a: Blastoise|Hydro Pump|p1a: Pikachu
|-damage|p1a: Pikachu|5/100
|upkeep
|turn|2
|move|p1a: Pikachu|Quick Attack|p2a: Blastoise
|-damage|p2a: Blastoise|0 fnt
|faint|p2a: Blastoise
|upkeep
|turn|3
|switch|p2a: Dragonite|Dragonite, L50, M|100/100
|move|p1a: Pikachu|Thunder Wave|p2a: Dragonite
|-damage|p2a: Dragonite|100/100
|move|p2a: Dragonite|Outrage|p1a: Pikachu
|-damage|p1a: Pikachu|0 fnt
|faint|p1a: Pikachu
|upkeep
|turn|4
|switch|p1a: Charizard|Charizard, L50, M|100/100
|move|p2a: Dragonite|Earthquake|p1a: Charizard
|-supereffective|p1a: Charizard
|-damage|p1a: Charizard|10/100
|upkeep
|win|Player2`
}
