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

generate:
	cd cmd/dipdup-gen && go run . abi ${GOPATH}/src/github.com/dipdup-net/indexer-sdk/cmd/dipdup-gen/ \
		-o ${GOPATH}/src/github.com/dipdup-io/uniswap \
		-c uniswap \
		-a 0xb78d5b29d50d0bb1ec8f9143bc64de0f4d1225df \
		-p github.com/dipdup-io/uniswap

test:
	go test ./...