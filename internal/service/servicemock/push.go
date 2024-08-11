package servicemock

import (
	"context"

	"github.com/hedon954/go-matcher/internal/entry"
	"github.com/hedon954/go-matcher/internal/pto"
)

type PushMock struct{}

func (p *PushMock) PushPlayerOnlineState(context.Context, []string, entry.PlayerOnlineState) {}
func (p *PushMock) PushGroupInfo(context.Context, []string, *pto.GroupInfo)                  {}
func (p *PushMock) PushInviteMsg(context.Context, *pto.InviteMsg)                            {}
func (p *PushMock) PushAcceptInvite(context.Context, string, string)                         {}
func (p *PushMock) PushRefuseInvite(context.Context, string, string, string)                 {}
func (p *PushMock) PushGroupDissolve(context.Context, []string, int64)                       {}
func (p *PushMock) PushGroupState(context.Context, []string, int64, entry.GroupState)        {}
func (p *PushMock) PushVoiceState(context.Context, []string, *pto.UserVoiceState)            {}
func (p *PushMock) PushKick(context.Context, string, int64)                                  {}
func (p *PushMock) PushMatchInfo(context.Context, []string, *pto.MatchInfo)                  {}
func (p *PushMock) PushCancelMatch(context.Context, []string, string)                        {}
func (p *PushMock) PushReady(context.Context, []string, string)                              {}
func (p *PushMock) PushUnReady(context.Context, []string, string)                            {}
