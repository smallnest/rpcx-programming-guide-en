# Metrics plugin

**Example:** [metrics](https://github.com/rpcx-ecosystem/rpcx-examples3/tree/master/metrics)


Metrics plugin uses [go-metrics](github.com/rcrowley/go-metrics) to calculate metrics of the server.

It contains some metrics:

1. serviceCounter
2. clientMeter
3. "service_"+servicePath+"."+serviceMethod+"_Read_Qps"
4. "service_"+servicePath+"."+serviceMethod+"_Write_Qps"
5. "service_"+servicePath+"."+serviceMethod+"_CallTime"
