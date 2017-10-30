# QuickStart

## Installation

First, you should install rpcx:

```go
go get -u -v github.com/smallnest/rpcx/...
```

It only install the basic features of rpcx. If you want to use etcd, you should add `etcd` tag:

```go
go get -u -v -tags "etcd" github.com/smallnest/rpcx/...
```

And if you want to use `quic`, you also need to add `quic` tag:

```go
go get -u -v -tags "quic etcd" github.com/smallnest/rpcx/...
```

I recommend you to install all tags even though you won't use them so far:

```go
go get -u -v -tags "reuseport quic kcp zookeeper etcd consul ping" github.com/smallnest/rpcx/...
```

**tags** means:

- **quic**: support quic transport
- **kcp**: support kcp transport
- **zookeeper**: support zookeeper register
- **etcd**: support etcd register
- **consul**: support consul register
- **ping**: support network quality load balancing
- **reuseport**: support reuseport

## Implement Service

You can write your service just like writing a plain Go struct:

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


`Arith` is a Go type and it has a one method `Mul`.
The first parameter of `Mul` method is `context.Context`。
The second parameter of `Mul` method is `args`。 It contains request data: `A` and `B`。
The third parameter of `Mul` method is `reply`。It is a pointer to `Reply`.
The `Mul` method returns an error (can be nil).
`Mul` sets `Reply.C`  with result of `A * B`.

Now you has defined one service `Arith` and implemented its `Mul` method. we will introduce how to explore it to server and how clients to call it in next steps.


## Implement Server

You can use three lines to explore the above service:

```go
    s := server.NewServer()
	s.RegisterName("Arith", new(Arith), "")
	s.Serve("tcp", ":8972")
```

You register your service with Name `Arith`。


You can register service by the below:
```go
s.Register(new(example.Arith), "")
```

It uses the service type name as service name.


## Implement Client

```go
    // #1
    d := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")
    // #2
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

    // #3
	args := &example.Args{
		A: 10,
		B: 20,
	}

    // #4
    reply := &example.Reply{}
    
    // #5
	err := xclient.Call(context.Background(), "Mul", args, reply)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}

	log.Printf("%d * %d = %d", args.A, args.B, reply.C)
```

`#1` defines a service discovery implementation. In this example we use the simplest discovery: `Peer2PeerDiscovery`. Client get the server address by accessing this discovery and connect the server directly.

`#2` creates a `XClient`, and passes FailMode, SelectMode and default option.
FailMode indicates the client how to handle call failures: retry, return fast or retry another server?
SelectMode indicates the client how to select a server if there are multiple servers for one service.

`#3` defines the request： we want to get result of `10 * 20`. Of course we know the result must be `200` but we will see whether reply from the server is `200`.

`#4` defines the result object, it is a zero value and actually rpcx can use it to know type of the result and then unmarshal the result to it.

`#5` call the remote service and gets the result synchronously。

## Call the service asynchronously

You can call the service synchronously:

```go
    d := client.NewPeer2PeerDiscovery("tcp@"+*addr2, "")
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	args := &example.Args{
		A: 10,
		B: 20,
	}

	reply := &example.Reply{}
	call, err := xclient.Go(context.Background(), "Mul", args, reply, nil)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}

	replyCall := <-call.Done
	if replyCall.Error != nil {
		log.Fatalf("failed to call: %v", replyCall.Error)
	} else {
		log.Printf("%d * %d = %d", args.A, args.B, reply.C)
	}
```

You must use `xclient.Go` to replace `xclient.Call` and ti returns a channel. You can read the result from the channel.