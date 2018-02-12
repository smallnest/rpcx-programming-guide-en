# FailMode

In the distribution architecture such SOA or micro service, you can not guarantee services always work well. Maybe  the server is down or network is slow or broken. So you need to handle invocation failures.

rpcx has three fail modes to handle invocation failures and you can set it when create XClient.


FailMode is meaningful only for sync call (`XClient.Call`).


## Failfast {#failfast}

**Example:** [failfast](https://github.com/rpcx-ecosystem/rpcx-examples3/tree/master/failmode/failfast)


In this mode, rpcx will finished this call if it gets a failure. The failure is not business failure and business failures are normal responses. It may be caused by server crashing or network issues.


## Failover {#failover}

**Example:** [failover](https://github.com/rpcx-ecosystem/rpcx-examples3/tree/master/failmode/failover)


In this mode, rpcx will try another nodes if it gets a failure. it will try `retries` until the service returns normal response. The `retries` is defined in `Option` and it is `3` in the defaultOption.

## Failtry {#failtry}

**Example:** [failtry](https://github.com/rpcx-ecosystem/rpcx-examples3/tree/master/failmode/failtry)


In this mode, rpcx still try current selected nodes if it gets a failure. it will try `retries` until the service returns normal response. The `retries` is defined in `Option` and it is `3` in the defaultOption.

## Failbackup {#failbackup}

**Example:** [failbackup](https://github.com/rpcx-ecosystem/rpcx-examples3/tree/master/failmode/failbackup)

In this mode, rpcx will send another request to one of servers if the first request has not returned in a given time. rpcx uses the the fastest returned response of the two requests. The given time is configured in `Option.BackupLatency`.

If ypu want to learn more details about backup request pattern, you can read Jeff Dean's [Achieving Rapid Response Times in Large Online Services](https://static.googleusercontent.com/media/research.google.com/zh-CN//pubs/archive/44875.pdf)

> From http://highscalability.com/blog/2012/6/18/google-on-latency-tolerant-systems-making-a-predictable-whol.html
> Backup requests are the idea of sending requests out to multiple replicas, but in a particular way. Here’s the example for a read operation for a distributed file system client:
> 
> send request to first replica
> wait 2 ms, and send to second replica
> servers cancel request on other replica when starting read
> A request could wait in a queue stuck behind an expensive query or a packet could be dropped, so if a reply is not returned quickly other replicas are tried. Responses come back faster if requests sit in multiple queues.
