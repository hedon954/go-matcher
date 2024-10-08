package main

import (
	"fmt"

	"github.com/hedon954/go-matcher/pkg/zinx/ziface"
)

func HelloHandle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle")
	fmt.Println("receive from client, msgID=", request.GetMsgID(), ", data=", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("Hello Zinx Router V0.6!"))
	if err != nil {
		fmt.Println("call back in hello error: ", err)
	}
}
