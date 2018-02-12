# Group

**Example:** [group](https://github.com/rpcx-ecosystem/rpcx-examples3/tree/master/group)

When you register a service, maybe you have noticed there is third parameter and we set it as empty string at most examples. Actually you add some meta for those services.

You can manage rpcx service via [rpcx-ui](https://github.com/smallnest/rpcx-ui)).


`group` is a meta data. If you set `group` metadata for some services, only clients in this `group` can access those services.

```go server.go

func main() {
	flag.Parse()

	go createServer1(*addr1, "")
	go createServer2(*addr2, "group=test")

	select {}
}

func createServer1(addr, meta string) {
	s := server.NewServer()
	s.RegisterName("Arith", new(example.Arith), meta)
	s.Serve("tcp", addr)
}

func createServer2(addr, meta string) {
	s := server.NewServer()
	s.RegisterName("Arith", new(Arith), meta)
	s.Serve("tcp", addr)
}
```

Discovery can find the group. Client can use `option.Group` to set group.

If you have not set `option.Group`, clients can access any services whether services set group or not.

```go client.go
	option := client.DefaultOption
	option.Group = "test"
	xclient := client.NewXClient("Arith", client.Failover, client.RoundRobin, d, option)
	defer xclient.Close()

	args := &example.Args{
		A: 10,
		B: 20,
	}

	for {
		reply := &example.Reply{}
		err := xclient.Call(context.Background(), "Mul", args, reply)
		if err != nil {
			log.Fatalf("failed to call: %v", err)
		}

		log.Printf("%d * %d = %d", args.A, args.B, reply.C)
		time.Sleep(1e9)
	}
```