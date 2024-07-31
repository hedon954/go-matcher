package config

type GroupConfig struct {
	PlayerLimit             int
	BroadcastUsersTimeoutMs int64
	QueueTimeoutMs          int64
	WaitAttrTimeoutMs       int64
	InviteTimeoutMs         int64
	GroupBroadcastTimeoutMs int64
}
