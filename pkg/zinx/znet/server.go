package znet

import (
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"github.com/hedon954/go-matcher/pkg/safe"
	"github.com/hedon954/go-matcher/pkg/zinx/utils"
	"github.com/hedon954/go-matcher/pkg/zinx/ziface"
)

type Server struct {
	Name        string
	IPVersion   string
	IP          string
	Port        int
	ConnIDGen   atomic.Uint64 // TODO: use distributed unique id to support multiple servers
	msgHandler  ziface.IMsgHandle
	ConnMgr     ziface.IConnManager
	onConnStart func(conn ziface.IConnection)
	onConnStop  func(conn ziface.IConnection)
}

func NewServer(conf string) ziface.IServer {
	utils.GlobalObject.Reload(conf)

	s := &Server{
		IPVersion:  "tcp4",
		Name:       utils.GlobalObject.Name,
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TCPPort,
		ConnIDGen:  atomic.Uint64{},
		msgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

func (s *Server) Start() {
	fmt.Printf("[START] Server listener at IP: %s, Port: %d, is starting\n", s.IP, s.Port)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d, WorkerPoolSize: %d, MaxWorkerTaskLen: %d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn,
		utils.GlobalObject.WorkPoolSize, utils.GlobalObject.MaxWorkerTaskLen)

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

	safe.Go(func() {
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept error: ", err)
				continue
			}

			if utils.GlobalObject.MaxConn > 0 && s.ConnMgr.Len() > utils.GlobalObject.MaxConn {
				fmt.Println("too many connections, MaxConn = ", utils.GlobalObject.MaxConn)
				_ = conn.Close()
				continue
			}

			dealConn := NewConnection(s, conn, s.ConnIDGen.Add(1), s.msgHandler)
			safe.Go(dealConn.Start)
		}
	})
}

func (s *Server) Stop() {
	fmt.Printf("[STOP] Zinx server, name: %s\n", s.Name)

	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	s.Start()
	select {}
}

func (s *Server) AddRouter(msgID uint32, route ziface.IRouter) {
	s.msgHandler.AddRouter(msgID, route)

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
