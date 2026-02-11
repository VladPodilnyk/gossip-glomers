.PHONY: echo
echo:
	./maelstrom-runner/maelstrom test -w echo --bin ~/go/bin/echo --node-count 1 --time-limit 10

.PHONY: unique-ids
unique-ids:
	./maelstrom-runner/maelstrom test -w unique-ids --bin ~/go/bin/unique-ids --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition

.PHONY: broadcast
broadcast: clean-logs
	cd broadcast && go install . && cd ../ && ./maelstrom-runner/maelstrom test -w broadcast --bin ~/go/bin/broadcast --node-count 25 --time-limit 20 --rate 100 --latency 100

.PHONY: gcounter
gcounter: clean-logs
	cd gcounter && go install . && cd ../ && ./maelstrom-runner/maelstrom test -w g-counter --bin ~/go/bin/gcounter --node-count 3 --rate 100 --time-limit 20 --nemesis partition

.PHONE: replicated-log
replicated-log: clean-logs
	cd replicated-log && go install . && cd ../ && ./maelstrom-runner/maelstrom test -w kafka --bin ~/go/bin/replicated-log --node-count 1 --concurrency 2n --time-limit 20 --rate 1000

.PHONY: clean-logs
clean-logs:
	rm -rf ${HOME}/maelstrom
