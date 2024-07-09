package merr

import (
	"errors"
)

var (
	ServerMaintain = errors.New("系统停服维护中")
	GameNotOpen    = errors.New("玩法暂未开放")

	GroupNotInQueue      = errors.New("队伍不在队列中")
	InGroup              = errors.New("玩家正在队伍中")
	InQueuing            = errors.New("玩家正在队列中")
	InGaming             = errors.New("玩家正在游戏中")
	InSettling           = errors.New("玩家正在结算中")
	AlreadyMatched       = errors.New("已匹配成功")
	NeedCreateGroupFirst = errors.New("请退出组队界面重新进入")
	PlayerNotInGroup     = errors.New("用户不在队伍中")
	PlayerNotExists      = errors.New("玩家不存在，请重新组队")
	GroupFull            = errors.New("队伍已经满了")
	GroupDissolved       = errors.New("队伍已解散")
	GroupNotExist        = errors.New("队伍不存在")
	GroupPlayerNotReady  = errors.New("队伍成员未准备")
	GroupEmpty           = errors.New("队伍玩家为空")
	InviteExpired        = errors.New("邀请已过期")
	MessageInvalid       = errors.New("此内容包含违规信息，发送失败")
	VersionNotSame       = errors.New("游戏版本或玩法不匹配")
	CantKickSelf         = errors.New("不能踢自己")
	PermissionDeny       = errors.New("权限不足")
	OnlyOwnerCanMatch    = errors.New("只能由队长开始匹配")
)

// 以下错误需重点关注，应该永远不会出现
var (
	UnknownUserOnlineState   = errors.New("未知玩家状态")
	UnknownGroupState        = errors.New("未知队伍状态")
	UnsupportedMatchStrategy = errors.New("不支持的匹配策略")
)
