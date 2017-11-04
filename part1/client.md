# Client

Clients use the same transport protocol to send requests to services and get responses from them.


```go
type Client struct {
    Conn net.Conn

    Plugins PluginContainer
    // contains filtered or unexported fields
}
```

`Conn` is the connection between client and server and `Plugins` contains client plugins.


It has methods:

```go
    func (client *Client) Call(ctx context.Context, servicePath, serviceMethod string, args interface{}, reply interface{}) error
    func (client *Client) Close() error
    func (c *Client) Connect(network, address string) error
    func (client *Client) Go(ctx context.Context, servicePath, serviceMethod string, args interface{}, reply interface{}, done chan *Call) *Call
    func (client *Client) IsClosing() bool
    func (client *Client) IsShutdown() bool
```

`Call` is a synchronous service invocation。 Clients are blocked until response are received or some errors are returned. While `Go` is an asynchronous invocation. It returns a Call pointer and you can check the `*Call` to get the result or the error.

`Close` closes the connnection to the service. It won't wait unfinished requests and closes the connection inmmediately.

`IsClosing` indicates the client is closing and won't accept new invocations.
`IsShutdown` indicates the client won't receive responses from the service.


`Client` uses the default [CircuitBreaker (circuit.NewRateBreaker(0.95, 100))](https://godoc.org/github.com/rubyist/circuitbreaker#NewRateBreaker) to handle errors. This is a poplular rpc error handling style. When the error rate hits the threshold, this service is marked unavailable in 10 second window. You can implement your customzied CircuitBreaker.

There is a client example:


```go
	client := &Client{
		option: DefaultOption,
	}

	err := client.Connect("tcp", addr)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer client.Close()

	args := &Args{
		A: 10,
		B: 20,
	}

	reply := &Reply{}
	err = client.Call(context.Background(), "Arith", "Mul", args, reply)
	if err != nil {
		t.Fatalf("failed to call: %v", err)
	}

	if reply.C != 200 {
		t.Fatalf("expect 200 but got %d", reply.C)
	}
```


## XClient

`XClient` is a client wrapper and adds some service discovery and service governance features.

```go
type XClient interface {
    SetPlugins(plugins PluginContainer)
    ConfigGeoSelector(latitude, longitude float64)
    Auth(auth string)

    Go(ctx context.Context, serviceMethod string, args interface{}, reply interface{}, done chan *Call) (*Call, error)
    Call(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error
    Broadcast(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error
    Fork(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error
    Close() error
}
```

`SetPlugins` set Plugin container and `Auth` set authentication token.

`ConfigGeoSelector` is special method to set latitude and longitude of this client and use geography selector.

One XClient is responsable for one service and can call all methods of this service via `serviceMethod` argument. If you want to call multiple services, you must create one xclient for each service.

Only one shared xclient is needed for one corresponding service in a application. It can be shared by goroutines and it is goroutine-safe.

`Go` is for async invocations and `Call` is fo sync invocations.

XClient use the single connection to one service node and it caches the connection until the connection is broken or closed.

### service discovery

rpcx support a lot of service discovery and you can aslo implement your service discovery.

- * [Peer to Peer](part2/registry.md#peer2peer): the client connects the single service directly. It acts like the `client` type.
  * [Peer to Multiple](part2/registry.md#multiple): the client can connnect multiple services. Ther services are configured programmatically.
  * [Zookeeper](part2/registry.md#zookeeper): find the services via zookeeper.
  * [Etcd](part2/registry.md#etcd)： find the service via etcd.
  * [Consul](part2/registry.md#consul): find the service via consul.
  * [mDNS](part2/registry.md#mdns): find the service via mDNS that support local service discovery.
  * [In process](part2/registry.md#inprocess): find services in the same process. Clients call services in process and they don't commnucaite via TCP or UDP. It is convenient in test.


There is a sync rpcx example:

```go
package main

import (
	"context"
	"flag"
	"log"

	example "github.com/rpcx-ecosystem/rpcx-examples3"
	"github.com/smallnest/rpcx/client"
)

var (
	addr = flag.String("addr", "localhost:8972", "server address")
)

func main() {
	flag.Parse()

	d := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
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

### service governance (Fail mode and Load balancing)

In a large scale rpc system, there are a lot of service nodes to provide a same service. How clients select one of most appropriate node to call? If on call failed, clients should select another node or return the error immediately? They are fail mode issue and load balancing issue.


rpcx support fail mode:

- [Failfast](part3/failmode.md#failfast)：  return the error immediately if call failed
- [Failover](part3/failmode.md#failover) : select another node until reach max retries
- [Failtry](part3/failmode.md#failtry): select the same node and retry, until reach max retrirs


For the load balancing, rpcs provides a lot of selectors:

 * [Random](part3/random_selector.md): select the nodes randomly
  * [Roundrobin](part3/roundrobin_selector.md) : select the node by roundrobin
  * [Consistent hashing](part3/hash_selector.md): select the same node if servicePath, serviceMethod and Args. It uses the [jump consistent hash](https://arxiv.org/abs/1406.2294) and is very fast.
  * [Weighted](part3/weighted_selector.md): use the weighted which is configured in metadata of services(`weight=xxx`). The algorithm is like nginx implementation(smooth weighted algorithm)
  * [Network quality](part3/ping_selector.md): It uses the `ping` results. The better the quality of the network, the higher the probability of election of nodes
  * [Geography](part3/geo_selector.md): If there are multiple datacenters, clients perfer to connect service in the same datacenter.
  * [Customized Selector](part3/user_selector.md): If the above selectors are not suitable for you, you can use your customized selector. For example, one of rpcx customers has defined its selector for two datacenters because they can't use `Network quality` because of some limits.



There is an async rpcx example:

```go
package main

import (
	"context"
	"flag"
	"log"

	example "github.com/rpcx-ecosystem/rpcx-examples3"
	"github.com/smallnest/rpcx/client"
)

var (
	addr2 = flag.String("addr", "localhost:8972", "server address")
)

func main() {
	flag.Parse()

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
}
```

the client uses `Failtry` mode and selects the nodes randomly.

### Broadcast and Fork

And for rare cases, you can use `Broadcast`、`Fork` methods of XClient

```go
    Broadcast(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error
	Fork(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error
```

`Broadcast` sends requests to all servers and Success only when all servers return OK.FailMode and SelectMode are meanless for this method.Please set timeout to avoid hanging.

`Fork` sends requests to all servers and Success once any one server returns OK.
FailMode and SelectMode are meanless for this method.


You can use `NewXClient` to get a xclient instance.

```go
func NewXClient(servicePath string, failMode FailMode, selectMode SelectMode, discovery ServiceDiscovery, option Option) XClient
```

`NewXClient` must use the service name as the first argument, then failmode, selector, discovery and other options.


