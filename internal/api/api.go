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
	timermock "github.com/hedon954/go-matcher/pkg/timer/mock"

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
func Start(conf *config.Config) (api *API, shutdown func()) {
	var (
		groupChannel = make(chan entry.Group, 1024)
		roomChannel  = make(chan entry.Room, 1024)
		mgrs         = NewEntryManagers()
	)

	// init delay timer
	dt, err := NewDelayTime(conf)
	if err != nil {
		panic(err)
	}

	// init api
	api = NewAPI(conf, groupChannel, roomChannel, dt, NewGlicko2Matcher(roomChannel, conf, mgrs), mgrs)

	// start delay timer and match service
	go dt.Start()
	go api.M.Start()

	return api, func() {
		dt.Stop()
		api.M.Stop()
	}
}

func NewAPI(conf *config.Config,
	groupChannel chan entry.Group, roomChannel chan entry.Room,
	dt timer.Operator[int64], gm *glicko2.Matcher, mgrs *repository.Mgrs) *API {
	api := &API{
		PM: mgrs.PlayerMgr,
		GM: mgrs.GroupMgr,
		TM: mgrs.TeamMgr,
		RM: mgrs.RoomMgr,
		M:  matcher.New(groupChannel, gm),
		MS: matchimpl.NewDefault(conf.GroupPlayerLimit,
			mgrs, groupChannel, roomChannel, dt, conf),
	}
	return api
}

func NewDelayTime(conf *config.Config) (timer.Operator[int64], error) {
	switch conf.DelayTimerType {
	case config.DelayTimerTypeAsynq:
		return timerasynq.NewTimer[int64](&asynq.RedisClientOpt{
			Addr:     conf.AsynqRedis.Addr,
			Password: conf.AsynqRedis.Password,
			DB:       conf.AsynqRedis.DB,
		}), nil
	case config.DelayTimerTypeNative:
		return timermock.NewTimer(), nil
	default:
		return nil, fmt.Errorf("unsupported delay timer type: %s", conf.DelayTimerType)
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

func NewGlicko2Matcher(roomChannel chan entry.Room, conf *config.Config, mgrs *repository.Mgrs) *glicko2.Matcher {
	return glicko2.New(roomChannel, conf, conf.MatchInterval(), mgrs)
}
