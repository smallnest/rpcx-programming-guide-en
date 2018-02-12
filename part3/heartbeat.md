# Heartbeat

**Example:** [heartbeat](https://github.com/rpcx-ecosystem/rpcx-examples3/tree/master/heartbeat)

You can set automatic heartbeat ping/pong to keep connection alive.

Rpcx servers can handle heartbeat requests automatically and you don't need to do more configurations.

Clients should enable heartbeat option in Client option and set the heartbeat interval:

```go
 	option := client.DefaultOption
	option.Heartbeat = true
	option.HeartbeatInterval = time.Second
```