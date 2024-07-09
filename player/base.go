package player

import "sync"

// Base 是 Player 的基础类，所有游戏模式和所有匹配策略共用
type Base struct {
	// 在 Base 的内部方法中不要进行同步处理，统一交给外部方法调用
	sync.RWMutex
	uid     string
	GroupID int64
}

func (b *Base) Inner() *Base {
	return b
}

func (b *Base) UID() string {
	return b.uid
}
