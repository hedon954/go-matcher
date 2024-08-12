package ziface

import (
	"context"
	"net"
)

type IConnection interface {
	Start(ctx context.Context)
	Stop()
	GetTCPConnection() *net.TCPConn
	GetConnID() uint64
	RemoteAddr() net.Addr
	SendMsg(id uint32, data []byte) error
	SetProperty(key string, value any)
	GetProperty(key string) (any, bool)
	RemoveProperty(key string)
}
