.PHONY: echo
echo:
	./maelstrom-runner/maelstrom test -w echo --bin ~/go/bin/maelstrom-echo --node-count 1 --time-limit 10

.PHONY: unique-ids
unique-ids:
	./maelstrom-runner/maelstrom test -w unique-ids --bin ~/go/bin/maelstrom-unique-ids --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition

.PHONY: broadcast
broadcast:
	./maelstrom-runner/maelstrom test -w broadcast --bin ~/go/bin/maelstrom-broadcast --node-count 1 --time-limit 20 --rate 10
