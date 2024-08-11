package znet

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/hedon954/go-matcher/pkg/zinx/ziface"
	"github.com/hedon954/go-matcher/pkg/zinx/zutils"
)

const (
	idLen   = 4
	dataLen = 4
)

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	return dataLen + idLen
}

func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, fmt.Errorf("write data len occurs error %w", err)
	}

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, fmt.Errorf("write msg id occurs error %w", err)
	}

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, fmt.Errorf("write msg data occurs error %w", err)
	}

	return dataBuff.Bytes(), nil
}

func (dp *DataPack) Unpack(data []byte) (ziface.IMessage, error) {
	dataBuff := bytes.NewBuffer(data)
	msg := &Message{}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, fmt.Errorf("read data len occurs error %w", err)
	}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.ID); err != nil {
		return nil, fmt.Errorf("read msg id occurs error %w", err)
	}

	if zutils.GlobalObject.MaxPacketSize > 0 && msg.DataLen > zutils.GlobalObject.MaxPacketSize {
		return nil, fmt.Errorf("too large msg data len %d, limit %d", msg.DataLen, zutils.GlobalObject.MaxPacketSize)
	}

	// ...here we do not read msg data

	return msg, nil
}
