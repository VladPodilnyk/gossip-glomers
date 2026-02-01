### Broadcast
In this part the task is to implement a broadcast system that gossip messages
between nodes. The task is split into three parts that progressively increase in difficulty.
Find the prompt for the problem [here](https://fly.io/dist-sys/3a/)

__Note__: In this task it's OK to rely on the fact that all nodes are connected to each other, and there are
no separate clusters of nodes.

#### Part A
This part is quite simple. We need to implement a few RPC handlers on our server.
Since we begin with 1 node only, concurrency and topology do not matter.

#### Part B
The part B takes it one step further and we tasked to add support for multiple nodes.
The solution from the part A requires a bunch of changes. A few obvious are:
- Implementing the "topology" handler
- Implementing the "broadcast" handler.

While implementing the "topology" handler is fairly straightforward, since a node gets this message only once
during cluster initialization, there is not need to handle concurrency.
The "broadcast" handler is a bit more tricky. First we need to use locking to ensure that node's state is consistent and updated correctly. Then we need to figure out how to propagate each message through the cluster.

To do this, when a node receives a broadcast message, it should send "broadcast" message using `.Send` RPC to all its neighbors. To prevent loops we need to check the following conditions:
- The message should not be sent to the node that sent the original message.
- The message should not be sent to a node that has already received the message.

**Implementation note**: For 3b, the values are stored in an array of float64. While for this toy problem this is enough and not bad, it is not an efficient way to store the values. A more efficient way would be to use a set or in the case of Golang - map. This way lookups are fast and efficient.

#### Part C
🚧

#### Part D
🚧

#### Part E
🚧
