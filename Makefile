.PHONY: run
run: build
	./bin/proc-gen | ./bin/simulate

.PHONY: srr
srr: build
	cat input1.json | ./bin/simulate -a srr 

.PHONY: srr
srr-test: build
	cat test-input.json | ./bin/simulate -a srr 

.PHONY: run-example
run-example: build
	./bin/proc-gen -ioMean 3 -ioStd 1 -maxIoTick 4 -procCount 3 -procTicks 10 -seed 42 | ./bin/simulate -a sjr

build:
	go build -o ./bin/proc-gen cmd/proc-gen/main.go & go build -o ./bin/simulate cmd/simulate/main.go
