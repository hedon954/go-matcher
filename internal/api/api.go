package api

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hedon954/go-matcher/internal/config"
	"github.com/hedon954/go-matcher/internal/constant"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/entry/modes"
	"github.com/hedon954/go-matcher/internal/matcher"
	"github.com/hedon954/go-matcher/internal/matcher/common"
	"github.com/hedon954/go-matcher/internal/matcher/glicko2"
	"github.com/hedon954/go-matcher/internal/service"
	"github.com/hedon954/go-matcher/internal/service/matchimpl"
	"github.com/hedon954/go-matcher/pkg/timer"

	timerasynq "github.com/hedon954/go-matcher/pkg/timer/asynq"
	timernative "github.com/hedon954/go-matcher/pkg/timer/native"

	"github.com/hibiken/asynq"
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

	// start delay timer and match service
	go dt.Start()
	go api.M.Start()

	return api, func() {
		dt.Stop()
		api.M.Stop()
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
	rooms := make(map[constant.GameMode][][]byte)
	teams := make(map[constant.GameMode][][]byte)
	groups := make(map[constant.GameMode][][]byte)
	players := make(map[constant.GameMode][][]byte)

	api.RM.Range(func(id int64, r entry.Room) bool {
		bs, err := r.Encode()
		if err != nil {
			return true
		}
		rooms[r.GetMatchInfo().GameMode] = append(rooms[r.GetMatchInfo().GameMode], bs)
		return true
	})

	api.TM.Range(func(id int64, t entry.Team) bool {
		bs, err := t.Encode()
		if err != nil {
			return true
		}
		teams[t.Base().GameMode] = append(teams[t.Base().GameMode], bs)
		return true
	})

	api.GM.Range(func(id int64, g entry.Group) bool {
		bs, err := g.Encode()
		if err != nil {
			return true
		}
		groups[g.Base().GameMode] = append(groups[g.Base().GameMode], bs)
		return true
	})

	api.PM.Range(func(id string, p entry.Player) bool {
		bs, err := p.Encode()
		if err != nil {
			return true
		}
		players[p.Base().GameMode] = append(players[p.Base().GameMode], bs)
		return true
	})

	// create dir and save files
	dir := filepath.Join(os.TempDir(), "matcher")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	bs, err := json.Marshal(rooms)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "rooms.json"), bs, 0600); err != nil {
		return err
	}

	bs, err = json.Marshal(teams)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "teams.json"), bs, 0600); err != nil {
		return err
	}

	bs, err = json.Marshal(groups)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "groups.json"), bs, 0600); err != nil {
		return err
	}

	bs, err = json.Marshal(players)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "players.json"), bs, 0600); err != nil {
		return err
	}
	return nil
}

// ReloadEntries reloads the entries when the server starts.
//
//nolint:gocyclo
func (api *API) ReloadEntries() error {
	// read files
	dir := filepath.Join(os.TempDir(), "matcher")

	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		bs, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return err
		}

		data := make(map[constant.GameMode][][]byte)
		if err := json.Unmarshal(bs, &data); err != nil {
			return err
		}

		switch file.Name() {
		case "rooms.json":
			for _, rs := range data {
				for _, rbs := range rs {
					var r = new(entry.Room)
					if err := (*r).Decode(rbs); err != nil {
						return err
					}
					api.RM.Add((*r).ID(), *r)
				}
			}
		case "teams.json":
			for _, ts := range data {
				for _, tbs := range ts {
					var t = new(entry.Team)
					if err := (*t).Decode(tbs); err != nil {
						return err
					}
					api.TM.Add((*t).ID(), *t)
				}
			}
		case "groups.json":
			for _, gs := range data {
				for _, gbs := range gs {
					var g = new(entry.Group)
					if err := (*g).Decode(gbs); err != nil {
						return err
					}
					api.GM.Add((*g).ID(), *g)
				}
			}
		case "players.json":
			for _, ps := range data {
				for _, pbs := range ps {
					var p = new(entry.Player)
					if err := (*p).Decode(pbs); err != nil {
						return err
					}
					api.PM.Add((*p).Base().UID(), *p)
				}
			}
		}
	}

	return nil
}
