### Unique ID generator

🚧 short explanation of multiple potential approaches 🚧

Notice that since [here](https://fly.io/dist-sys/2/) there is no
constraint on how the generated ID should look like, it's OK to 
take a shortcut and generate unique string (using the fact that each node has an ID).

However, if let's say IDs should be 64 bit integers, than this approach will not work.
Also, the first approach might not work in case machines have been added/removed from a cluster.

The true way is to assign ranges of ids (64 bit integers) to each node. When a range is exhausted a node
requests another free (un-claimed) range.
