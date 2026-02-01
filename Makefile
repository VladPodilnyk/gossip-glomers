.PHONY: echo
echo:
	./maelstrom-runner/maelstrom test -w echo --bin ~/go/bin/echo --node-count 1 --time-limit 10

.PHONY: unique-ids
unique-ids:
	./maelstrom-runner/maelstrom test -w unique-ids --bin ~/go/bin/unique-ids --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition

.PHONY: broadcast
broadcast: clean-logs
	cd broadcast && go install . && cd ../ && ./maelstrom-runner/maelstrom test -w broadcast --bin ~/go/bin/broadcast --node-count 5 --time-limit 20 --rate 10 --nemesis partition

.PHONY: clean-logs
	rm /tmp/maelstrom*
