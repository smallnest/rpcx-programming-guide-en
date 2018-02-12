# Client

**Example:** [function](https://github.com/rpcx-ecosystem/rpcx-examples3/tree/master/function)

In the normal we register methods as services'methods. Methods must follow the below rules:

- exported method of exported type
- three arguments, the first is of context.Context, both of exported type for three arguments
- the third argument is a pointer
- one return value, of type error

And rpcx can also register raw functions as services and functions must follow the below rules:

- the function can be exported or not
- three arguments, the first is of context.Context, both of exported type for three arguments
- the third argument is a pointer
- one return value, of type error


There is an example.

Server must use `RegisterFunction` to register a function and provides a service name.

```go server.go
type Args struct {
	A int
	B int
}

type Reply struct {
	C int
}

func mul(ctx context.Context, args *Args, reply *Reply) error {
	reply.C = args.A * args.B
	return nil
}

func main() {
	flag.Parse()

	s := server.NewServer()
	s.RegisterFunction("a.fake.service", mul, "")
	s.Serve("tcp", *addr)
}
```

Client use the service name and function name to access:

```go client.go
	d := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")
	xclient := client.NewXClient("a.fake.service", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	args := &example.Args{
		A: 10,
		B: 20,
	}

	reply := &example.Reply{}
	err := xclient.Call(context.Background(), "mul", args, reply)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}

	log.Printf("%d * %d = %d", args.A, args.B, reply.C)
```