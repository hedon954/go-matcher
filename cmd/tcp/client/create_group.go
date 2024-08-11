package main

import (
	"fmt"
	"net"

	"github.com/hedon954/go-matcher/internal/pb"
	"github.com/hedon954/go-matcher/pkg/zinx/ziface"
	"github.com/hedon954/go-matcher/pkg/zinx/znet"

	"google.golang.org/protobuf/proto"
)

func sendCreateGroup(conn net.Conn, dp ziface.IDataPack) {
	var req = &pb.CreateGroupReq{
		PlayerInfo: &pb.PlayerInfo{
			Uid:         "uid",
			GameMode:    pb.GameMode_GAME_MODE_GOAT_GAME,
			ModeVersion: 1,
			Star:        0,
			Rank:        0,
			Glicko2Info: &pb.Glicko2Info{
				Mmr:  0.0,
				Star: 0,
				Rank: 0,
			},
		},
	}
	bs, _ := proto.Marshal(req)
	msg, _ := dp.Pack(znet.NewMsgPackage(uint32(pb.ReqType_REQ_TYPE_CREATE_GROUP), bs))
	_, err := conn.Write(msg)
	if err != nil {
		fmt.Println("write to server error", err)
		return
	}
}
