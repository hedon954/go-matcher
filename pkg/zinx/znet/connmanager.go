package znet

import (
	"fmt"
	"sync"

	"github.com/hedon954/go-matcher/pkg/zinx/ziface"
)

type ConnManager struct {
	connLock    sync.RWMutex
	connections map[uint64]ziface.IConnection
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint64]ziface.IConnection),
	}
}

func (cm *ConnManager) Add(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	cm.connections[conn.GetConnID()] = conn

	fmt.Println("connection added to connection manager successfully: connID = ", conn.GetConnID())
}

func (cm *ConnManager) Remove(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	delete(cm.connections, conn.GetConnID())

	fmt.Println("connection removed from connection manager successfully: connID = ", conn.GetConnID())
}

func (cm *ConnManager) Get(connID uint64) (ziface.IConnection, error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	}
	return nil, fmt.Errorf("connection not found: connID = %d", connID)
}

func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

func (cm *ConnManager) ClearConn() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	num := len(cm.connections)
	for _, conn := range cm.connections {
		conn.Stop()
	}
	cm.connections = make(map[uint64]ziface.IConnection)

	fmt.Println("clear all connections successfully, num = ", num)
}
