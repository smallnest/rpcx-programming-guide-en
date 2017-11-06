# Registry

Registry implements service discovery.

![](registry.png)

Services registered themself into registries and clients can find those service via registries.

There are some poplular products that can be used as registies, for example, zookeeper, etcd, consul. And many companies implement their registry center.

rpcx provides out-of-box registry for zookeeper, etcd, consul, mDNS, and it aslo provides peer-to-peer, peer-to-multiple registies. For the test purpose, rpcx aslo has a in-process registgry.


## peer2peer {#peer2peer}

Actually it has no registry. Clients connect services directly. It can be used in a small system and there is only one node for services. You can do few changes to change another registry when your system scales.

Server doesn't need to do more configurations.

Client uses the `Peer2PeerDiscovery` and only set the network and address of that service.

Since there is one node, the selector is meanless.

```
    d := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
    defer xclient.Close()
```

** Notice: rpcx uses `network@Host:port` format to indicates one service. the `network` can be `tcp`、`http`、`unix`、`quic` or `kcp`. The `Host` can be host name or ip address.**


`NewXClient` must use the service name as the first argument, then failmode, selector, discovery and other options.


## peer2multiple {#multiple}

If you has multiple services but has no a center registry. You can configure addresses of services in client programatically.

Server doesn't need to do more configurations.

Client uses the `MultipleServersDiscovery` and only set the network and address of that service.

```go
    d := client.NewMultipleServersDiscovery([]*client.KVPair{{Key: *addr1}, {Key: *addr2}})
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
    defer xclient.Close()
```

You must set services info and metadata in MultipleServersDiscovery。 If some services are added or removed, you can call `MultipleServersDiscovery.Update` to update service dynamically.

```
func (d *MultipleServersDiscovery) Update(pairs []*KVPair)
```

## zookeeper {#zookeeper}

Apache ZooKeeper is a software project of the Apache Software Foundation. It is essentially a distributed hierarchical key-value store, which is used to provide a distributed configuration service, synchronization service, and naming registry for large distributed systems.

Services must use `ZooKeeperRegisterPlugin` plugin to register their info into zookeeper, while clients must use `ZookeeperDiscovery` to get services info.

You must set ServiceAddress (`network@address`) and zookeeper address. You should set the basePath so that your services won't be conflict with other serviced developed by others.

You can set `Metrics` and rpcx will update throughputs periodically. If you want to update metrics you must add [Metrics] plugin in server.

You can set `UpdateInterval` to refresh the lease. rpcx sets the lease to `UpdateInterval * 3`. That means if server has not update the lease in  `UpdateInterval * 3`, this node will be removed from zookeeper.


```go server.go
func main() {
	flag.Parse()

	s := server.NewServer()
	addRegistryPlugin(s)

	s.RegisterName("Arith", new(example.Arith), "")
	s.Serve("tcp", *addr)
}

func addRegistryPlugin(s *server.Server) {

	r := &serverplugin.ZooKeeperRegisterPlugin{
		ServiceAddress:   "tcp@" + *addr,
		ZooKeeperServers: []string{*zkAddr},
		BasePath:         *basePath,
		Metrics:          metrics.NewRegistry(),
		UpdateInterval:   time.Minute,
	}
	err := r.Start()
	if err != nil {
		log.Fatal(err)
	}
	s.Plugins.Add(r)
}
```




Configuration for Client is simple.
You need to set the zookeeper address and the same basePath to service.

```go client.go
    d := client.NewZookeeperDiscovery(*basePath, "Arith",[]string{*zkAddr}, nil)
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
    defer xclient.Close()
```


## etcd {#etcd}

etcd is a distributed, reliable key-value store for the most critical data of a distributed system.

Services must use `EtcdRegisterPlugin` plugin to register their info into etcd, while clients must use `EtcdDiscovery` to get services info.

You must set ServiceAddress (`network@address`) and etcd address. You should set the basePath so that your services won't be conflict with other serviced developed by others.

You can set `Metrics` and rpcx will update throughputs periodically. If you want to update metrics you must add [Metrics] plugin in server.

You can set `UpdateInterval` to refresh the lease. rpcx sets the lease to `UpdateInterval * 3`. That means if server has not update the lease in  `UpdateInterval * 3`, this node will be removed from etcd.


```go server.go
func main() {
	flag.Parse()

	s := server.NewServer()
	addRegistryPlugin(s)

	s.RegisterName("Arith", new(example.Arith), "")
	s.Serve("tcp", *addr)
}

func addRegistryPlugin(s *server.Server) {

	r := &serverplugin.EtcdRegisterPlugin{
		ServiceAddress: "tcp@" + *addr,
		EtcdServers:    []string{*etcdAddr},
		BasePath:       *basePath,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Minute,
	}
	err := r.Start()
	if err != nil {
		log.Fatal(err)
	}
	s.Plugins.Add(r)
}
```




Configuration for Client is simple.
You need to set the etcd address and the same basePath to service.

```go client.go
    d := client.NewEtcdDiscovery(*basePath, "Arith",[]string{*etcdAddr}, nil)
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
    defer xclient.Close()
```


## consul {#consul}
Consul is a highly available and distributed service discovery and KV store designed with support for the modern data center to make distributed systems and configuration easy.

Services must use `ConsulRegisterPlugin` plugin to register their info into consul, while clients must use `ConsulDiscovery` to get services info.

You must set ServiceAddress (`network@address`) and consul address. You should set the basePath so that your services won't be conflict with other serviced developed by others.

You can set `Metrics` and rpcx will update throughputs periodically. If you want to update metrics you must add [Metrics] plugin in server.

You can set `UpdateInterval` to refresh the lease. rpcx sets the lease to `UpdateInterval * 3`. That means if server has not update the lease in  `UpdateInterval * 3`, this node will be removed from consul.


```go server.go
func main() {
	flag.Parse()

	s := server.NewServer()
	addRegistryPlugin(s)

	s.RegisterName("Arith", new(example.Arith), "")
	s.Serve("tcp", *addr)
}

func addRegistryPlugin(s *server.Server) {

	r := &serverplugin.ConsulRegisterPlugin{
		ServiceAddress: "tcp@" + *addr,
		ConsulServers:  []string{*consulAddr},
		BasePath:       *basePath,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Minute,
	}
	err := r.Start()
	if err != nil {
		log.Fatal(err)
	}
	s.Plugins.Add(r)
}
```




Configuration for Client is simple.
You need to set the etcd address and the same basePath to service.

```go client.go
    d := client.NewConsulDiscovery(*basePath, "Arith",[]string{*consulAddr}, nil)
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
    defer xclient.Close()
```



## mDNS {#mdns}
mDNS resolves host names to IP addresses within small networks that do not include a local name server. It is a zero-configuration service, using essentially the same programming interfaces, packet formats and operating semantics as the unicast Domain Name System (DNS). Although Stuart Cheshire designed mDNS to be stand-alone capable, it can work in concert with unicast DNS servers.

The mDNS protocol is published as [RFC 6762](https://tools.ietf.org/html/rfc6762), uses IP multicast User Datagram Protocol (UDP) packets, and is implemented by the Apple Bonjour, Spotify Connect, Philips Hue, Google Chromecast, and open source Avahi (software) software packages. Android contains an mDNS implementation. mDNS has also been implemented in Windows 10, but its use is limited to discovering networked printers.

mDNS can work in conjunction with DNS Service Discovery (DNS-SD), a companion zero-configuration technique specified separately in [RFC 6763](https://tools.ietf.org/html/rfc6763).


Services must use `MDNSRegisterPlugin` plugin to register their info into mDNS, while clients must use `MDNSDiscovery` to get services info.

You must set ServiceAddress (`network@address`). You should set the basePath so that your services won't be conflict with other serviced developed by others.

You can set `Metrics` and rpcx will update throughputs periodically. If you want to update metrics you must add [Metrics] plugin in server.

You can set `UpdateInterval` to refresh tps. If the server is down, it will be removed from mDNS automatically.


```go server.go
func main() {
	flag.Parse()

	s := server.NewServer()
	addRegistryPlugin(s)

	s.RegisterName("Arith", new(example.Arith), "")
	s.Serve("tcp", *addr)
}

func addRegistryPlugin(s *server.Server) {

	r := serverplugin.NewMDNSRegisterPlugin("tcp@"+*addr, 8972, metrics.NewRegistry(), time.Minute, "")
	err := r.Start()
	if err != nil {
		log.Fatal(err)
	}
	s.Plugins.Add(r)
}
```




Configuration for Client is simple.
You need to set the etcd address and the same basePath to service.

```go client.go
    d := client.NewMDNSDiscovery("Arith", 10*time.Second, 10*time.Second, "")
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()
```

## In process {#inprocess}
Inprocess registry is a registry for test purpose. You won't use it in production.

When you develop clients, you want to create mock services to integration test. You you can use Inprocess registry and change it to other registries when clients are deployed in production environment.

This plugin uses `reflect` to invoke services in the same process.


You add `client.InprocessClient` into services' plugin container. And clients use `InpreocessDiscovery` to find registered services to invoke.


```go server.go
func main() {
	flag.Parse()

	s := server.NewServer()
	addRegistryPlugin(s)

	s.RegisterName("Arith", new(example.Arith), "")

	go func() {
		s.Serve("tcp", *addr)
	}()

	d := client.NewInpreocessDiscovery()
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	args := &example.Args{
		A: 10,
		B: 20,
	}

	for i := 0; i < 100; i++ {

		reply := &example.Reply{}
		err := xclient.Call(context.Background(), "Mul", args, reply)
		if err != nil {
			log.Fatalf("failed to call: %v", err)
		}

		log.Printf("%d * %d = %d", args.A, args.B, reply.C)

	}
}

func addRegistryPlugin(s *server.Server) {

	r := client.InprocessClient
	s.Plugins.Add(r)
}
```