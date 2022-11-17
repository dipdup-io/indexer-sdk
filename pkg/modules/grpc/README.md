# gRPC module

The package contains basic implementation of gRPC server and client.

## Usage

Inherit `Server` module to implement custom gRPC server module. It realizes default `Module` interface. For example:

```go
type Server struct {
	*grpc.Server
}
```


Inherit `Client` module to implement custom gRPC client module. It realizes default `Module` interface. For example:

```go
type Client struct {
	*grpc.Client
}
```