### Replicated log 🚧

Find task prompt [here](https://fly.io/dist-sys
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
