package example

import (
	"fmt"
	"sync"
	"time"

	"github.com/hedon954/go-matcher/pkg/algorithm/glicko2"
	"github.com/montanaflynn/stats"
)

const (
	// 车队方差阈值
	MaliciousTeamVarianceMin  = 100000
	UnfriendlyTeamVarianceMin = 10000
)

type Group struct {
	sync.RWMutex

	ID         string              `json:"id"`
	State      glicko2.GroupState  `json:"state"`
	PlayersMap map[string]struct{} `json:"players_map"`
	Players    []*Player           `json:"players"`

	modeVersion int
	platform    int

	startMatchTimeSec int64
}

func (g *Group) IsNewer() bool {
	return false
}

func (g *Group) GetModeVersion() int {
	return g.modeVersion
}

func (g *Group) GetPlatform() int {
	return g.platform
}

func (g *Group) ForceCancelMatch(reason string, waitSec int64) {}

func NewGroup(id string, players []*Player) *Group {
	g := &Group{
		RWMutex:    sync.RWMutex{},
		ID:         id,
		State:      glicko2.GroupStateUnready,
		PlayersMap: make(map[string]struct{}),
		Players:    players,
	}
	for _, p := range g.Players {
		g.PlayersMap[p.GetID()] = struct{}{}
		g.startMatchTimeSec = p.GetStartMatchTimeSec()
	}
	return g
}

func (g *Group) GetID() string {
	return g.ID
}

func (g *Group) GetState() glicko2.GroupState {
	g.RLock()
	defer g.RUnlock()

	return g.State
}

func (g *Group) SetState(state glicko2.GroupState) {
	g.Lock()
	defer g.Unlock()

	g.State = state
}

func (g *Group) PlayerCount() int {
	g.RLock()
	defer g.RUnlock()
	return len(g.Players)
}

func (g *Group) GetPlayers() []glicko2.Player {
	g.RLock()
	defer g.RUnlock()
	res := make([]glicko2.Player, len(g.Players))
	for i := 0; i < len(res); i++ {
		res[i] = g.Players[i]
	}
	return res
}

func (g *Group) AddPlayers(players ...glicko2.Player) {
	g.Lock()
	defer g.Unlock()

	for _, p := range players {
		_, ok := g.PlayersMap[p.GetID()]
		if ok {
			continue
		}
		g.PlayersMap[p.GetID()] = struct{}{}
		g.Players = append(g.Players, p.(*Player))
	}
}

func (g *Group) RemovePlayers(players ...glicko2.Player) {
	g.Lock()
	defer g.Unlock()

	for _, p := range players {
		for i, gp := range g.Players {
			if gp == p {
				g.Players = append(g.Players[:i], g.Players[i+1:]...)
				delete(g.PlayersMap, p.GetID())
			}
		}
	}
}

// AverageMMR 算出队伍的平均 MMR
func (g *Group) AverageMMR() float64 {
	total := 0.0
	for _, player := range g.Players {
		total += player.GetMMR()
	}
	return total / float64(len(g.Players))
}

// MMR 算出队伍的最大的 MMR
func (g *Group) BiggestMMR() float64 {
	mmr := 0.0
	for _, p := range g.Players {
		pMMR := p.GetMMR()
		if pMMR > mmr {
			mmr = pMMR
		}
	}
	return mmr
}

// MMR 算出队伍的 MMR
func (g *Group) GetMMR() float64 {
	teamType := g.Type()
	switch teamType {
	case glicko2.GroupTypeUnfriendlyTeam:
		mmr := g.AverageMMR() * 1.5
		bMmr := g.BiggestMMR()
		if mmr > bMmr {
			mmr = bMmr
		}
		return mmr
	case glicko2.GroupTypeMaliciousTeam:
		return g.BiggestMMR()
	default:
		return g.AverageMMR()
	}
}

// Rank 队伍段位要弄平均值替代
func (g *Group) GetStar() int {
	if len(g.Players) == 0 {
		return 0
	}
	rank := 0
	for _, p := range g.Players {
		rank += p.GetStar()
	}
	return rank / len(g.Players)
}

// Group 算出队伍的 MMR 方差
func (g *Group) MMRVariance() float64 {
	data := stats.Float64Data{}
	for _, p := range g.Players {
		data = append(data, p.GetMMR())
	}
	variance, _ := stats.Variance(data)
	return variance
}

// Type 确定车队类型
func (g *Group) Type() glicko2.GroupType {
	if len(g.Players) != 5 {
		return glicko2.GroupTypeNotTeam
	}
	variance := g.MMRVariance()
	if variance >= MaliciousTeamVarianceMin {
		return glicko2.GroupTypeMaliciousTeam
	} else if variance >= UnfriendlyTeamVarianceMin {
		return glicko2.GroupTypeUnfriendlyTeam
	} else {
		return glicko2.GroupTypeNormalTeam
	}
}

func (g *Group) CanFillAi() bool {
	return false
	// TODO: 读取配置，根据条件判断是否可以返回 isAi
	now := time.Now().Unix()
	if now-g.GetStartMatchTimeSec() > 60 {
		return true
	}
	return false
}

// Print 打印 group 信息
func (g *Group) Print() {
	fmt.Printf("\t\t%s\t\t\t%d\t\t%.2f\t\t%ds\t\t\n", g.GetID(), len(g.Players), g.GetMMR(),
		time.Now().Unix()-g.GetStartMatchTimeSec())
}

func (g *Group) GetFinishMatchTimeSec() int64 {
	if len(g.Players) == 0 {
		return 0
	}
	return g.Players[0].GetFinishMatchTimeSec()
}

func (g *Group) SetFinishMatchTimeSec(t int64) {
	for _, p := range g.Players {
		p.SetFinishMatchTimeSec(t)
	}
}

func (g *Group) GetStartMatchTimeSec() int64 {
	return g.startMatchTimeSec
}

func (g *Group) SetStartMatchTimeSec(t int64) {
	g.startMatchTimeSec = t
	for _, p := range g.Players {
		p.SetStartMatchTimeSec(t)
	}
}
