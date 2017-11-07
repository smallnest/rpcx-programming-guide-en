# Selector

In a scaled system, there are many nodes to provide same service. They are deployed at the datacenter or multiple datacenter.

How do clients select a node to invoke? You can use a selector to do this and it is just like a load balancing with routing algorithm.

rpcx supports multiple selectors and  you can set it when you create XClient.



## RandomSelect {#random_selector}

This selector will pick a node randomly.


## Roundrobin {#roundrobin_selector}

select the node in a round-robin fashion.

## WeightedRoundRobin {#weighted_selector}

use the [smooth weighted round-robin balancing](https://github.com/phusion/nginx/commit/27e94984486058d73157038f7950a0a36ecc6e35), implemented the same algorithm in Nginx.


For edge case weights like { 5, 1, 1 } we now produce { a, a, b, a, c, a, a }
sequence instead of { c, b, a, a, a, a, a } produced previously.

Algorithm is as follows: on each peer selection we increase current_weight
of each eligible peer by its weight, select peer with greatest current_weight
and reduce its current_weight by total number of weight points distributed
among peers.


## WeightedICMP {#ping_selector}
use the result of ping (ICMP) to set weight of each node. The shorter ping time, the higher weight of the node.

assume `t` is cost time of ping:

-  weight=191       if t <= 10
-  weight=201 -t    if 10 < t <=200
-  weight=1         if 200 < t < 1000
-  weight=0         if t >= 1000


## ConsistentHash {#hash_selector}

Use [JumpConsistentHash](https://arxiv.org/abs/1406.2294) to router the same node for the same servicePath, serviceMethod and arguments. It is a fast consistent hash algorithm.

If nodes have changed, the algorithm will re-calculate consistent hash.

## Geography {#geo_selector}

We want to clients to select a closest service to invoke.

Services must set the geo location in their metadata.

If they are multiple closest services, select one randomly.

You should use the below method set geo location of clients and it will change selector to Geo selector.

```go
func (c *xClient) ConfigGeoSelector(latitude, longitude float64)
```

## SelectByUser Selector {#user_selector}

If the above selectors are not suitable for you, you can develop a customized selector.

The selector is an interface:

```go
type Selector interface {
	Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string
	UpdateServer(servers map[string]string)
}
```


-`Select`: defines the select algorithm.
-`UpdateServer`: clients init the nodes and update if nodes change.
