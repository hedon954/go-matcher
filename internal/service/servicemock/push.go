package servicemock

import (
	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/pto"
)

type PushMock struct{}

func (p *PushMock) PushPlayerOnlineState([]string, entry.PlayerOnlineState) {}
func (p *PushMock) PushGroupPlayers([]string, *pto.GroupPlayers)            {}
func (p *PushMock) PushInviteMsg(*pto.InviteMsg)                            {}
func (p *PushMock) PushAcceptInvite(string, string)                         {}
func (p *PushMock) PushRefuseInvite(string, string, string)                 {}
func (p *PushMock) PushGroupDissolve([]string, int64)                       {}
func (p *PushMock) PushGroupState([]string, int64, entry.GroupState)        {}
func (p *PushMock) PushVoiceState([]string, *pto.UserVoiceState)            {}
func (p *PushMock) PushKick(string, int64)                                  {}
func (p *PushMock) PushMatchInfo([]string, *pto.MatchInfo)                  {}
func (p *PushMock) PushCancelMatch([]string, string)                        {}
func (p *PushMock) PushReady([]string, string)                              {}
func (p *PushMock) PushUnReady([]string, string)                            {}
