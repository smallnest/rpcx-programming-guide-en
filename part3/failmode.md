# FailMode

In the distribution architecture such SOA or micro service, you can not guarantee services always work well. Maybe  the server is down or network is slow or broken. So you need to handle invocation failures.

rpcx has three fail modes to handle invocation failures and you can set it when create XClient.


FailMode is meaningful only for sync call (`XClient.Call`).


## Failfast {#failfast}

In this mode, rpcx will finished this call if it gets a failure. The failure is not business failure and business failures are normal responses. It may be caused by server crashing or network issues.


## Failover {#failover}

This this mode, rpcx will try another nodes if it gets a failure. it will try `retries` until the service returns normal response. The `retries` is defined in `Option` and it is `3` in the defaultOption.

## Failtry {#failtry}

This this mode, rpcx still try current selected nodes if it gets a failure. it will try `retries` until the service returns normal response. The `retries` is defined in `Option` and it is `3` in the defaultOption.