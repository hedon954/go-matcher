http:
	swag init
	go run ./cmd/http/main.go

genpb:
	protoc --go_out=. ./protos/*

docker-up:
	make build-ubuntu
	docker compose -f scripts/apm/docker-compose.yml up -d

docker-down:
	docker compose -f scripts/apm/docker-compose.yml down

docker-restart:
	make docker-down
	make docker-up

build-ubuntu:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o scripts/apm/build/go-matcher-http-api cmd/http/main.go
