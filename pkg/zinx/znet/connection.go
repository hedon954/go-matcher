package znet

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"

	"github.com/hedon954/go-matcher/pkg/safe"
	"github.com/hedon954/go-matcher/pkg/zinx/ziface"
	"github.com/hedon954/go-matcher/pkg/zinx/zutils"
)

type Connection struct {
	TCPServer    ziface.IServer
	Conn         *net.TCPConn
	ConnID       uint64
	MsgHandler   ziface.IMsgHandle
	msgChan      chan []byte
	isClosed     atomic.Bool
	ExitBuffChan chan struct{}

	propertyLock sync.RWMutex
	properties   map[string]any
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint64, mh ziface.IMsgHandle) *Connection {
	c := &Connection{
		TCPServer:    server,
		Conn:         conn,
		ConnID:       connID,
		MsgHandler:   mh,
		ExitBuffChan: make(chan struct{}, 1),
		msgChan:      make(chan []byte, zutils.GlobalObject.MaxMsgChanLen),
		properties:   make(map[string]any),
	}

	c.TCPServer.GetConnMgr().Add(c)
	return c
}

func (c *Connection) Start(ctx context.Context) {
	go c.startReader(ctx)
	go c.startWriter()

	c.TCPServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	if !c.isClosed.CompareAndSwap(false, true) {
		return
	}

	// do hook
	c.TCPServer.CallOnConnStop(c)

	// close socket connection
	_ = c.Conn.Close()

	// notify tcp server to remove connection
	c.TCPServer.NotifyClose(c)

	c.ExitBuffChan <- struct{}{}
	// close all channels
	close(c.ExitBuffChan)
	close(c.msgChan) // TODO: check, if close msg chan then send msg in connection, would panic
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint64 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(id uint32, data []byte) error {
	if c.isClosed.Load() {
		return nil
	}

	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(id, data))
	if err != nil {
		fmt.Printf("pack msg %d, error %v\n", id, err)
		return err
	}

	c.msgChan <- msg
	return nil
}

func (c *Connection) SetProperty(key string, value any) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.properties[key] = value
}

func (c *Connection) GetProperty(key string) (any, bool) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	value, ok := c.properties[key]
	return value, ok
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.properties, key)
}

func (c *Connection) startReader(ctx context.Context) {
	fmt.Println("[Reader Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit")
	defer c.Stop()

	dp := NewDataPack()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// read header
			headData := make([]byte, dp.GetHeadLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
				fmt.Println("read msg head error ", err)
				return
			}

			// unpack to get data length and msg id
			msg, err := dp.Unpack(headData)
			if err != nil {
				fmt.Println("unpack error ", err)
				return
			}

			// read body according to data len
			data := make([]byte, msg.GetDataLen())
			if msg.GetDataLen() > 0 {
				if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
					fmt.Println("read msg data error ", err)
					return
				}
			}
			msg.SetData(data)

			// handle request
			req := Request{conn: c, msg: msg}
			if zutils.GlobalObject.WorkPoolSize > 0 {
				c.MsgHandler.SendMsgToTaskQueue(&req)
			} else {
				safe.Go(func() { c.MsgHandler.DoMsgHandle(&req) })
			}
		}
	}
}

func (c *Connection) startWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), " conn writer exit")
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("write data error: ", err)
				return
			}
		case <-c.ExitBuffChan:
			return
		}
	}
}
