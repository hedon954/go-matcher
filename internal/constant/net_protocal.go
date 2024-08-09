package constant

type NetProtocal string

const (
	TCP   NetProtocal = "tcp"
	UDP   NetProtocal = "udp"
	WS    NetProtocal = "ws"
	WSS   NetProtocal = "wss"
	KCP   NetProtocal = "kcp"
	GRPC  NetProtocal = "grpc"
	GRPCS NetProtocal = "grpcs"
)
