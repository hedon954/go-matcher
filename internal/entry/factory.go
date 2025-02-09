package entry

import (
	"fmt"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/pto"
)

type Factory interface {
	CreateRoom(mgr *Mgrs, base *RoomBase) (Room, error)
	CreatePlayer(mgr *Mgrs, base *PlayerBase, pInfo *pto.PlayerInfo) (Player, error)
	CreateTeam(mgr *Mgrs, base *TeamBase) (Team, error)
	CreateGroup(mgr *Mgrs, base *GroupBase) (Group, error)
}

var factories = make(map[constant.GameMode]Factory)

func RegisterFactory(gameMode constant.GameMode, factory Factory) {
	factories[gameMode] = factory
}

func GetFactory(gameMode constant.GameMode) Factory {
	return factories[gameMode]
}

func PrintModes() {
	if len(factories) == 0 {
		fmt.Println("┌─────────────────────────────┐")
		fmt.Println("│ No game modes registered    │")
		fmt.Println("└─────────────────────────────┘")
		return
	}

	fmt.Println("┌─────────────────────────────────┐")
	fmt.Println("│        Registered Modes         │")
	fmt.Println("├──────────┬────────────────────┤")
	fmt.Println("│    ID    │        Name        │")
	fmt.Println("├──────────┼────────────────────┤")
	for gameMode := range factories {
		fmt.Printf("│  %-6d  │  %-15s │\n", gameMode, constant.GameModeNames[gameMode])
	}
	fmt.Println("└──────────┴────────────────────┘")
}
