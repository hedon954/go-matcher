package ziface

type IMessage interface {
	GetDataLen() uint32
	GetMsgID() uint32
	GetData() []byte
	SetData([]byte)
}
