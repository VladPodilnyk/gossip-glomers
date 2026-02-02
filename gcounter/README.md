### Grow-Only counter
In this challenge, the task is to implement stateless grow-only counter.
For more information, please check [this](https://fly.io/dist-sys/4/) page.

During this challenge, it's also allowed to use [sequentially-consistent ](https://jepsen.io/consistency/models/sequential) KV store provided by maelstrom.

So, the solution here is to store a counter for each node in our KV store. This way each node can receive messages and update its counter independently.
When a client requests a counter value, a node should get all counters from the KV store and sum them up.

Overall sounds quite straightforward, especially the read part.
However the `add` operation is a bit more tricky. We need to ensure that the counter is incremented atomically and that our data is not violated during concurrent updates. 
To solve this problem we can leverage compare-and-swap operation provided by maelstrom.SeqKV. Sounds easy.. right?

Well, actually the signature of the compare-and-swap operation looks like this:
```go
func (kv *KV) CompareAndSwap(ctx context.Context, key string, from, to any, createIfNotExists bool) error
````
It requires users to pass previous value and if this value is the actual one, then SeqKV performs the operation.
This operation is definitely useful, but it introduces another problem. Before using it, first we need to get the current value of the "local" counter. Which essentially means that the whole update process is not atomic and prone to race conditions.

Now to solve this "tiny" issue 😅, it's sufficient enough to use some retry logic in case a node fails to update the counter due to stale data (e.g. another request has passed and updated the counter).
In my solution, I've just used a simple loop that retries the operation until it succeeds. Although it's not the best solution for production use-cases, it's a good enough solution for this challenge.

Notice, that the solution doesn't use locking mechanisms and instead implements a what is called optimistic concurrency control updates where counter values serve as versions and hence we can track if a node has stale or recent data. Be aware, that OCC is not a silver bullet and doesn't fit systems with high contention, as it can lead to a lot of retries and performance degradation.
