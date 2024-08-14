package ztest

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/hedon954/go-matcher/pkg/zinx/zconfig"
	"github.com/hedon954/go-matcher/pkg/zinx/ziface"
	"github.com/hedon954/go-matcher/pkg/zinx/znet"
)

func TestServer(t *testing.T) {
	conf := zconfig.Load("zinx.yml")
	s := znet.NewServer(conf)

	s.SetOnConnStart(DoConnectionStart)
	s.SetOnConnStop(DoConnectionStop)

	s.AddRouter(0, func(request ziface.IRequest) {
		fmt.Println("[Call Router Handle]")
		fmt.Printf("receive from client, msgID=%d, data=%s\n", request.GetMsgID(), string(request.GetData()))

		err := request.GetConnection().SendMsg(0, []byte("ping...ping...ping..."))
		if err != nil {
			fmt.Println("call back in ping error: ", err)
		}
	})
	s.AddRouter(1, func(request ziface.IRequest) {
		fmt.Println("Call HelloZinxRouter Handle")
		fmt.Println("receive from client, msgID=", request.GetMsgID(), ", data=", string(request.GetData()))

		err := request.GetConnection().SendMsg(1, []byte("Hello Zinx Router V0.6!"))
		if err != nil {
			fmt.Println("call back in hello error: ", err)
		}
	})

	go s.Serve()
	time.Sleep(10 * time.Millisecond)
	go StartClient()
	time.Sleep(20 * time.Millisecond)
	s.Stop()
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

func StartClient() {
	fmt.Println("Client Test ...start")
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	i := 0

	for i <= 10 {
		dp := znet.NewDataPack(zconfig.DefaultConfig)

		i++

		// send to server
		msg, _ := dp.Pack(znet.NewMsgPackage(uint32(i%2), []byte("Hello Zinx!")))
		_, err = conn.Write(msg)
		if err != nil {
			fmt.Println("write to server error", err)
			return
		}

		// read from server
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			fmt.Println("read head error", err)
			break
		}

		// read msg head
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			// read msg body
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetDataLen())
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("read msg data error", err)
				return
			}
			fmt.Println("-----> Receive Server Msg: ID=", msg.GetMsgID(), ", len=", msg.GetDataLen(), ", data=", string(msg.GetData()))
		}

		time.Sleep(1 * time.Millisecond)
	}

	_ = conn.Close()
}
