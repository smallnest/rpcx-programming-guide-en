# Fork

**Example:** [broadcast](https://github.com/rpcx-ecosystem/rpcx-examples3/tree/master/broadcast)


`Broadcast` is a method of `XClient` and you can use it to send a request to all servers that contains this service.

If all servers return response whithout an error, `Broadcast` will return OK for this XClient. If any of servers returns an error, `Broadcast` returns an error of those errors.


```go

func main() {
    ……

	xclient := client.NewXClient("Arith", client.Failover, client.RoundRobin, d, client.DefaultOption)
	defer xclient.Close()

	args := &example.Args{
		A: 10,
		B: 20,
	}

	for {
		reply := &example.Reply{}
		err := xclient.Broadcast(context.Background(), "Mul", args, reply)
		if err != nil {
			log.Fatalf("failed to call: %v", err)
		}

		log.Printf("%d * %d = %d", args.A, args.B, reply.C)
		time.Sleep(1e9)
	}

}
```