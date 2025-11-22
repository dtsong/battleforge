package httpapi

// sampleShowdownLog returns a minimal valid Showdown battle log for testing.
func sampleShowdownLog() string {
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
|switch|p1a: Charizard|Charizard, L50, M|100/100
|move|p2a: Blastoise|Protect|p2a: Blastoise
|-singleturn|p2a: Blastoise|Protect
|upkeep
|turn|3
|move|p1a: Charizard|Flamethrower|p2a: Blastoise
|-resisted|p2a: Blastoise
|-damage|p2a: Blastoise|45/100
|move|p2a: Blastoise|Ice Beam|p1a: Charizard
|-supereffective|p1a: Charizard
|-damage|p1a: Charizard|20/100
|faint|p1a: Charizard
|upkeep
|
|switch|p1a: Pikachu|Pikachu, L50, M|30/100
|turn|4
|move|p1a: Pikachu|Quick Attack|p2a: Blastoise
|-damage|p2a: Blastoise|40/100
|move|p2a: Blastoise|Waterfall|p1a: Pikachu
|-supereffective|p1a: Pikachu
|-damage|p1a: Pikachu|0 fnt
|faint|p1a: Pikachu
|upkeep
|
|win|Player2`
}
