# Timeout

**Example:** [timeout](https://github.com/rpcx-ecosystem/rpcx-examples3/tree/master/timeout)

The timeout pattern prevents remote procedure calls from waiting indefinitely for a response. A timeout specifies how long the RPC has to return. If there is no response by the timeout, the RPC invokes a fallback mechanism, whether it’s retrying the request, throwing an exception, or something else. This pattern is perhaps the most basic, fundamental resilience pattern used in RPCs.


## Server

You can use `OptionFn` to set servers' `readTimeout` and `writeTimeout`.

```go server struct
type Server struct {
	……
	readTimeout  time.Duration
	writeTimeout time.Duration
	……
}
```

`OptionFn` methods are:

```go
func WithReadTimeout(readTimeout time.Duration) OptionFn
func WithWriteTimeout(writeTimeout time.Duration) OptionFn 
```

## Client

There are two methods to set timeout for client.

One is set read/write deadline for connections and another is using `context.Context`.

### read/write deadline

Client's Option has three options to set timeout:

```go

type Option struct {
	……
	//ConnectTimeout sets timeout for dialing
	ConnectTimeout time.Duration
	// ReadTimeout sets readdeadline for underlying net.Conns
	ReadTimeout time.Duration
	// WriteTimeout sets writedeadline for underlying net.Conns
	WriteTimeout time.Duration
	……
}
```

`DefaultOption` sets the ConnectTimeout to 10 seconds but has not set ReadTimeout and WriteTimeout. If the timeout has not been set, there is no timeout for accessing.


### context.Context

`context.Context` since Go 1.7 can be used to control timeout.

You need to use `context.WithTimeout` to wrap a go context to invoke client.

```go
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
```