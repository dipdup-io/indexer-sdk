-include .env
export $(shell sed 's/=.*//' .env)

lint:
	golangci-lint run

build-proto:
	protoc -I=. -I=${GOPATH}/src --go-grpc_out=${GOPATH}/src --go_out=${GOPATH}/src ${GOPATH}/src/github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/proto/*.proto

build-example-proto:
	protoc -I=. -I=${GOPATH}/src --go-grpc_out=${GOPATH}/src --go_out=${GOPATH}/src ${GOPATH}/src/github.com/dipdup-net/indexer-sdk/examples/grpc/proto/*.proto

example-grpc:
	cd examples/grpc && go run .