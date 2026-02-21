### Replicated log 🚧

Find task prompt [here](https://fly.io/dist-sys)
The task is split into 3 parts.

#### Part 1: Single-node Kafka-Style log
The one is quite straightforward. No need to implement multiple nodes.
A single node has defined handlers and is responsible for storing values and managing offsets.
For the first part I went with the most basic solution, the log is stored in memory. The log state also
tracks offsets and last committed offsets per key.

The log state is protected by a mutex to ensure safe concurrent access. This was enough to pass tests.

#### Part 2: Multi-node Kafka-Style log
🚧

#### Part 3: Efficient Kafka challenge
🚧

__TODO__:
The idea for 2&3 is to implement a distributed log using lin-kv that guarantees linearizability.
Lin-kv should be used for orchestration and coordination. The log itself should be stored locally on disk on each node in the cluster. The leader node should be responsible for handling writes. Hence followers should forward writes to the leader.

- [x] implement persistence and store logs for each node locally on disk
- [ ] implement leader management:
    - auto select leader (can be just a random node or the first node that successfully put its claim to link-kv)
    - implement failover (heartbeats and last updated timestamps or on failed broadcasts messages)
- [ ] implement replication protocol (leader should wait for thea ACKs message from the majority of followers)
- [ ] return only committed values
