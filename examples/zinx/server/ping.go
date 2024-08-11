package main

import (
	"fmt"

	"github.com/hedon954/go-matcher/pkg/zinx/ziface"
)

func PingHandle(request ziface.IRequest) {
	fmt.Println("[Call Router Handle]")
	fmt.Printf("receive from client, msgID=%d, data=%s\n", request.GetMsgID(), string(request.GetData()))

	err := request.GetConnection().SendMsg(0, []byte("ping...ping...ping..."))
	if err != nil {
		fmt.Println("call back in ping error: ", err)
	}
}
