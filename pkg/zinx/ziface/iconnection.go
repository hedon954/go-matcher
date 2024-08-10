package ziface

import (
	"net"
)

type IConnection interface {
	Start()
	Stop()
	GetTCPConnection() *net.TCPConn
	GetConnID() uint64
	RemoteAddr() net.Addr
	SendMsg(id uint32, data []byte) error
	SetProperty(key string, value any)
	GetProperty(key string) (any, bool)
	RemoveProperty(key string)
}
