-include .env
export $(shell sed 's/=.*//' .env)

lint:
	golangci-lint run

build-proto:
	protoc -I=. --go-grpc_out=./pkg ./pkg/modules/grpc/proto/*.proto
	protoc -I=. --go_out=./pkg ./pkg/modules/grpc/proto/*.proto
