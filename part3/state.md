# State

**Example:** [state](https://github.com/rpcx-ecosystem/rpcx-examples3/tree/master/state)

`state` is another meta data. If you set `state=inactive` metadata for some services, only clients can not access those services even if those services are alive.

You can use it to disable a service temporally and don't shutdown the server.

You can manage rpcx service via [rpcx-ui](https://github.com/smallnest/rpcx-ui)).

```go server.go

func main() {
	flag.Parse()

	go createServer1(*addr1, "")
	go createServer2(*addr2, "state=inactive")

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


```go client.go
	xclient := client.NewXClient("Arith", client.Failover, client.RoundRobin, d, client.DefaultOption)
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
