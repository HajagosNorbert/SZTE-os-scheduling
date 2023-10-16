.PHONY build
build:
	go build -o ./bin/proc-gen cmd/proc-gen/main.go & go build -o ./bin/simulate cmd/simulate/main.go 
