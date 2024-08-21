package znet

import (
	"fmt"
	"strconv"

	"github.com/hedon954/go-matcher/pkg/zinx/zconfig"
	"github.com/hedon954/go-matcher/pkg/zinx/ziface"
)

type MsgHandle struct {
	Apis      map[uint32]ziface.HandleFunc
	TaskQueue []chan ziface.IRequest
	config    *zconfig.ZConfig
}

func NewMsgHandle(config *zconfig.ZConfig) *MsgHandle {
	return &MsgHandle{
		Apis:      make(map[uint32]ziface.HandleFunc),
		config:    config,
		TaskQueue: make([]chan ziface.IRequest, config.WorkPoolSize),
	}
}

func (mh *MsgHandle) DoMsgHandle(request ziface.IRequest) {
	handle, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), " is not FOUND")
		return
	}

	handle(request)
}

func (mh *MsgHandle) AddRouter(msgID uint32, handle ziface.HandleFunc) {
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeated api, msgID = " + strconv.Itoa(int(msgID)))
	}
	mh.Apis[msgID] = handle
}

func (mh *MsgHandle) StarWorkerPool() {
	for i := 0; i < int(mh.config.WorkPoolSize); i++ {
		mh.TaskQueue[i] = make(chan ziface.IRequest, mh.config.MaxWorkerTaskLen)
		go mh.startOneWorker(i, mh.TaskQueue[i])
	}
}

func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	workerID := request.GetConnection().GetConnID() % uint64(mh.config.WorkPoolSize)
	mh.TaskQueue[workerID] <- request
}

func (mh *MsgHandle) startOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	for request := range taskQueue {
		mh.DoMsgHandle(request)
	}
}
