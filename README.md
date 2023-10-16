# Process Scheduling algorithm simulation

I measure various aspects of different process schedulers found in operating systems with this simulation.
These processes have a blocking I/O operations that also take some time to complete. Everything is mesured in CPU ticks, or with a simple integer counter for the numer of context switches.

I implement the 
 - Scheduling algorithms, of which these are done
    - First Come First Serve
    - Shortest Job Remaining
 - Methods to generate the processes and their blocking I/O operations probablistically
 - Reporting system that gives insights into the result of the simulation, with data like
    - context switches, idle ticks, latency, throughput

### What you need to run it
 - the **go toolchain** (I specified 1.21.1, but try modifying the go.mod file to your version if you can't run it)
 - GNU make (optional)

### How to run it ASAP

 The default process generation and simulation, with no parameterization can be run simply like this
```sh
make
```

#### How to use it 

There are 2 executables (at the moment result visualization is = 0). One for producing the input, one for processing it.

See available algorithms with:
```sh
go run cmd/simulate/main.go --help
```

For generating the processes and the I/O operations associated with them, see the options with:
```sh
go run cmd/proc-gen/main.go --help
```
You can use these two execuatbles together with some basic command piping to make them work well with eachother. They are only separated in case you want to generate your own input.json file.

```sh
go run cmd/proc-gen/main.go | go run cmd/simulate/main.go
```
Take a look at the `run-example` target of the [Makefile](./Makefile).

Everything writes and reads from standard input and standard output, so if you are interested in the generated json input, one way is to do this:
```sh
go run cmd/proc-gen/main.go | tee input.json | go run cmd/simulate/main.go
```
