# Transport

rpcx can commnunicate via TCP, HTTP Connect, UnixDomain, QUIC and KCP. You can also access rpcx by a http client via Gateway or HTTP_Invoke.

## TCP

It is the most used transport. High performance and easy. You can use TLS for TCP security.

**Example:** [101basic](https://github.com/rpcx-ecosystem/rpcx-examples3/tree/master/101basic)

Servers use `tcp` as network name and they registers services as`serviceName/tcp@ipaddress:port` in the registry.
```go server.go
s.Serve("tcp", *addr)
```

Client can access services by:
```go
d := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")
xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
defer xclient.Close()
```

## HTTP Connect

You can send `HTTP CONNECT` method to rpcx servers. Rpcx servers hijacks this connection and use it as a TCP connection.
Notice clients and servers don't use http requests/responses to communicate and they still use the binary protocol based on this connection.

The network name is `http` and the registered format is `serviceName/http@ipaddress:port`

It is not recommended and TCP is the first choice.

If you want to use http requests/responses to access, you should use gateway or http_invoke.

## Unixdomain

the network name is `unix`.

**Example:** [unix](https://github.com/rpcx-ecosystem/rpcx-examples3/tree/master/unixdomain)

## QUIC

From [wikipedia](https://en.wikipedia.org/wiki/QUIC)

> QUIC (Quick UDP Internet Connections, pronounced quick) is an experimental transport layer network protocol designed by Jim Roskind at Google, initially implemented in 2012, and announced publicly in 2013 as experimentation broadened. QUIC supports a set of multiplexed connections between two endpoints over User Datagram Protocol (UDP), and was designed to provide security protection equivalent to TLS/SSL, along with reduced connection and transport latency, and bandwidth estimation in each direction to avoid congestion. QUIC's main goal is to improve perceived performance of connection-oriented web applications that are currently using TCP. It also moves control of the congestion avoidance algorithms into the application space at both endpoints, rather than the kernel space, which it is claimed will allow these algorithms to improve more rapidly.

> In June 2015, an Internet Draft of a specification for QUIC was submitted to the IETF for standardization. A QUIC working group was established in 2016. The QUIC working group foresees multipath support and optional forward error correction (FEC) as the next step. The working group also focuses on network management issues that QUIC may introduce, aiming to produce an applicability and manageability statement in parallel to the actual protocol work. Internet statistics in 2017 suggest that QUIC now accounts for more than 5% of Internet traffic.

the network name is `quic`.

**Example:** [quic](https://github.com/rpcx-ecosystem/rpcx-examples3/tree/master/quic)

## KCP

[KCP](https://github.com/skywind3000/kcp) is a fast and reliable ARQ protocol.

the network is `kcp`.

**Example:** [kcp](https://github.com/rpcx-ecosystem/rpcx-examples3/tree/master/kcp)


## reuseport

the network name is `reuseport`.

**Example:** [reuseport](https://github.com/rpcx-ecosystem/rpcx-examples3/tree/master/reuseport)

It uses `tcp` protocol and set [SO_REUSEPORT](https://lwn.net/Articles/542629/) socket option for linux and unix servers.


## TLS

**Example:** [TLS](https://github.com/rpcx-ecosystem/rpcx-examples3/tree/master/tls)


You can set TLS in server:

```go
func main() {
	flag.Parse()

	cert, err := tls.LoadX509KeyPair("server.pem", "server.key")
	if err != nil {
		log.Print(err)
		return
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	s := server.NewServer(server.WithTLSConfig(config))
	s.RegisterName("Arith", new(example.Arith), "")
	s.Serve("tcp", *addr)
}
```

And set TLS in client:

```go
func main() {
	flag.Parse()

	d := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")

	option := client.DefaultOption

	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	option.TLSConfig = conf

	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, option)
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

}
```