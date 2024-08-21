package znet

import (
	"context"
	"fmt"
	"net"
	"sync/atomic"

	"github.com/hedon954/go-matcher/pkg/zinx/zconfig"
	"github.com/hedon954/go-matcher/pkg/zinx/ziface"
)

type Server struct {
	config          *zconfig.ZConfig
	ConnIDGen       atomic.Uint64 // TODO: use distributed unique id to support multiple servers
	msgHandler      ziface.IMsgHandle
	ConnMgr         ziface.IConnManager
	notifyConnClose chan ziface.IConnection
	onConnStart     func(conn ziface.IConnection)
	onConnStop      func(conn ziface.IConnection)
	cancelCtx       context.Context
	cancelFunc      func()
}

func NewServer(conf *zconfig.ZConfig) ziface.IServer {
	ctx, cancelFunc := context.WithCancel(context.Background())
	s := &Server{
		config:          conf,
		ConnIDGen:       atomic.Uint64{},
		notifyConnClose: make(chan ziface.IConnection, 1),
		msgHandler:      NewMsgHandle(conf),
		ConnMgr:         NewConnManager(),
		cancelCtx:       ctx,
		cancelFunc:      cancelFunc,
	}
	return s
}

func (s *Server) Start() {
	fmt.Printf("[START] Server listener at IP: %s, Port: %d, is starting\n", s.config.Host, s.config.TCPPort)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d, WorkerPoolSize: %d, MaxWorkerTaskLen: %d\n",
		s.config.Version, s.config.MaxConn,
		s.config.WorkPoolSize, s.config.MaxWorkerTaskLen)

	addr, err := net.ResolveTCPAddr(s.config.IPVersion, fmt.Sprintf("%s:%d", s.config.Host, s.config.TCPPort))
	if err != nil {
		fmt.Println("resolve tcp addr error: ", err)
		return
	}

	listener, err := net.ListenTCP(s.config.IPVersion, addr)
	if err != nil {
		fmt.Println("listen", s.config.IPVersion, "error: ", err)
		return
	}

	s.msgHandler.StarWorkerPool()
	fmt.Printf("start Zinx server: %s successfully, now listening\n", s.config.Name)

	go func() {
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept error: ", err)
				continue
			}

			if s.config.MaxConn > 0 && s.ConnMgr.Len() > s.config.MaxConn {
				fmt.Println("too many connections, MaxConn = ", s.config.MaxConn)
				_ = conn.Close()
				continue
			}

			dealConn := NewConnection(s, conn, s.ConnIDGen.Add(1), s.msgHandler)
			go dealConn.Start(s.cancelCtx)
		}
	}()

	go func() {
		for conn := range s.notifyConnClose {
			s.ConnMgr.Remove(conn)
		}
	}()
}

func (s *Server) Stop() {
	fmt.Printf("[STOP] Zinx server, name: %s\n", s.config.Name)
	s.ConnMgr.ClearConn()
	s.cancelFunc()
}

func (s *Server) Serve() {
	s.Start()
	select {} // TODO: graceful shutdown
}

func (s *Server) AddRouter(msgID uint32, handle ziface.HandleFunc) {
	s.msgHandler.AddRouter(msgID, handle)
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

func (s *Server) SetOnConnStart(f func(conn ziface.IConnection)) {
	s.onConnStart = f
}

func (s *Server) SetOnConnStop(f func(conn ziface.IConnection)) {
	s.onConnStop = f
}

func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.onConnStart != nil {
		s.onConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.onConnStop != nil {
		s.onConnStop(conn)
	}
}

func (s *Server) NotifyClose(c ziface.IConnection) {
	s.notifyConnClose <- c
}

func (s *Server) Config() *zconfig.ZConfig {
	return s.config
}
