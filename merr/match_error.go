package merr

import (
	"errors"
)

var (
	ErrServerMaintain = errors.New("系统停服维护中")
	ErrGameNotOpen    = errors.New("玩法暂未开放")

	ErrGroupNotInQueue      = errors.New("队伍不在队列中")
	ErrInGroup              = errors.New("玩家正在队伍中")
	ErrInQueuing            = errors.New("玩家正在队列中")
	ErrInGaming             = errors.New("玩家正在游戏中")
	ErrInSettling           = errors.New("玩家正在结算中")
	ErrAlreadyMatched       = errors.New("已匹配成功")
	ErrNeedCreateGroupFirst = errors.New("请退出组队界面重新进入")
	ErrPlayerNotInGroup     = errors.New("用户不在队伍中")
	ErrPlayerNotExists      = errors.New("玩家不存在，请重新组队")
	ErrGroupFull            = errors.New("队伍已经满了")
	ErrGroupDissolved       = errors.New("队伍已解散")
	ErrGroupNotExist        = errors.New("队伍不存在")
	ErrGroupPlayerNotReady  = errors.New("队伍成员未准备")
	ErrGroupEmpty           = errors.New("队伍玩家为空")
	ErrInviteExpired        = errors.New("邀请已过期")
	ErrMessageInvalid       = errors.New("此内容包含违规信息，发送失败")
	ErrVersionNotSame       = errors.New("游戏版本或玩法不匹配")
	ErrCantKickSelf         = errors.New("不能踢自己")
	ErrPermissionDeny       = errors.New("权限不足")
	ErrOnlyOwnerCanMatch    = errors.New("只能由队长开始匹配")
	ErrInvalidMessage       = errors.New("此内容包含违规信息，发送失败")
	ErrUserBanned           = errors.New("你被禁赛了")
	ErrGameOffline          = errors.New("玩法已下线")
)

// 以下错误需重点关注，应该永远不会出现
var (
	ErrUnknownUserOnlineState   = errors.New("未知玩家状态")
	ErrUnknownGroupState        = errors.New("未知队伍状态")
	ErrUnsupportedMatchStrategy = errors.New("不支持的匹配策略")
)
