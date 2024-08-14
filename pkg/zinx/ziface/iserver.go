package ziface

import (
	"github.com/hedon954/go-matcher/pkg/zinx/zconfig"
)

type IServer interface {
	Start()
	Stop()
	Serve()
	Config() *zconfig.ZConfig
	AddRouter(msgID uint32, handle HandleFunc)
	GetConnMgr() IConnManager
	SetOnConnStart(func(conn IConnection))
	SetOnConnStop(func(conn IConnection))
	CallOnConnStart(conn IConnection)
	CallOnConnStop(conn IConnection)
	NotifyClose(conn IConnection)
}

type HandleFunc = func(request IRequest)
