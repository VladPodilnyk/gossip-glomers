### Totally available transactions

The task it to implement a simple key/value storage which implements transactions
with some traits. Find full task prompt [here](https://fly.io/dist-sys/6a/).

This problem is split into 3 parts.

#### Part A: Single-Node, Totally-Available Transactions
This one is quite straightforward and overall doesn't require much explanation.

#### Part B&C: Totally Available, Read Uncommitted & Read Committed.
This time we have to fulfill specific transaction semantics.

Read uncommitted is the most permissive consistency model and it forbids only "dirty writes".
Meaning writes from different transaction "compete"/interleave with each other. However, notice that 
there is no constrains on eventual ordering of values.

Solution for this part almost the same as for the previous part. When transaction is executed 
a mutex is taken to ensure that only operation from the current transaction modify the storage.
All write operations are bundled together and broadcasted to other nodes in a cluster using `sync`
message. Since our solution should withstand network partitions, I added a simple retry mechanism on top of
message broadcasting.

And, as it turned out the same approach passes tests for part C, where we are tasked to implement Read Committed semantics, which means that only committed state of transactions can be observed.
