# Benchmark

 The benchmark code is at [rpcx-ecosystem/rpcx-benchmark](https://github.com/rpcx-ecosystem/rpcx-benchmark).

 Use the same test environment, the same test data and the same test parameters, test  grpc, rpcx, dubbo, motan, thrift and go-micro.

 Based on my prior test, dubbo, motan and go-micro have poor performance, so their latest test have not been listed here, you can use the benchmark code to test them.


 ## Test Logic

 Use protobuf as the codec for all test. The proto file is [benchmark.proto](https://github.com/rpcx-ecosystem/rpcx-benchmark/blob/master/rpcx/benchmark.proto):

 ```proto
 syntax = "proto2";
 
 package main;
 
 option optimize_for = SPEED;
 
 
 message BenchmarkMessage {
   required string field1 = 1;
   optional string field9 = 9;
   optional string field18 = 18;
   optional bool field80 = 80 [default=false];
   optional bool field81 = 81 [default=true];
   required int32 field2 = 2;
   required int32 field3 = 3;
   
   ......
 }
 ```

Client generates a requests by setting each fields and the size of this request is 518 bytes. 

Server receives this request and sets the first field to "OK" and the second field to 100, then Server returns this request to Client.


The below two options can set concurrency and the total requests.

```
var concurrency = flag.Int("c", 1, "concurrency")
var total = flag.Int("n", 1, "total requests for all clients")
```

## Test Environment

- CPU: Intel(R) Xeon(R) CPU E5-2630 v3 @ 2.40GHz, 32 cores
- Memory: 32G
- Go: 1.9.2
- OS: CentOS 7 / 3.10.0-229.el7.x86_64

Client and Server are installed on the same machine.


## Test Result

### TPS

for 5000 concurrency, rpcx can archive 176,894 transations/second TPS, but grpc-go only gets 105219 transations/second

|_concurrency_| RPCX | GRPC-GO|
|----|----|------|
|5000|176894|105219|
|2000|161660|108245|
|1000|148227|111351|
|100|145479|93447|


### Latency: mean time

|_concurrency_| RPCX | GRPC-GO|
|----|----|------|
|5000|27|47|
|2000|12|18|
|1000|6|8|
|100|0|1|


### Latency: median time

|_concurrency_| RPCX | GRPC-GO|
|----|----|------|
|5000|3|42|
|2000|7|15|
|1000|5|7|
|100|0|0|