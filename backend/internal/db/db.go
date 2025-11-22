package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Database wraps a SQL database connection with helper methods.
type Database struct {
	conn *sql.DB
}

// NewDatabase creates a new Database instance.
func NewDatabase(connString string) (*Database, error) {
	conn, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{conn: conn}, nil
}

// Close closes the database connection.
func (db *Database) Close() error {
	return db.conn.Close()
}

// WithTx executes a function within a database transaction.
func (db *Database) WithTx(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback() // Ignore rollback error as we're already returning an error
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Exec executes a query without returning rows.
func (db *Database) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := db.conn.ExecContext(ctx, query, args...)
	return err
}

// QueryRow queries a single row.
func (db *Database) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return db.conn.QueryRowContext(ctx, query, args...)
}

// Query queries multiple rows.
func (db *Database) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return db.conn.QueryContext(ctx, query, args...)
}

// StoreBattle saves a battle and related data to the database.
// Returns the battle ID.
func (db *Database) StoreBattle(ctx context.Context, battle *Battle) (string, error) {
	var battleID string

	err := db.WithTx(ctx, func(tx *sql.Tx) error {
		// Insert battle
		err := tx.QueryRowContext(ctx,
			`INSERT INTO battles (format, timestamp, duration_sec, winner, player1_id, player2_id, battle_log, is_private, created_at, updated_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
			 RETURNING id`,
			battle.Format, battle.Timestamp, battle.DurationSec, battle.Winner,
			battle.Player1ID, battle.Player2ID, battle.BattleLog, battle.IsPrivate,
		).Scan(&battleID)

		if err != nil {
			return fmt.Errorf("failed to insert battle: %w", err)
		}

		// Insert analysis results
		if battle.Analysis != nil {
			err = insertBattleAnalysis(ctx, tx, battleID, battle.Analysis)
			if err != nil {
				return err
			}
		}

		// Insert key moments
		for _, moment := range battle.KeyMoments {
			err = insertKeyMoment(ctx, tx, battleID, moment)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return battleID, err
}

// GetBattle retrieves a battle by ID.
func (db *Database) GetBattle(ctx context.Context, battleID string) (*Battle, error) {
	var b Battle
	err := db.QueryRow(ctx,
		`SELECT id, format, timestamp, duration_sec, winner, player1_id, player2_id, battle_log, is_private, created_at, updated_at
		 FROM battles WHERE id = $1`,
		battleID,
	).Scan(&b.ID, &b.Format, &b.Timestamp, &b.DurationSec, &b.Winner, &b.Player1ID, &b.Player2ID, &b.BattleLog, &b.IsPrivate, &b.CreatedAt, &b.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Get analysis data
	analysis, err := getBattleAnalysis(ctx, db, battleID)
	if err != nil {
		return nil, err
	}
	b.Analysis = analysis

	// Get key moments
	moments, err := getKeyMoments(ctx, db, battleID)
	if err != nil {
		return nil, err
	}
	b.KeyMoments = moments

	return &b, nil
}

// ListBattles retrieves battles with optional filtering.
func (db *Database) ListBattles(ctx context.Context, filter *BattleFilter, limit int, offset int) ([]*Battle, int, error) {
	query := `SELECT id, format, timestamp, duration_sec, winner, player1_id, player2_id, is_private FROM battles WHERE 1=1`
	var args []interface{}
	argIndex := 1

	if filter != nil {
		if filter.Format != "" {
			query += fmt.Sprintf(" AND format = $%d", argIndex)
			args = append(args, filter.Format)
			argIndex++
		}
		if filter.IsPrivate != nil {
			query += fmt.Sprintf(" AND is_private = $%d", argIndex)
			args = append(args, *filter.IsPrivate)
			argIndex++
		}
	}

	// Get total count
	countQuery := "SELECT COUNT(*) FROM battles WHERE 1=1"
	if len(args) > 0 {
		countQuery = query
		countQuery = "SELECT COUNT(*) FROM (" + countQuery + ") AS filtered"
	}

	var total int
	err := db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query += fmt.Sprintf(" ORDER BY timestamp DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var battles []*Battle
	for rows.Next() {
		var b Battle
		err := rows.Scan(&b.ID, &b.Format, &b.Timestamp, &b.DurationSec, &b.Winner, &b.Player1ID, &b.Player2ID, &b.IsPrivate)
		if err != nil {
			return nil, 0, err
		}
		battles = append(battles, &b)
	}

	return battles, total, rows.Err()
}

// Helper functions

func insertBattleAnalysis(ctx context.Context, tx *sql.Tx, battleID string, analysis *BattleAnalysis) error {
	_, err := tx.ExecContext(ctx,
		`INSERT INTO battle_analysis (battle_id, total_turns, avg_damage_per_turn, avg_heal_per_turn, moves_used_count, switches_count, super_effective_moves, not_very_effective_moves, critical_hits, player1_damage_dealt, player1_damage_taken, player1_healing_done, player2_damage_dealt, player2_damage_taken, player2_healing_done, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, NOW())`,
		battleID, analysis.TotalTurns, analysis.AvgDamagePerTurn, analysis.AvgHealPerTurn,
		analysis.MovesUsedCount, analysis.SwitchesCount, analysis.SuperEffectiveMoves,
		analysis.NotVeryEffectiveMoves, analysis.CriticalHits, analysis.Player1DamageDealt,
		analysis.Player1DamageTaken, analysis.Player1HealingDone, analysis.Player2DamageDealt,
		analysis.Player2DamageTaken, analysis.Player2HealingDone,
	)
	return err
}

func insertKeyMoment(ctx context.Context, tx *sql.Tx, battleID string, moment *KeyMoment) error {
	_, err := tx.ExecContext(ctx,
		`INSERT INTO key_moments (battle_id, turn_number, moment_type, description, significance, created_at)
		 VALUES ($1, $2, $3, $4, $5, NOW())`,
		battleID, moment.TurnNumber, moment.MomentType, moment.Description, moment.Significance,
	)
	return err
}

func getBattleAnalysis(ctx context.Context, db *Database, battleID string) (*BattleAnalysis, error) {
	var analysis BattleAnalysis
	err := db.QueryRow(ctx,
		`SELECT battle_id, total_turns, avg_damage_per_turn, avg_heal_per_turn, moves_used_count, switches_count, super_effective_moves, not_very_effective_moves, critical_hits, player1_damage_dealt, player1_damage_taken, player1_healing_done, player2_damage_dealt, player2_damage_taken, player2_healing_done
		 FROM battle_analysis WHERE battle_id = $1`,
		battleID,
	).Scan(&analysis.BattleID, &analysis.TotalTurns, &analysis.AvgDamagePerTurn, &analysis.AvgHealPerTurn,
		&analysis.MovesUsedCount, &analysis.SwitchesCount, &analysis.SuperEffectiveMoves,
		&analysis.NotVeryEffectiveMoves, &analysis.CriticalHits, &analysis.Player1DamageDealt,
		&analysis.Player1DamageTaken, &analysis.Player1HealingDone, &analysis.Player2DamageDealt,
		&analysis.Player2DamageTaken, &analysis.Player2HealingDone)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &analysis, nil
}

func getKeyMoments(ctx context.Context, db *Database, battleID string) ([]*KeyMoment, error) {
	rows, err := db.Query(ctx,
		`SELECT turn_number, moment_type, description, significance FROM key_moments WHERE battle_id = $1 ORDER BY turn_number`,
		battleID,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var moments []*KeyMoment
	for rows.Next() {
		var m KeyMoment
		err := rows.Scan(&m.TurnNumber, &m.MomentType, &m.Description, &m.Significance)
		if err != nil {
			return nil, err
		}
		moments = append(moments, &m)
	}

	return moments, rows.Err()
}
