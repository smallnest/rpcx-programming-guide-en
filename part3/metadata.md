# Metadata

**Example:** [metadata](https://github.com/rpcx-ecosystem/rpcx-examples3/tree/master/metadata)

Clients and servers can communicate with extra data.

Metadata is a key-vaue pair and key/value must be string.

## Client

If you want to send metadata, you **must** set `share.ReqMetaDataKey` in context.
If you want to receive metadata, you **must** set `share.ResMetaDataKey` in context.

```go client.go
	reply := &example.Reply{}
	ctx := context.WithValue(context.Background(), share.ReqMetaDataKey, map[string]string{"aaa": "from client"})
	ctx = context.WithValue(ctx, share.ResMetaDataKey, make(map[string]string))
	err := xclient.Call(ctx, "Mul", args, reply)
```

## Server

Server can get `share.ReqMetaDataKey` and `share.ResMetaDataKey` from context:

```go server.go
	reqMeta := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	resMeta := ctx.Value(share.ResMetaDataKey).(map[string]string)
```