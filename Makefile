http:
	swag init
	go run -toolexec=/Users/wangjiahan/Downloads/skywalking ./cmd/http/main.go

genpb:
	protoc --go_out=. ./protos/*