# gRPC Module

Basic gRPC server and client implementation with built-in subscription management and Prometheus metrics.

## Usage

```go
import "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc"
```

Usage is described by the [gRPC example](/examples/grpc/).

## Server

Embed `grpc.Server` and your generated unimplemented server to create a custom gRPC server module:

```go
type Server struct {
    *grpc.Server
    pb.UnimplementedTimeServiceServer

    subscriptions *grpc.Subscriptions[time.Time, *pb.Response]
    wg            *sync.WaitGroup
}
```

Register the generated server in `Start`:

```go
func (server *Server) Start(ctx context.Context) {
    pb.RegisterTimeServiceServer(server.Server.Server(), server)
    server.Server.Start(ctx)
    // start background work...
}
```

### Subscriptions

`Subscriptions[T, P]` manages stream subscriptions with generic types:
- `T` — source data type
- `P` — protobuf response type

Implement the `Subscription[T, P]` interface for custom filtering and conversion:

```go
func (server *Server) SubscribeOnTime(req *pb.Request, stream pb.TimeService_SubscribeOnTimeServer) error {
    return grpc.DefaultSubscribeOn[time.Time, *pb.Response](
        stream, server.subscriptions, NewTimeSubscription(),
    )
}

func (server *Server) UnsubscribeFromTime(ctx context.Context, req *generalPB.UnsubscribeRequest) (*generalPB.UnsubscribeResponse, error) {
    return grpc.DefaultUnsubscribe(ctx, server.subscriptions, req.Id)
}
```

### Server Configuration

```go
serverCfg := grpc.ServerConfig{
    Bind: "127.0.0.1:8889",
}
```

## Client

Embed `grpc.Client` and the generated client:

```go
type Client struct {
    *grpc.Client
    output *modules.Output
    client pb.TimeServiceClient
    wg     *sync.WaitGroup
}
```

Initialize the generated client in `Start` (after `Connect` opens the connection):

```go
func (client *Client) Start(ctx context.Context) {
    client.client = pb.NewTimeServiceClient(client.Connection())
}
```

### Subscribing

```go
func (client *Client) SubscribeOnTime(ctx context.Context) (uint64, error) {
    stream, err := client.client.SubscribeOnTime(ctx, new(pb.Request))
    if err != nil {
        return 0, err
    }
    return grpc.Subscribe[*pb.Response](stream, client.handleTime, client.wg)
}

func (client *Client) handleTime(ctx context.Context, data *pb.Response, id uint64) error {
    client.output.Push(data)
    return nil
}
```

## Protocol

Define your proto file importing the SDK's general proto:

```protobuf
syntax = "proto3";
import "github.com/dipdup-net/indexer-sdk/pkg/modules/grpc/proto/general.proto";

service TimeService {
    rpc SubscribeOnTime(Request) returns (stream Response);
    rpc UnsubscribeFromTime(UnsubscribeRequest) returns (UnsubscribeResponse);
}
```

Generate Go code:

```bash
protoc -I=. -I=${GOPATH}/src \
    --go-grpc_out=${GOPATH}/src \
    --go_out=${GOPATH}/src \
    your_proto.proto
```

## Inputs and Outputs

Server and client modules do not define default inputs/outputs — they must be implemented by the developer based on project-specific data structures and notification logic.

Full example: [`examples/grpc/`](/examples/grpc/)
