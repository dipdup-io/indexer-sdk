# gRPC module

The package contains basic implementation of gRPC server and client.

## Usage

Inherit `Server` module to implement custom gRPC server module. It realizes default `Module` interface and handshake logic. For example:

```go
type Server struct {
	*grpc.Server
}
```


Inherit `AuthClient` module to implement custom gRPC client module. It realizes default `Module` interface and handshake logic. For example:

```go
type Client struct {
	*grpc.AuthClient
}
```