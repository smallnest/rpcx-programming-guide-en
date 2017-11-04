# Server

You can implements services in Server side.

The type of services is not important. You can use customized type to keep states or use `struct{}` or `int` simplely.

You should start a TCP or UDP server to serve services.

And you can aslo set some plugins to add features to the server.


## Service

As service providers, you should define some service first.
Currently rpcx only support exported `methods` as functions of service.
And the exported method must satisfy the below rules:

- exported method of exported type
- three arguments, the first is of context.Context, both of exported type or three arguments
- the third argument is a pointer
- one return value, of type error

in rpcx 3.1 we want to add register raw functions that satisfy last three rules, but for rpcx 3.0, you must use `method` to implement services.

you can use `RegisterName` to register methods of `rcvr` and name of this service is `name`.
If you use `Register`, the generated service name is type name of `rcvr`.
You can add meta data of this service in registry, that can be used by clients or service manager, for example, `weight`、`geolocation`、`metrics`. 

```
func (s *Server) Register(rcvr interface{}, metadata string) error
func (s *Server) RegisterName(name string, rcvr interface{}, metadata string) error
```

There is an example that implements a `Mul` method:

```go
import "context"

type Args struct {
	A int
	B int
}

type Reply struct {
	C int
}

type Arith int

func (t *Arith) Mul(ctx context.Context, args *Args, reply *Reply) error {
	reply.C = args.A * args.B
	return nil
}
```

In this example, you can define `Arith` as `struct{}` and it doesn't affect the this service.

And you can define `args` as `Args`, this service aslo works.


## Server

After you define services, you should want to explore to users. You should start a TCP server or UDP server to listen.

The server supports below approaches to start, serve and close:

```go
    func NewServer(options ...OptionFn) *Server
    func (s *Server) Close() error
    func (s *Server) RegisterOnShutdown(f func())
    func (s *Server) Serve(network, address string) (err error)
    func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request)
```

First you use `NewServer` to create a server instance with options. Second you can invokde `Serve` or `ServeHTTP` to serve.

The server contains some fields (and unexported fields):
```go
type Server struct {
    Plugins PluginContainer
    // AuthFunc can be used to auth.
    AuthFunc func(ctx context.Context, req *protocol.Message, token string) error
    // contains filtered or unexported fields
}
```

`Plugins` caontains all plugins in this server. We will introduce it in later chapters.

`AuthFunc` is a authentication function which can check clients are authorized. It is also introduced later.


rpcx provides three `OptionFn` to set options:

```go

    func WithReadTimeout(readTimeout time.Duration) OptionFn
    func WithTLSConfig(cfg *tls.Config) OptionFn
    func WithWriteTimeout(writeTimeout time.Duration) OptionFn
```

It can set readTimeout, writeTimeout and tls.Config.

`ServeHTTP` uses a HTTP connection to serve services. It hijacks the connnection to communicate with clients.

`Serve` uses TCP or UDP ptotocol to communicate with clients.

rpcx supports the below `network`:

- tcp: recommended network
- http: hijack http connection
- unix: unix domain sockets
- reuseport: use `SO_REUSEPORT` socket option, only support Linux kernel 3.9+
- quic: support [quic protocol](https://en.wikipedia.org/wiki/QUIC)
- kcp: sopport [kcp protocol](https://github.com/skywind3000/kcp)


There is a server example:

```go
package main

import (
	"flag"

	example "github.com/rpcx-ecosystem/rpcx-examples3"
	"github.com/smallnest/rpcx/server"
)

var (
	addr = flag.String("addr", "localhost:8972", "server address")
)

func main() {
	flag.Parse()

	s := server.NewServer()
	//s.RegisterName("Arith", new(example.Arith), "")
	s.Register(new(example.Arith), "")
	s.Serve("tcp", *addr)
}
```