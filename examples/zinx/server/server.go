package main

import (
	"fmt"
	"time"

	"github.com/hedon954/go-matcher/pkg/zinx/zconfig"
	"github.com/hedon954/go-matcher/pkg/zinx/ziface"
	"github.com/hedon954/go-matcher/pkg/zinx/znet"
)

func main() {
	conf := zconfig.Load("conf/zinx.yml")
	s := znet.NewServer(conf)

	s.SetOnConnStart(DoConnectionStart)
	s.SetOnConnStop(DoConnectionStop)

	s.AddRouter(0, PingHandle)
	s.AddRouter(1, HelloHandle)

	s.Serve()
}

func DoConnectionStart(conn ziface.IConnection) {
	fmt.Println("DoConnectionStart is Called ... ConnID = ", conn.GetConnID())

	conn.SetProperty("name", "hedon")
	conn.SetProperty("now", time.Now().Unix())
}

func DoConnectionStop(conn ziface.IConnection) {
	fmt.Println("DoConnectionStop is Called ... ConnID = ", conn.GetConnID())

	fmt.Println(conn.GetProperty("name"))
	fmt.Println(conn.GetProperty("now"))
}
