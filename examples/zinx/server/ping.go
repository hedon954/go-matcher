package main

import (
	"fmt"

	"github.com/hedon954/go-matcher/pkg/zinx/ziface"
	"github.com/hedon954/go-matcher/pkg/zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (p *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("[Call Router PreHandle]")

	err := request.GetConnection().SendMsg(0, []byte("before ping..."))
	if err != nil {
		fmt.Println("call back before ping error: ", err)
	}
}

func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("[Call Router Handle]")
	fmt.Printf("receive from client, msgID=%d, data=%s\n", request.GetMsgID(), string(request.GetData()))

	err := request.GetConnection().SendMsg(0, []byte("ping...ping...ping..."))
	if err != nil {
		fmt.Println("call back in ping error: ", err)
	}
}

func (p *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("[Call Router PostHandle]")

	err := request.GetConnection().SendMsg(0, []byte("after ping..."))
	if err != nil {
		fmt.Println("call back after ping error: ", err)
	}
}
