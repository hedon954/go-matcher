package ziface

type IMsgHandle interface {
	DoMsgHandle(request IRequest)
	AddRouter(msgID uint32, handle HandleFunc)
	StarWorkerPool()
	SendMsgToTaskQueue(request IRequest)
}
