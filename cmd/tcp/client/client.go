package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/hedon954/go-matcher/internal/pb"
	"github.com/hedon954/go-matcher/pkg/typeconv"
	"github.com/hedon954/go-matcher/pkg/zinx/znet"

	"google.golang.org/protobuf/proto"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for {
		dp := znet.NewDataPack()

		sendCreateGroup(conn, dp)

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
			fmt.Println("-----> Receive Server Msg: ID=", msg.GetMsgID(), ", len=", msg.GetDataLen(), ", data=", resolveData(msg.Data))
		}

		time.Sleep(1 * time.Second)
	}
}

func resolveData(data []byte) string {
	var rsp = new(pb.CommonRsp)
	err := proto.Unmarshal(data, rsp)
	if err != nil {
		fmt.Println("unmarshal error", err)
		return ""
	}

	if rsp.Code != pb.RspCode_RSP_CODE_SUCCESS {
		return rsp.Message
	}

	var bs []byte
	switch rsp.ReqType {
	case pb.ReqType_REQ_TYPE_CREATE_GROUP:
		bs, err = json.Marshal(typeconv.MustFromProto[pb.CreateGroupRsp](rsp.Data))
	default:
	}
	if err != nil {
		fmt.Println("json marshal error", err)
	}
	return string(bs)
}
