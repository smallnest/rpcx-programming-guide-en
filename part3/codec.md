# Codec

Currently rpcx supports four type of codec:


```
// SerializeType defines serialization type of payload.
type SerializeType byte

const (
	// SerializeNone uses raw []byte and don't serialize/deserialize
	SerializeNone SerializeType = iota
	// JSON for payload.
	JSON
	// ProtoBuffer for payload.
	ProtoBuffer
	// MsgPack for payload
	MsgPack
)
```


Service use the same codec with clients. Clients set SerializeType in options. The default option use msgpack.

```go
var DefaultOption = Option{
	Retries:        3,
	RPCPath:        share.DefaultRPCPath,
	ConnectTimeout: 10 * time.Second,
	Breaker:        CircuitBreaker,
	SerializeType:  protocol.MsgPack,
	CompressType:   protocol.None,
}
```

The option is set when create XClient:

```go
func NewXClient(servicePath string, failMode FailMode, selectMode SelectMode, discovery ServiceDiscovery, option Option) 
```

## SerializeNone

No codec for data and use raw slice of bytes.

Clients and serices can encode/decode payload by their customized codec, for example, [Avro](https://github.com/linkedin/goavro).



## JSON

JSON is a general data format and can be used by many programming language.

But its performance is not better than protobuf and messagepack.


## Protobuf

[Protobuf](https://developers.google.com/protocol-buffers/) is a performant codec by Google and it is used in grpc and many projects.


## MsgPack

[messagepack](https://msgpack.org/index.html) is another performant codec and it is cross-platform too.