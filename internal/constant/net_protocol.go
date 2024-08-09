package constant

type NetProtocol string

const (
	TCP   NetProtocol = "tcp"
	UDP   NetProtocol = "udp"
	WS    NetProtocol = "ws"
	WSS   NetProtocol = "wss"
	KCP   NetProtocol = "kcp"
	GRPC  NetProtocol = "grpc"
	GRPCS NetProtocol = "grpcs"
)
