### Replicated log

Find task prompt [here](https://fly.io/dist-sys)
The task is split into 3 parts.

#### Part 1: Single-node Kafka-Style log
The one is quite straightforward. No need to implement multiple nodes.
A single node has defined handlers and is responsible for storing values and managing offsets.
For the first part I went with the most basic solution, the log is stored in memory. The log state also
tracks offsets and last committed offsets per key.

The log state is protected by a mutex to ensure safe concurrent access. This was enough to pass tests.

#### Part 2&3: Multi-node Kafka-Style log / Efficient Kafka challenge
The idea for 2&3 is to implement a distributed log using lin-kv that guarantees linearizability.
Lin-kv should be used for orchestration and coordination. The log itself should be stored locally on disk on each node in the cluster. The leader node should be responsible for handling writes. Hence followers should forward writes to the leader. Here is how everything works together:

- Leader election happens during cluster initialization. Each node in a cluster tries to claim its right to
be a leader. The first node which successfully does so becomes the leader. In this challenge we don't need to
worry about failover or making re-election in case the leader dies hence the strategy as simple as possible.
- Leader is responsible for ALL writes in the system. While reads can happen from different nodes.
- When a `send` request lands on a regular node, it should proxy this request to the leader. Leader then persists the message and broadcasts it to all nodes in the system. For this particular exercises, I didn't bother with un-reliable network and other failures in this process hence leader just wait until all messages are propagated to all nodes. However this approach is not feasible in real world. To make it more production-like, leader should wait for the majority of nodes to have the correct and up-to-date information. But then we have to deal with a situation when an unavailable node comes back online.
