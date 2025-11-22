# Entity-Relationship Diagram & Data Schema

## ER Diagram

```
┌────────────────┐
│    battles     │
├────────────────┤
│ id (PK)        │◄──┐
│ format         │   │
│ timestamp      │   │
│ duration_sec   │   │
│ winner         │   │
│ player1_id (FK)├──┐│
│ player2_id (FK)├──┤│
│ created_at     │   ││
│ updated_at     │   ││
└────────────────┘   ││
         ▲            ││
         │            ││
    ┌────┴────────────┘│
    │                  │
    │   ┌──────────────┘
    │   │
    │   ▼
┌───┴───────────┐
│  battle_      │
│  turns        │
├───────────────┤
│ id (PK)       │
│ battle_id(FK) │◄──┐
│ turn_num      │   │
│ created_at    │   │
└───────────────┘   │
    ▲               │
    │               │
    ├───────────────┤
    │               │
    ▼               │
┌───────────────┐   │
│ battle_       │   │
│ actions       │   │
├───────────────┤   │
│ id (PK)       │   │
│ turn_id (FK)──┴───┘
│ player        │
│ action_type   │
│ move_id (FK)  │
│ pokemon_id(FK)│
│ item_id (FK)  │
│ created_at    │
└───────────────┘
    ▲     ▲     ▲
    │     │     │
    │     │     └──────┐
    │     │            │
┌───┴──┐  │      ┌─────┴─────┐
│moves │  │      │   items   │
├──────┤  │      ├───────────┤
│id(PK)│  │      │id (PK)    │
│name  │  │      │name       │
│type  │  │      │effect     │
│power │  │      │created_at │
│pp    │  │      └───────────┘
└──────┘  │
       ┌──┴──────┐
       │          │
       ▼          ▼
  ┌────────┐  ┌──────────────┐
  │pokemon │  │  pokemon_    │
  │        │  │  moves       │
  ├────────┤  ├──────────────┤
  │id (PK) │  │id (PK)       │
  │name    │  │pokemon_id(FK)│
  │species │  │move_id (FK)  │
  │level   │  │slot          │
  │ability │  └──────────────┘
  │item_id │
  │gender  │
  │shiny   │
  └────────┘

┌──────────────┐
│players       │
├──────────────┤
│id (PK)       │
│username      │
│created_at    │
│updated_at    │
└──────────────┘

┌──────────────────┐
│battle_analysis   │
├──────────────────┤
│id (PK)           │
│battle_id (FK)    │
│total_turns       │
│avg_damage_turn   │
│avg_heal_turn     │
│moves_used_count  │
│switches_count    │
│super_effective   │
│not_very_effective│
│critical_hits     │
│created_at        │
└──────────────────┘

┌──────────────────┐
│key_moments       │
├──────────────────┤
│id (PK)           │
│battle_id (FK)    │
│turn_number       │
│moment_type       │
│description       │
│significance      │
│created_at        │
└──────────────────┘
```

## Schema Definition

```sql
-- Players table
CREATE TABLE players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Battles table
CREATE TABLE battles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    format VARCHAR(100) NOT NULL,                    -- e.g., "Regulation H"
    timestamp TIMESTAMP NOT NULL,
    duration_sec INT NOT NULL,
    winner VARCHAR(20),                              -- "player1", "player2", "draw"
    player1_id UUID NOT NULL REFERENCES players(id),
    player2_id UUID NOT NULL REFERENCES players(id),
    battle_log TEXT,                                 -- Original .log content
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_battles_format (format),
    INDEX idx_battles_timestamp (timestamp)
);

-- Pokémon species reference (static data)
CREATE TABLE pokemon_species (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    pokedex_number INT NOT NULL UNIQUE,
    type_1 VARCHAR(50) NOT NULL,
    type_2 VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Pokémon instances (specific Pokémon in a battle)
CREATE TABLE pokemon (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    battle_id UUID NOT NULL REFERENCES battles(id) ON DELETE CASCADE,
    species_id INT NOT NULL REFERENCES pokemon_species(id),
    nickname VARCHAR(100),
    level INT NOT NULL,
    gender VARCHAR(1),                               -- 'M', 'F', ''
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
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Moves reference (static data)
CREATE TABLE moves (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    type VARCHAR(50) NOT NULL,
    category VARCHAR(20),                           -- "physical", "special", "status"
    power INT,
    accuracy INT,
    pp_max INT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Pokémon-Move mapping
CREATE TABLE pokemon_moves (
    id SERIAL PRIMARY KEY,
    pokemon_id UUID NOT NULL REFERENCES pokemon(id) ON DELETE CASCADE,
    move_id INT NOT NULL REFERENCES moves(id),
    slot INT NOT NULL,                              -- 1-4
    UNIQUE(pokemon_id, slot)
);

-- Items reference (static data)
CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    category VARCHAR(50),                           -- "held", "consumable", etc.
    effect VARCHAR(255),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Battle turns
CREATE TABLE battle_turns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    battle_id UUID NOT NULL REFERENCES battles(id) ON DELETE CASCADE,
    turn_number INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(battle_id, turn_number)
);

-- Battle actions (moves, switches, items)
CREATE TABLE battle_actions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    battle_turn_id UUID NOT NULL REFERENCES battle_turns(id) ON DELETE CASCADE,
    player_number INT NOT NULL,                     -- 1 or 2
    action_type VARCHAR(50) NOT NULL,               -- "move", "switch", "item"
    move_id INT REFERENCES moves(id),
    pokemon_id UUID REFERENCES pokemon(id),
    item_id INT REFERENCES items(id),
    switch_to_pokemon_id UUID REFERENCES pokemon(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Battle analysis results
CREATE TABLE battle_analysis (
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
CREATE TABLE key_moments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    battle_id UUID NOT NULL REFERENCES battles(id) ON DELETE CASCADE,
    turn_number INT NOT NULL,
    moment_type VARCHAR(50) NOT NULL,               -- "switch", "ko", "status", "weather", "critical"
    description TEXT NOT NULL,
    significance INT DEFAULT 5,                     -- 1-10 scale
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_key_moments_battle (battle_id)
);

-- Indexes for common queries
CREATE INDEX idx_battle_turns_battle ON battle_turns(battle_id);
CREATE INDEX idx_battle_actions_turn ON battle_actions(battle_turn_id);
CREATE INDEX idx_pokemon_battle ON pokemon(battle_id);
CREATE INDEX idx_pokemon_moves_pokemon ON pokemon_moves(pokemon_id);
```

## Data Relationships

| Table | Primary Key | Foreign Keys | Purpose |
|-------|-------------|--------------|---------|
| **players** | id (UUID) | — | User accounts |
| **battles** | id (UUID) | player1_id, player2_id | Battle records |
| **pokemon** | id (UUID) | battle_id, species_id | Pokémon instances in battles |
| **pokemon_species** | id (INT) | — | Static Pokémon data |
| **moves** | id (INT) | — | Static move data |
| **pokemon_moves** | id (INT) | pokemon_id, move_id | Many-to-many: which moves a Pokémon has |
| **items** | id (INT) | — | Static item data |
| **battle_turns** | id (UUID) | battle_id | Turn records |
| **battle_actions** | id (UUID) | battle_turn_id, move_id, pokemon_id, item_id | Individual actions within turns |
| **battle_analysis** | id (UUID) | battle_id (unique) | Aggregate stats for a battle |
| **key_moments** | id (UUID) | battle_id | Notable moments in the battle |

## Normalization Notes

- **pokemon_species**, **moves**, and **items** are reference tables (static data from Pokédex/game data)
- **pokemon** instances are battle-specific (same species can appear in multiple battles)
- **battle_analysis** is denormalized for performance (could be computed from turns/actions, but stored for quick retrieval)
- **key_moments** is denormalized (derived from actions, but stored for quick retrieval)

## Query Patterns

```sql
-- Get a complete battle with all turns and actions
SELECT b.*, bt.turn_number, ba.action_type, m.name as move_name
FROM battles b
JOIN battle_turns bt ON b.id = bt.battle_id
JOIN battle_actions ba ON bt.id = ba.battle_turn_id
LEFT JOIN moves m ON ba.move_id = m.id
WHERE b.id = ?
ORDER BY bt.turn_number, ba.created_at;

-- Get key moments for a battle
SELECT * FROM key_moments
WHERE battle_id = ?
ORDER BY turn_number;

-- Get aggregate stats for a battle
SELECT * FROM battle_analysis
WHERE battle_id = ?;

-- Get recent battles by a player
SELECT b.* FROM battles b
WHERE b.player1_id = ? OR b.player2_id = ?
ORDER BY b.timestamp DESC
LIMIT 10;
```
