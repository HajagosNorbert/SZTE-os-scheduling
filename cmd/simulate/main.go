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
	alg, algName := chooseAlg(algFlag)
	procsJsonInput, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Could not read stdin until end of file.")
	}
	var procs []Proc
	if err := json.Unmarshal(procsJsonInput, &procs); err != nil {
		log.Fatalf("Input was not in the correct Json format, could not decode it. Execute `make run` at the root of the project, or take a look at the makefile to know how to simulate.\n")
	}
	// procs := ReadProcs()
	result := SimulateScheduling(procs, alg)

	CreateResultReport(result, algName)
}

func ReadFlag() *string {
	alg := flag.String("a", AlgFcfs, "Choose scheduling algorithm. possible values: \n\"fcfs\" - First Come First Serve \n\"sjr\" - Shortest Job Remaining\n")
	flag.Parse()
	return alg
}

// Returns the scheduling algorithm function and the display name of the algorithm
func chooseAlg(alg *string) (func([]Proc, int) (int, bool), string) {
	switch *alg {
	case AlgFcfs:
		return FirstComeFirstServe, "First Come First Serve"
	case AlgSjr:
		return ShortestJobRemaining, "Shortest Job Remaining"
	default:
		log.Fatalf("Error: No algorithm implemented with value '%s'", *alg)
		return nil, "err"
	}
}
