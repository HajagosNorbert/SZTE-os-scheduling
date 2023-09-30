package main

import (
	"flag"
	"log"

	. "github.com/HajagosNorbert/SZTE-os-scheduling/internal/simulation"
)

func main() {
	algFlag := ReadFlag()
	alg, algName := chooseAlg(algFlag)
	procs := ReadProcs()
	result := alg(procs)
	CreateResultReport(result, algName)
}

func ReadFlag() *string {
	alg := flag.String("a", AlgFcfs, "Choose scheduling algorithm. possible values: \n\"fcfs\" - First Come First Serve \n\"sjr\" - Shortest Job First\n")
	flag.Parse()
	return alg
}

// Returns the scheduling algorithm function and the display name of the algorithm
func chooseAlg(alg *string) (func([]Proc) SimResult, string) {
	switch *alg {
	case AlgFcfs:
		return FirstComeFirstServe, "First Come First Serve"
	case AlgSjr:
		return SortestJobFirst, "Shortest Job First"
	default:
		log.Fatalf("Error: No algorithm implemented with value '%s'", *alg)
		return nil, "err"
	}
}
