# Selector

In a scaled system, there are many nodes to provide same service. They are deployed at the datacenter or multiple datacenter.

How do clients select a node to invoke? You can use a selector to do this and it is just like a load balancing with routing algorithm.

rpcx supports multiple selectors and  you can set it when you create XClient.



## random_selector {#random_selector}

This selector will pick a node randomly.