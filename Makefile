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
		-o ${GOPATH}/src/github.com/dipdup-net/indexer-sdk/cmd/example \
		-c app \
		-a 0x002b1ee9B1CF77233F9f96Fc9ee6191D2b857Be2 \
		-p github.com/dipdup-net/indexer-sdk/cmd/example