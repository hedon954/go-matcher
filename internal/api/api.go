package api

import (
	"fmt"

	"github.com/hedon954/go-matcher/internal/config"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/matcher"
	"github.com/hedon954/go-matcher/internal/matcher/glicko2"
	"github.com/hedon954/go-matcher/internal/repository"
	"github.com/hedon954/go-matcher/internal/service"
	"github.com/hedon954/go-matcher/internal/service/matchimpl"
	"github.com/hedon954/go-matcher/pkg/timer"

	timerasynq "github.com/hedon954/go-matcher/pkg/timer/asynq"
	timernative "github.com/hedon954/go-matcher/pkg/timer/native"

	"github.com/hibiken/asynq"
)

type API struct {
	MS service.Match
	M  *matcher.Matcher
	PM *repository.PlayerMgr
	GM *repository.GroupMgr
	TM *repository.TeamMgr
	RM *repository.RoomMgr
}

// Start initializes the api components and starts them.
func Start(
	sc config.Configer[config.ServerConfig],
	mc config.Configer[config.MatchConfig],
) (api *API, shutdown func()) {
	var (
		groupChannel = make(chan entry.Group, 1024)
		roomChannel  = make(chan entry.Room, 1024)
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
	groupChannel chan entry.Group, roomChannel chan entry.Room,
	dt timer.Operator[int64], gm *glicko2.Matcher, mgrs *repository.Mgrs) *API {

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

func NewEntryManagers() *repository.Mgrs {
	return &repository.Mgrs{
		PlayerMgr: repository.NewPlayerMgr(),
		GroupMgr:  repository.NewGroupMgr(0),
		TeamMgr:   repository.NewTeamMgr(0),
		RoomMgr:   repository.NewRoomMgr(0),
	}
}

func NewGlicko2Matcher(roomChannel chan entry.Room, conf *config.MatchConfig, mgrs *repository.Mgrs) *glicko2.Matcher {
	return glicko2.New(roomChannel, conf, conf.MatchInterval(), mgrs)
}
