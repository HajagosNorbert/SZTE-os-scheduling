package main

import (
	"flag"
	. "github.com/HajagosNorbert/SZTE-os-scheduling/internal/simulation"
)

func main() {
	procs := ReadProcs()
	alg := chooseAlg()
	res := alg(procs)
	CreateResultReport(res)
}

func chooseAlg() func([]Proc) SimResult {
	alg := flag.String("a", "fcfs", "Choose scheduling algorithm. possible values: \n\"fcfs\" - First Come First Serve \n\"sjr\" - Shortest Job First\n")
	flag.Parse()
	switch *alg {
	case "fcfs":
		return FirstComeFirstServe
	case "sjr":
		return SortestJobFirst
	default:
		return FirstComeFirstServe
	}
}
