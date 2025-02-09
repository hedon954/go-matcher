package api

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hibiken/asynq"

	"github.com/hedon954/go-matcher/internal/config"
	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/entry/goat_game"
	"github.com/hedon954/go-matcher/internal/entry/modes"
	"github.com/hedon954/go-matcher/internal/entry/test_game"
	"github.com/hedon954/go-matcher/internal/log"
	"github.com/hedon954/go-matcher/internal/matcher"
	"github.com/hedon954/go-matcher/internal/matcher/common"
	"github.com/hedon954/go-matcher/internal/matcher/glicko2"
	"github.com/hedon954/go-matcher/internal/service"
	"github.com/hedon954/go-matcher/internal/service/matchimpl"
	"github.com/hedon954/go-matcher/pkg/timer"
	timerasynq "github.com/hedon954/go-matcher/pkg/timer/asynq"
	timernative "github.com/hedon954/go-matcher/pkg/timer/native"
)

func init() {
	modes.Init()
}

type API struct {
	MS service.Match
	M  *matcher.Matcher
	PM *entry.PlayerMgr
	GM *entry.GroupMgr
	TM *entry.TeamMgr
	RM *entry.RoomMgr
}

// Start initializes the api components and starts them.
func Start(
	sc config.Configer[config.ServerConfig],
	mc config.Configer[config.MatchConfig],
) (api *API, shutdown func()) {
	var (
		groupChannel = make(chan entry.Group, 1024)
		roomChannel  = make(chan common.Result, 1024)
		mgrs         = NewEntryManagers()
	)

	matchConf := mc.Get()
	serverConf := sc.Get()

	// init delay timer
	dt, err := NewDelayTime(matchConf.DelayTimerType, serverConf.AsynqRedis)
	if err != nil {
		panic(err)
	}

	// init api
	api = NewAPI(mc, groupChannel, roomChannel, dt, NewGlicko2Matcher(roomChannel, matchConf, mgrs), mgrs)

	// if not in testing mode, reload entries
	// TODO: find a better way.
	if flag.Lookup("test.v") == nil {
		if err := api.ReloadEntries(); err != nil {
			panic(fmt.Errorf("failed to reload entries: %v", err))
		}
	}

	// start delay timer and match service
	go dt.Start()
	go api.M.Start()

	return api, func() {
		dt.Stop()
		api.M.Stop()

		// TODO: find a better way.
		if flag.Lookup("test.v") == nil {
			if err := api.SaveEntries(); err != nil {
				log.Error().Err(err).Msg("failed to save entries")
			}
		}
	}
}

func NewAPI(configer config.Configer[config.MatchConfig],
	groupChannel chan entry.Group, roomChannel chan common.Result,
	dt timer.Operator[int64], gm *glicko2.Matcher, mgrs *entry.Mgrs) *API {
	api := &API{
		PM: mgrs.PlayerMgr,
		GM: mgrs.GroupMgr,
		TM: mgrs.TeamMgr,
		RM: mgrs.RoomMgr,
		M:  matcher.New(groupChannel, gm),
		MS: matchimpl.NewDefault(configer, mgrs, groupChannel, roomChannel, dt),
	}
	return api
}

func NewDelayTime(t config.DelayTimerType, r *config.RedisOpt) (timer.Operator[int64], error) {
	switch t {
	case config.DelayTimerTypeAsynq:
		return timerasynq.NewTimer[int64](&asynq.RedisClientOpt{
			Addr:     r.Addr,
			Password: r.Password,
			DB:       r.DB,
		}), nil
	case config.DelayTimerTypeNative:
		return timernative.NewTimer(), nil
	default:
		return nil, fmt.Errorf("unsupported delay timer type: %s", t)
	}
}

func NewEntryManagers() *entry.Mgrs {
	return &entry.Mgrs{
		PlayerMgr: entry.NewPlayerMgr(),
		GroupMgr:  entry.NewGroupMgr(0),
		TeamMgr:   entry.NewTeamMgr(0),
		RoomMgr:   entry.NewRoomMgr(0),
	}
}

func NewGlicko2Matcher(roomChannel chan common.Result, conf *config.MatchConfig, mgrs *entry.Mgrs) *glicko2.Matcher {
	return glicko2.New(roomChannel, conf, conf.MatchInterval(), mgrs)
}

// SaveEntries saves the entries when the server stops.
func (api *API) SaveEntries() error {
	return save(map[string]any{
		"rooms.json":   api.RM.Encode(),
		"teams.json":   api.TM.Encode(),
		"groups.json":  api.GM.Encode(),
		"players.json": api.PM.Encode(),
	})
}

func save(data map[string]any) error {
	dir := filepath.Join(".", "tmp_entries")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	for filename, v := range data {
		bs, err := json.Marshal(v)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dir, filename), bs, 0600); err != nil {
			return err
		}
	}
	return nil
}

// ReloadEntries reloads the entries when the server starts.
func (api *API) ReloadEntries() error {
	// read files
	dir := filepath.Join(".", "tmp_entries")

	playerData := getDatFromFile(filepath.Join(dir, "players.json"))
	if err := api.reloadPlayers(playerData); err != nil {
		return err
	}

	groupData := getDatFromFile(filepath.Join(dir, "groups.json"))
	if err := api.reloadGroups(groupData); err != nil {
		return err
	}

	teamData := getDatFromFile(filepath.Join(dir, "teams.json"))
	if err := api.reloadTeams(teamData); err != nil {
		return err
	}

	roomData := getDatFromFile(filepath.Join(dir, "rooms.json"))
	if err := api.reloadRooms(roomData); err != nil {
		return err
	}

	// delete old backup files
	backupDir := filepath.Join(".", "tmp_entries", "matcher_backup")
	if err := os.RemoveAll(backupDir); err != nil {
		return err
	}

	// move current files to backup
	backupDir = filepath.Join(".", "tmp_entries", "matcher_backup")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return err
	}

	// move files to backup
	files := []string{"rooms.json", "teams.json", "groups.json", "players.json"}
	for _, file := range files {
		_ = os.Rename(filepath.Join(dir, file), filepath.Join(backupDir, file))
	}
	return nil
}

func getDatFromFile(file string) map[constant.GameMode][][]byte {
	data := make(map[constant.GameMode][][]byte)
	bs, err := os.ReadFile(file)
	if err != nil {
		return nil
	}
	if err := json.Unmarshal(bs, &data); err != nil {
		return nil
	}
	return data
}

//nolint:dupl
func (api *API) reloadPlayers(playerData map[constant.GameMode][][]byte) error {
	for mode, ps := range playerData {
		for _, pbs := range ps {
			var p entry.Player
			if mode == constant.GameModeTest {
				p = &test_game.Player{}
			} else if mode == constant.GameModeGoatGame {
				p = &goat_game.Player{}
			} else {
				log.Error().Any("mode", mode).Msg("unsupported game mode")
				continue
			}
			if err := p.Decode(pbs); err != nil {
				return err
			}
			api.PM.Add(p.Base().UID(), p)
		}
	}
	return nil
}

//nolint:dupl
func (api *API) reloadGroups(groupData map[constant.GameMode][][]byte) error {
	for mode, gs := range groupData {
		for _, gbs := range gs {
			var g entry.Group
			if mode == constant.GameModeTest {
				g = &test_game.Group{}
			} else if mode == constant.GameModeGoatGame {
				g = &goat_game.Group{}
			} else {
				log.Error().Any("mode", mode).Msg("unsupported game mode")
				continue
			}
			if err := g.Decode(gbs); err != nil {
				return err
			}
			if goatGroup, ok := g.(*goat_game.Group); ok {
				goatGroup.SetPlayerMgr(api.PM)
			}
			api.GM.Add(g.ID(), g)
		}
	}
	return nil
}

//nolint:dupl
func (api *API) reloadTeams(teamData map[constant.GameMode][][]byte) error {
	for mode, ts := range teamData {
		for _, tbs := range ts {
			var t entry.Team
			if mode == constant.GameModeTest {
				t = &test_game.Team{}
			} else if mode == constant.GameModeGoatGame {
				t = &goat_game.Team{}
			} else {
				log.Error().Any("mode", mode).Msg("unsupported game mode")
				continue
			}
			if err := t.Decode(tbs); err != nil {
				return err
			}
			if goatTeam, ok := t.(*goat_game.Team); ok {
				goatTeam.SetGroupMgr(api.GM)
			}
			api.TM.Add(t.ID(), t)
		}
	}
	return nil
}

//nolint:dupl
func (api *API) reloadRooms(roomData map[constant.GameMode][][]byte) error {
	for mode, rs := range roomData {
		for _, rbs := range rs {
			var r entry.Room
			if mode == constant.GameModeTest {
				r = &test_game.Room{}
			} else if mode == constant.GameModeGoatGame {
				r = &goat_game.Room{}
			} else {
				log.Error().Any("mode", mode).Msg("unsupported game mode")
				continue
			}
			if err := r.Decode(rbs); err != nil {
				return err
			}
			if goatRoom, ok := r.(*goat_game.Room); ok {
				goatRoom.SetTeamMgr(api.TM)
				goatRoom.FillGlicko2Teams()
			}
			api.RM.Add(r.ID(), r)
		}
	}
	return nil
}
