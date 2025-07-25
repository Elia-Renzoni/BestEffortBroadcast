# BestEffortBroadcast
Best Effort Broadcast Algorithm Implemented in Golang

## Intro
A common operation in distributed systems is the dissemination of messages to the entire cluster.
Broadcasting is a way to spread information across a cluster via point-to-point TCP/UDP links.
The simplest and most unreliable broadcast algorithm is the Best Effort Broadcast.

This algorithm is based on a fail-silent distributed system model, where: <br>

* The process abstraction is crash-stop.

* The link abstraction is perfect (i.e., exactly-once delivery, no duplication).

The algorithm allows a sender to disseminate a message to the entire cluster, but it does not guarantee anything if the sender crashes.
In such cases, the agreement property is violated due to an inconsistent data view across nodes.

## Test
First starts a cluster: <br>
```
go build best_effort_broadcast.go
./best_effort_broadcast -host=<your host> -port=<a different port for each node>
```

Then Start the client: <br>
```
cd BestEffortBroadcast/client
go beb_client.go
```
The client sends a SET request to the distributed key-value store.
The node that receives the message will eventually broadcast it to the entire cluster.
Note: The client sends all messages to the seed node, which is statically configured as 127.0.0.1:7878.
All nodes in the distributed system join the cluster by contacting the seed node.