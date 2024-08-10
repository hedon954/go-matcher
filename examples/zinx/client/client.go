package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/hedon954/go-matcher/pkg/zinx/znet"
)

func main() {
	fmt.Println("Client Test ...start")
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	i := 0

	for {
		dp := znet.NewDataPack()

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

		time.Sleep(1 * time.Second)
	}
}
