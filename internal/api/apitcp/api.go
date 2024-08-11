package apitcp

import (
	"time"

	"github.com/hedon954/go-matcher/internal/config/mock"
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/matcher"
	"github.com/hedon954/go-matcher/internal/matcher/glicko2"
	"github.com/hedon954/go-matcher/internal/repository"
	"github.com/hedon954/go-matcher/internal/service"
	"github.com/hedon954/go-matcher/internal/service/matchimpl"
	"github.com/hedon954/go-matcher/pkg/zinx/ziface"

	timermock "github.com/hedon954/go-matcher/pkg/timer/mock"
)

type API struct {
	ms service.Match
	m  *matcher.Matcher
	pm *repository.PlayerMgr
	gm *repository.GroupMgr
	tm *repository.TeamMgr
	rm *repository.RoomMgr
}

func NewAPI(groupPlayerLimit int, matchInterval time.Duration) *API {
	var (
		groupChannel = make(chan entry.Group, 1024)
		roomChannel  = make(chan entry.Room, 1024)
	)

	var (
		playerMgr = repository.NewPlayerMgr()
		groupMgr  = repository.NewGroupMgr(0)
		teamMgr   = repository.NewTeamMgr(0)
		roomMgr   = repository.NewRoomMgr(0)
		configer  = &glicko2.Configer{
			Glicko2: new(mock.Glicko2Mock), // TODO: change
		}
		glicko2Matcher = glicko2.New(roomChannel, configer, matchInterval, playerMgr, groupMgr, teamMgr, roomMgr)
	)

	delayTimer := timermock.NewTimer() // TODO: get from param

	api := &API{
		pm: playerMgr,
		gm: groupMgr,
		tm: teamMgr,
		rm: roomMgr,
		m:  matcher.New(groupChannel, glicko2Matcher),
		ms: matchimpl.NewDefault(groupPlayerLimit, playerMgr, groupMgr, teamMgr, roomMgr, groupChannel, roomChannel, delayTimer),
	}

	go delayTimer.Start()
	go api.m.Start()
	return api
}

func (api *API) Bind(request ziface.IRequest) {

}

func (api *API) CreateGroup(request ziface.IRequest) {

}

func (api *API) EnterGroup(request ziface.IRequest) {

}

func (api *API) ExitGroup(request ziface.IRequest) {

}

func (api *API) DissolveGroup(request ziface.IRequest) {

}

func (api *API) KickPlayer(request ziface.IRequest) {

}

func (api *API) ChangeRole(request ziface.IRequest) {

}

func (api *API) Invite(request ziface.IRequest) {

}

func (api *API) AcceptInvite(request ziface.IRequest) {

}

func (api *API) RefuseInvite(request ziface.IRequest) {

}

func (api *API) SetNearbyJoinGroup(request ziface.IRequest) {

}

func (api *API) SetRecentJoinGroup(request ziface.IRequest) {

}

func (api *API) SetVoiceState(request ziface.IRequest) {

}

func (api *API) StartMatch(request ziface.IRequest) {

}

func (api *API) CancelMatch(request ziface.IRequest) {

}

func (api *API) Ready(request ziface.IRequest) {

}

func (api *API) Unready(request ziface.IRequest) {

}
