-- Migration: Create BattleForge database schema
-- Version: 001_initial_schema.sql

-- Players table
CREATE TABLE IF NOT EXISTS players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Pokémon species reference table (static data)
CREATE TABLE IF NOT EXISTS pokemon_species (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    pokedex_number INT NOT NULL UNIQUE,
    type_1 VARCHAR(50) NOT NULL,
    type_2 VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Battles table
CREATE TABLE IF NOT EXISTS battles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    format VARCHAR(100) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    duration_sec INT NOT NULL,
    winner VARCHAR(20),
    player1_id VARCHAR(255) NOT NULL,
    player2_id VARCHAR(255) NOT NULL,
    battle_log TEXT,
    is_private BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_battles_format (format),
    INDEX idx_battles_timestamp (timestamp)
);

-- Pokémon instances (specific Pokémon in a battle)
CREATE TABLE IF NOT EXISTS pokemon (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    battle_id UUID NOT NULL REFERENCES battles(id) ON DELETE CASCADE,
    species_id INT NOT NULL REFERENCES pokemon_species(id),
    nickname VARCHAR(100),
    level INT NOT NULL,
    gender VARCHAR(1),
    ability VARCHAR(100),
    item VARCHAR(100),
    hp_base INT,
    attack_base INT,
    defense_base INT,
    sp_atk_base INT,
    sp_def_base INT,
    speed_base INT,
    shiny BOOLEAN DEFAULT FALSE,
    happiness INT DEFAULT 255,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_pokemon_battle (battle_id)
);

-- Moves reference table (static data)
CREATE TABLE IF NOT EXISTS moves (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    type VARCHAR(50) NOT NULL,
    category VARCHAR(20),
    power INT,
    accuracy INT,
    pp_max INT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Pokémon-Move mapping
CREATE TABLE IF NOT EXISTS pokemon_moves (
    id SERIAL PRIMARY KEY,
    pokemon_id UUID NOT NULL REFERENCES pokemon(id) ON DELETE CASCADE,
    move_id INT NOT NULL REFERENCES moves(id),
    slot INT NOT NULL,
    UNIQUE(pokemon_id, slot),
    INDEX idx_pokemon_moves_pokemon (pokemon_id)
);

-- Items reference table (static data)
CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    category VARCHAR(50),
    effect VARCHAR(255),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Battle turns
CREATE TABLE IF NOT EXISTS battle_turns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    battle_id UUID NOT NULL REFERENCES battles(id) ON DELETE CASCADE,
    turn_number INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(battle_id, turn_number),
    INDEX idx_battle_turns_battle (battle_id)
);

-- Battle actions (moves, switches, items)
CREATE TABLE IF NOT EXISTS battle_actions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    battle_turn_id UUID NOT NULL REFERENCES battle_turns(id) ON DELETE CASCADE,
    player_number INT NOT NULL,
    action_type VARCHAR(50) NOT NULL,
    move_id INT REFERENCES moves(id),
    pokemon_id UUID REFERENCES pokemon(id),
    item_id INT REFERENCES items(id),
    switch_to_pokemon_id UUID REFERENCES pokemon(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_battle_actions_turn (battle_turn_id)
);

-- Battle analysis results
CREATE TABLE IF NOT EXISTS battle_analysis (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    battle_id UUID NOT NULL UNIQUE REFERENCES battles(id) ON DELETE CASCADE,
    total_turns INT,
    avg_damage_per_turn NUMERIC(10, 2),
    avg_heal_per_turn NUMERIC(10, 2),
    moves_used_count INT,
    switches_count INT,
    super_effective_moves INT,
    not_very_effective_moves INT,
    critical_hits INT,
    player1_damage_dealt INT,
    player1_damage_taken INT,
    player1_healing_done INT,
    player2_damage_dealt INT,
    player2_damage_taken INT,
    player2_healing_done INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Key moments in the battle
CREATE TABLE IF NOT EXISTS key_moments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    battle_id UUID NOT NULL REFERENCES battles(id) ON DELETE CASCADE,
    turn_number INT NOT NULL,
    moment_type VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    significance INT DEFAULT 5,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_key_moments_battle (battle_id)
);

-- Create indexes for common queries
CREATE INDEX IF NOT EXISTS idx_battles_player1 ON battles(player1_id);
CREATE INDEX IF NOT EXISTS idx_battles_player2 ON battles(player2_id);
CREATE INDEX IF NOT EXISTS idx_battles_is_private ON battles(is_private);
