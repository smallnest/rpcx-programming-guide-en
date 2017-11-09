# Alias

This plugin can alias a service and its method.

For example, the below code use `a.b.c.d` as the alias of `Arith`, and `Times` as the alias of `Mul`.

```go
func main() {
	flag.Parse()

	a := serverplugin.NewAliasPlugin()
	a.Alias("a.b.c.d", "Times", "Arith", "Mul")
	s := server.NewServer()
	s.Plugins.Add(a)
	s.RegisterName("Arith", new(example.Arith), "")
	err := s.Serve("reuseport", *addr)
	if err != nil {
		panic(err)
	}
}
```


Client can use alias to call the service:

```go
func main() {
	flag.Parse()

	d := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")

	option := client.DefaultOption
	option.ReadTimeout = 10 * time.Second

	xclient := client.NewXClient("a.b.c.d", client.Failtry, client.RandomSelect, d, option)
	defer xclient.Close()

	args := &example.Args{
		A: 10,
		B: 20,
	}

	reply := &example.Reply{}
	err := xclient.Call(context.Background(), "Times", args, reply)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}

	log.Printf("%d * %d = %d", args.A, args.B, reply.C)

}
```