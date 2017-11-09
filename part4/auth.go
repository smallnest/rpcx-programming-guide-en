# Auth

For the security, only authorized clients can invoke services.

Clients must set an authorization token which can gets from OAuth/OAuth2 or granted access token.

Servers receive requests, first check  the auth token. They reject requests with invalid tokens.

```go
func main() {
	flag.Parse()

	s := server.NewServer()
	s.RegisterName("Arith", new(example.Arith), "")
	s.AuthFunc = auth
	s.Serve("reuseport", *addr)
}

func auth(ctx context.Context, req *protocol.Message, token string) error {

	if token == "bearer tGzv3JOkF0XG5Qx2TlKWIA" {
		return nil
	}

	return errors.New("invalid token")
}
```

Server must define a AuthFunc to validate auth token. In the above example, only requests with token `bearer tGzv3JOkF0XG5Qx2TlKWIA` are valid.


Client must set the token:

```go

func main() {
	flag.Parse()

	d := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")

	option := client.DefaultOption
	option.ReadTimeout = 10 * time.Second

	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, option)
	defer xclient.Close()

	//xclient.Auth("bearer tGzv3JOkF0XG5Qx2TlKWIA")
	xclient.Auth("bearer abcdefg1234567")

	args := &example.Args{
		A: 10,
		B: 20,
	}

	reply := &example.Reply{}
	ctx := context.WithValue(context.Background(), share.ReqMetaDataKey, make(map[string]string))
	err := xclient.Call(ctx, "Mul", args, reply)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}

	log.Printf("%d * %d = %d", args.A, args.B, reply.C)

}
```