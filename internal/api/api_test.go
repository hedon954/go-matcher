package api

import (
	"encoding/json"
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/pto"
)

func TestSaveEntries_ReloadEntries(t *testing.T) {
	mgrs := NewEntryManagers()
	api := &API{
		PM: mgrs.PlayerMgr,
		GM: mgrs.GroupMgr,
		TM: mgrs.TeamMgr,
		RM: mgrs.RoomMgr,
	}

	players, groups, teams, rooms := prepareEntries(t, mgrs)

	// save entries
	if err := api.SaveEntries(); err != nil {
		t.Fatalf("save entries failed: %v", err)
	}

	// reload entries
	mgrs.PlayerMgr.Clear()
	mgrs.GroupMgr.Clear()
	mgrs.TeamMgr.Clear()
	mgrs.RoomMgr.Clear()
	newApi := &API{
		PM: mgrs.PlayerMgr,
		GM: mgrs.GroupMgr,
		TM: mgrs.TeamMgr,
		RM: mgrs.RoomMgr,
	}
	if err := newApi.ReloadEntries(); err != nil {
		t.Fatalf("reload entries failed: %v", err)
	}

	// check entries
	comparePlayers(t, players, mgrs.PlayerMgr.All())
	compareGroups(t, groups, mgrs.GroupMgr.All())
	compareTeams(t, teams, mgrs.TeamMgr.All())
	compareRooms(t, rooms, mgrs.RoomMgr.All())
}

func prepareEntries(t *testing.T, mgrs *entry.Mgrs) (
	[]entry.Player, []entry.Group, []entry.Team, []entry.Room,
) {
	players := []entry.Player{}
	groups := []entry.Group{}
	teams := []entry.Team{}
	rooms := []entry.Room{}

	for i := 0; i < 10; i++ {
		gameMode := constant.GameModeGoatGame
		if i%2 == 0 {
			gameMode = constant.GameModeTest
		}
		p, err := mgrs.CreatePlayer(&pto.PlayerInfo{
			UID:         fmt.Sprintf("1-%d", i),
			GameMode:    gameMode,
			Glicko2Info: &pto.Glicko2Info{},
		})
		if err != nil {
			panic(err)
		}

		g, err := mgrs.CreateGroup(5, p)
		if err != nil {
			panic(err)
		}

		t, err := mgrs.CreateTeam(g)
		if err != nil {
			panic(err)
		}

		r, err := mgrs.CreateRoom(3, t)
		if err != nil {
			panic(err)
		}

		mgrs.TeamMgr.Add(t.ID(), t)
		mgrs.RoomMgr.Add(r.ID(), r)

		players = append(players, p)
		groups = append(groups, g)
		teams = append(teams, t)
		rooms = append(rooms, r)
	}

	assert.Equal(t, mgrs.PlayerMgr.Len(), len(players))
	assert.Equal(t, mgrs.GroupMgr.Len(), len(groups))
	assert.Equal(t, mgrs.TeamMgr.Len(), len(teams))
	assert.Equal(t, mgrs.RoomMgr.Len(), len(rooms))

	return players, groups, teams, rooms
}

//nolint:dupl
func comparePlayers(t *testing.T, a, b []entry.Player) {
	sort.Slice(a, func(i, j int) bool {
		return a[i].Base().UID() < a[j].Base().UID()
	})
	sort.Slice(b, func(i, j int) bool {
		return b[i].Base().UID() < b[j].Base().UID()
	})
	jsonA, _ := json.Marshal(a)
	jsonB, _ := json.Marshal(b)
	assert.JSONEq(t, string(jsonA), string(jsonB))
}

//nolint:dupl
func compareGroups(t *testing.T, a, b []entry.Group) {
	sort.Slice(a, func(i, j int) bool {
		return a[i].ID() < a[j].ID()
	})
	sort.Slice(b, func(i, j int) bool {
		return b[i].ID() < b[j].ID()
	})
	jsonA, _ := json.Marshal(a)
	jsonB, _ := json.Marshal(b)
	assert.JSONEq(t, string(jsonA), string(jsonB))
}

//nolint:dupl
func compareTeams(t *testing.T, a, b []entry.Team) {
	sort.Slice(a, func(i, j int) bool {
		return a[i].ID() < a[j].ID()
	})
	sort.Slice(b, func(i, j int) bool {
		return b[i].ID() < b[j].ID()
	})
	jsonA, _ := json.Marshal(a)
	jsonB, _ := json.Marshal(b)
	assert.JSONEq(t, string(jsonA), string(jsonB))
}

//nolint:dupl
func compareRooms(t *testing.T, a, b []entry.Room) {
	sort.Slice(a, func(i, j int) bool {
		return a[i].ID() < a[j].ID()
	})
	sort.Slice(b, func(i, j int) bool {
		return b[i].ID() < b[j].ID()
	})
	jsonA, _ := json.Marshal(a)
	jsonB, _ := json.Marshal(b)
	assert.JSONEq(t, string(jsonA), string(jsonB))
}
