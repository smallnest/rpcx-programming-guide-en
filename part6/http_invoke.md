# HTTP Invoke

A rpcx server can be accessed via HTTP. The HTTP headers are same to the approach to access [gateway](https://github.com/rpcx-ecosystem/rpcx-gateway).

Obviously, `Invoke via HTTP` can not achive the same high performance to `Invoke via TCP/UDP`, but you can use it for cross-language development easily.

You can write a http client in your familiar programming language to access rpcx services directly. Requests and responses are in http body whatever codeces they are using. Some extra headers must be added into HTTP headers so that rpcx servers can parse services and methods.

The below is a an example:

```go
func main() {
	cc := &codec.MsgpackCodec{}

	args := &Args{
		A: 10,
		B: 20,
	}

	data, _ := cc.Encode(args)

	req, err := http.NewRequest("POST", "http://127.0.0.1:8972/", bytes.NewReader(data))
	if err != nil {
		log.Fatal("failed to create request: ", err)
		return
	}

	h := req.Header
	h.Set(gateway.XMessageID, "10000")
	h.Set(gateway.XMessageType, "0")
	h.Set(gateway.XSerializeType, "3")
	h.Set(gateway.XServicePath, "Arith")
	h.Set(gateway.XServiceMethod, "Mul")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("failed to call: ", err)
	}
	defer res.Body.Close()

	// handle http response
	replyData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("failed to read response: ", err)
	}

	reply := &Reply{}
	err = cc.Decode(replyData, reply)
	if err != nil {
		log.Fatal("failed to decode reply: ", err)
	}

	log.Printf("%d * %d = %d", args.A, args.B, reply.C)
}
```