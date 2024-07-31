package def

const (
	RdsTimerDefaultInterval = 100
	RdsTimerGroupInvite     = "snake:rdstimer_group_invite"    // 队伍存在定时器，超时会被解散
	RdsTimerGroupQueue      = "snake:rdstimer_group_queue"     // 组队匹配定时器，超时会取消匹配
	RdsTimerGroupWaitAttr   = "snake:rdstimer_group_wait_attr" // 点匹配到加入队列，等待用户上传属性
	RdsTimerGroupBroadcast  = "snake:rdstimer_group_broadcast" // 用户加入后，推送给所有人
	RdsTimerRoomEnd         = "snake:rdstimer_room_end"
	RdsTimerRoomClear       = "snake:rdstimer_room_clear"
)
