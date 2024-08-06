# Entry Design

This repository contains the implementation of a matchmaking system designed to support different matching strategies and game modes. The project is organized into various components that facilitate the creation and management of player, group, team, and room entries using different matchmaking algorithms.

## Directory Structure Demo
```bash
.
├── entry
│ ├── glicko2
│ │ ├── group.go
│ │ ├── player.go
│ │ ├── room.go
│ │ └── team.go
│ ├── goat_game
│ │ ├── group.go
│ │ ├── player.go
│ │ ├── room.go
│ │ └── team.go
│ ├── group.go
│ ├── player.go
│ ├── room.go
│ └── team.go
└── repository
  ├── group.go
  ├── player.go
  ├── room.go
  └── team.go
```

## Components

![the relastionship between entries](../../assets/img/entry.png)

### Root Directory

The entry root directory provides interfaces and base implementations for the four types of entries:
- **Interfaces**:
    - `Player`
    - `Group`
    - `Team`
    - `Room`
- **Base Implementations**:
    - `PlayerBase`
    - `GroupBase`
    - `TeamBase`
    - `RoomBase`

These interfaces define the core functionalities for different entities, while the bases provide foundational implementations that can be combinated.

### Match Strategy (`glicko2`)

The `glicko2` directory represents a specific match strategy based on the `Glicko2` rating system. It includes implementations for the different types of entries, each combining the respective base class:

- `PlayerBaseGlicko2` (combines `PlayerBase`)
- `GroupBaseGlicko2` (combines `GroupBase`)
- `TeamBaseGlicko2` (combines `TeamBase`)
- `RoomBaseGlicko2` (combines `RoomBase`)

### Game Mode (`goat_game`)

The `goat_game` directory represents a specific game mode. A game mode can support multiple match strategies. Here, it combines the Glicko2 strategy implementations to form its own entries:

- `goat_game.Player` (combines `PlayerBaseGlicko2`, `PlayerBaseGather`)
- `goat_game.Group` (combines `GroupBaseGlicko2`, `GroupBaseGather`)
- `goat_game.Team` (combines `TeamBaseGlicko2`, `TeamBaseGather`)
- `goat_game.Room` (combines `RoomBaseGlicko2`, `RoomBaseGather`)

### Repository

In the `repository` directory, there are functions to create specific entry implementations based on provided parameters (game mode and match strategy). This allows for dynamic creation of entries like `Player`, `Group`, `Team`, and `Room`.

```go
type PlayerMgr struct {
	*collection.Manager[string, entry.Player]
}

type GroupMgr struct {
	*collection.Manager[int64, entry.Group]
}

type TeamMgr struct {
	*collection.Manager[int64, entry.Team]
	teamIDIter atomic.Int64
}

type RoomMgr struct {
	*collection.Manager[int64, entry.Room]
	roomIDIter atomic.Int64
}
```

## Factory Method

To use this matchmaking system, you can instantiate the desired base and strategy combinations using the factory method. Below is an example of how to create a player entry:

```go
// repository/player.go
func (m *PlayerMgr) CreatePlayer(pInfo *pto.PlayerInfo) (p entry.Player, err error) {
	base := entry.NewPlayerBase(pInfo)

	switch base.GameMode {
	case constant.GameModeGoatGame:
		p, err = goat_game.CreatePlayer(base, pInfo)
	case constant.GameModeTest:
		p = base
	default:
		return nil, fmt.Errorf("unsupported game mode: %d", base.GameMode)
	}

	if err != nil {
		return nil, err
	}
	m.Add(p.UID(), p)
	return p, nil
}

// entry/goat_game/player.go
func CreatePlayer(base *entry.PlayerBase, pInfo *pto.PlayerInfo) (entry.Player, error) {
	p := &Player{}

	if err := p.withMatchStrategy(base, pInfo.Glicko2Info); err != nil {
		return nil, err
	}
	return p, nil
}
func (p *Player) withMatchStrategy(base *entry.PlayerBase, info *pto.Glicko2Info) error {
	switch base.MatchStrategy {
	case constant.MatchStrategyGlicko2:
		p.PlayerBaseGlicko2 = glicko2.CreatePlayerBase(base, info)
	default:
		return fmt.Errorf("unknown match strategy: %d", base.MatchStrategy)
	}
	return nil
}
```
