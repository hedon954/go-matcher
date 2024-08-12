package znet

import (
	"context"
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"github.com/hedon954/go-matcher/pkg/zinx/ziface"
	"github.com/hedon954/go-matcher/pkg/zinx/zutils"
)

type Server struct {
	Name            string
	IPVersion       string
	IP              string
	Port            int
	ConnIDGen       atomic.Uint64 // TODO: use distributed unique id to support multiple servers
	msgHandler      ziface.IMsgHandle
	ConnMgr         ziface.IConnManager
	notifyConnClose chan ziface.IConnection
	onConnStart     func(conn ziface.IConnection)
	onConnStop      func(conn ziface.IConnection)
	cancelCtx       context.Context
	cancelFunc      func()
}

func NewServer(conf string) ziface.IServer {
	if conf != "" {
		zutils.GlobalObject.Reload(conf)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	s := &Server{
		IPVersion:       "tcp4",
		Name:            zutils.GlobalObject.Name,
		IP:              zutils.GlobalObject.Host,
		Port:            zutils.GlobalObject.TCPPort,
		ConnIDGen:       atomic.Uint64{},
		notifyConnClose: make(chan ziface.IConnection, 1),
		msgHandler:      NewMsgHandle(),
		ConnMgr:         NewConnManager(),
		cancelCtx:       ctx,
		cancelFunc:      cancelFunc,
	}
	return s
}

func (s *Server) Start() {
	fmt.Printf("[START] Server listener at IP: %s, Port: %d, is starting\n", s.IP, s.Port)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d, WorkerPoolSize: %d, MaxWorkerTaskLen: %d\n",
		zutils.GlobalObject.Version, zutils.GlobalObject.MaxConn,
		zutils.GlobalObject.WorkPoolSize, zutils.GlobalObject.MaxWorkerTaskLen)

	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("resolve tcp addr error: ", err)
		return
	}

	listener, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		fmt.Println("listen", s.IPVersion, "error: ", err)
		return
	}

	s.msgHandler.StarWorkerPool()
	time.Sleep(time.Millisecond)
	fmt.Printf("start Zinx server: %s successfully, now listening\n", s.Name)

	go func() {
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept error: ", err)
				continue
			}

			if zutils.GlobalObject.MaxConn > 0 && s.ConnMgr.Len() > zutils.GlobalObject.MaxConn {
				fmt.Println("too many connections, MaxConn = ", zutils.GlobalObject.MaxConn)
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
	fmt.Printf("[STOP] Zinx server, name: %s\n", s.Name)
	s.ConnMgr.ClearConn()
	s.cancelFunc()
}

func (s *Server) Serve() {
	s.Start()
	select {} // TODO: graceful shutdown
}

func (s *Server) AddRouter(msgID uint32, handle ziface.HandleFunc) {
	s.msgHandler.AddRouter(msgID, handle)

	fmt.Printf("Add Router[%d] successfully!\n", msgID)
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
		fmt.Println("---> CallOnConnStart()")
		s.onConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.onConnStop != nil {
		fmt.Println("---> CallOnConnStop()")
		s.onConnStop(conn)
	}
}

func (s *Server) NotifyClose(c ziface.IConnection) {
	s.notifyConnClose <- c
}
