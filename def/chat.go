package def

type ChatType int

// 聊天类型
const (
	ChatTypeSingle    ChatType = 1
	ChatTypeChannel   ChatType = 2
	ChatTypeGroup     ChatType = 3
	ChatTypeBroadcast ChatType = 4
	ChatTypeClan      ChatType = 5
	ChatTypeRaceRoom  ChatType = 8
)
