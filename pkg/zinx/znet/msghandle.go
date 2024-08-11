package znet

import (
	"fmt"
	"strconv"

	"github.com/hedon954/go-matcher/pkg/safe"
	"github.com/hedon954/go-matcher/pkg/zinx/ziface"
	"github.com/hedon954/go-matcher/pkg/zinx/zutils"
)

type MsgHandle struct {
	Apis           map[uint32]ziface.IRouter
	WorkerPoolSize uint32
	TaskQueue      []chan ziface.IRequest
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: zutils.GlobalObject.WorkPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, zutils.GlobalObject.WorkPoolSize),
	}
}

func (mh *MsgHandle) DoMsgHandle(request ziface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), " is not FOUND")
		return
	}

	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeated api, msgID = " + strconv.Itoa(int(msgID)))
	}
	mh.Apis[msgID] = router
	fmt.Println("Add api msgID = ", msgID)
}

func (mh *MsgHandle) StarWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		mh.TaskQueue[i] = make(chan ziface.IRequest, zutils.GlobalObject.MaxWorkerTaskLen)
		safe.Go(func() { mh.startOneWorker(i, mh.TaskQueue[i]) })
	}
}

func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	workerID := request.GetConnection().GetConnID() % uint64(mh.WorkerPoolSize)

	fmt.Printf("Add ConnID = %d, request MsgID = %d to WorkerID = %d\n", request.GetConnection().GetConnID(), request.GetMsgID(), workerID)

	mh.TaskQueue[workerID] <- request
}

func (mh *MsgHandle) startOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("[WORKER] worker ID = ", workerID, " is started")
	for request := range taskQueue {
		mh.DoMsgHandle(request)
	}
}
