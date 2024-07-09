package player

import (
	"fmt"
	"sync"

	"matcher/common"
	"matcher/enum"
	"matcher/internal/gamemode/goatgame"
	"matcher/pto"
)

type Manager struct {
	sync.RWMutex
	players map[string]common.Player
}

func NewManager() *Manager {
	const capacity = 1024
	return &Manager{players: make(map[string]common.Player, capacity)}
}

func (m *Manager) CreatePlayer(pInfo *pto.PlayerInfo) (common.Player, error) {
	playerBase := common.NewBase(pInfo)

	var p common.Player
	var err error

	switch playerBase.GameMode {
	case enum.GameModeGoat:
		if p, err = goatgame.CreatePlayer(playerBase); err != nil {
			break
		}
	default:
		return nil, fmt.Errorf("unsupported game mode: %d", playerBase.GameMode)
	}

	if err != nil {
		return nil, err
	}

	m.AddPlayer(p)
	return p, nil
}

func (m *Manager) GetPlayer(uid string) common.Player {
	m.RLock()
	defer m.RUnlock()
	return m.players[uid]
}

func (m *Manager) PlayerExists(uid string) bool {
	m.RLock()
	defer m.RUnlock()
	return m.players[uid] != nil
}

func (m *Manager) AddPlayer(player common.Player) {
	m.Lock()
	defer m.Unlock()
	m.players[player.UID()] = player
}

func (m *Manager) DelPlayer(uid string) {
	m.Lock()
	defer m.Unlock()
	delete(m.players, uid)
}
