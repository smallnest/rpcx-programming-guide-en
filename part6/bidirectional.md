# 双向通讯
**例子:** [bidirectional](https://github.com/rpcx-ecosystem/rpcx-examples3/tree/master/bidirectional)

绝大部分的rpc框架, client发送请求， server处理请求再返回给client, 这是标准的 `request-response` rpc 模型.

但是对于一些用户， 尤其是 `IOT` 的用户，有时候他们想要发送命令给客户端。当然可以通过客户端也部署一套服务，两端都有客户端和服务器，但是这样做太复杂和臃肿。

rpcx实现了一个简单的通知模型，可以让服务端直接发送请求给客户端。

它要求你在服务端保存客户端的连接， 可能你需要将连接和客户端标识绑定起来，然后发送通知给对应的连接。

## Server

服务端使用`SendMessage`发送消息给客户端， 数据是`[]byte`， 你还可以设置`servicePath`和 `serviceMethod`方便客户端区分数据。

你可以在客户端首次调用的时候保存它的连接，在服务的实现中使用`ctx.Value(server.RemoteConnContextKey)` 得到 `net.Conn`。

```go
func (s *Server) SendMessage(conn net.Conn, servicePath, serviceMethod string, metadata map[string]string, data []byte) error
```


```go server.go
func main() {
	flag.Parse()

	ln, _ := net.Listen("tcp", ":9981")
	go http.Serve(ln, nil)

	s := server.NewServer()
	//s.RegisterName("Arith", new(example.Arith), "")
	s.Register(new(Arith), "")
	go s.Serve("tcp", *addr)

	for !connected {
		time.Sleep(time.Second)
	}

	fmt.Printf("start to send messages to %s\n", clientConn.RemoteAddr().String())
	for {
		if clientConn != nil {
			err := s.SendMessage(clientConn, "test_service_path", "test_service_method", nil, []byte("abcde"))
			if err != nil {
				fmt.Printf("failed to send messsage to %s: %v\n", clientConn.RemoteAddr().String(), err)
				clientConn = nil
			}
		}
		time.Sleep(time.Second)
	}
}
```

## Client

客户端需要使用`NewBidirectionalXClient`创建`XClient`, 你需要传入一个处理channel， 然后你从这个channel中读取服务器端的请求就可以了。

```go client.go
func main() {
	flag.Parse()

	ch := make(chan *protocol.Message)

	d := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")
	xclient := client.NewBidirectionalXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption, ch)
	defer xclient.Close()

	args := &example.Args{
		A: 10,
		B: 20,
	}

	reply := &example.Reply{}
	err := xclient.Call(context.Background(), "Mul", args, reply)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}

	log.Printf("%d * %d = %d", args.A, args.B, reply.C)

	for msg := range ch {
		fmt.Printf("receive msg from server: %s\n", msg.Payload)
	}
}
```