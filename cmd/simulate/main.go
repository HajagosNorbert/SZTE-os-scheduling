package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"

	. "github.com/HajagosNorbert/SZTE-os-scheduling/internal/simulation"
)

func main() {
	algFlag := ReadFlag()
	alg, algLongName := chooseAlg(algFlag)
	procs := readProcs()
	result := SimulateScheduling(procs, alg)
	CreateResultReport(procs, result, algLongName)
}

func readProcs() []Proc {
	procsJsonInput, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Could not read stdin until end of file.")
	}
	var procs []Proc
	if err := json.Unmarshal(procsJsonInput, &procs); err != nil {
		log.Fatalf("Input was not in the correct Json format, could not decode it. Execute `make run` at the root of the project, or take a look at the makefile to know how to create a proper input file.\n")
	}
	return procs
}

func ReadFlag() *string {
	alg := flag.String("a", AlgFcfs, `Choose scheduling algorithm. possible values: 
	"fcfs" - First Come First Serve 
	"sjr" - Shortest Job Remaining
	"lottery" - Lottery
	"srr" - Smart Round-Robin
	"rr" - Round - Robin`+"\n")
	flag.Parse()
	return alg
}

// Returns the scheduling algorithm function and the display name of the algorithm
func chooseAlg(alg *string) (func([]Proc, int) (int, bool), string) {
	switch *alg {
	case AlgFcfs:
		return FirstComeFirstServe, "First Come First Serve"
	case AlgLottery:
		return MakeLottery(), "Lottery"
	case AlgSjr:
		return ShortestJobRemaining, "Shortest Job Remaining"
	case AlgSrr:
		return MakeSmartRoundRobin(), "Smart Round - Robin"
	case AlgRoundRobin:
		return MakeRoundRobin(), "Round - Robin"
	default:
		log.Fatalf("Error: No algorithm implemented with value '%s'", *alg)
		return nil, "err"
	}
}
