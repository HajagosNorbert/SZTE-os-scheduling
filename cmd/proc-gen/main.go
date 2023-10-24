package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"time"

	"github.com/HajagosNorbert/SZTE-os-scheduling/internal/simulation"
)

type opts struct {
	procTicks int
	procCount int
	mean      int
	stddev    int
	maxIo     int
	seed      int64
}

func main() {
	opts := readFlag()
	r := rand.New(rand.NewSource(opts.seed))
	ioOpsPerProc := ioOpsCountPerProcNormalDist(opts.procTicks, float64(opts.stddev), float64(opts.mean), opts.procCount, r)

	procs := make([]simulation.Proc, opts.procCount)

	for i := 0; i < len(procs); i++ {
		ioOpCount:= ioOpsPerProc[i]
		priority := r.Intn(10)+1
		userId := r.Intn(5)
		procs[i] = simulation.Proc{TicksLeft: opts.procTicks, TotalTicks: opts.procTicks, SpawnedAt: 0, Priority: priority, UserId: userId}
		procs[i].IoOps = genIoOps(ioOpCount, opts.procTicks, opts.maxIo, r)
	}

	procsJson, err := json.MarshalIndent(procs, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf(string(procsJson))
}

func readFlag() opts {
	procTicks := flag.Int("procTicks", 100, "The amount of CPU ticks each process will take to complete.(waiting for I/O operations not included)")
	procCount := flag.Int("procCount", 5, "The amount of processes that will be started.")
	mean := flag.Int("ioMean", 5, "The mean of a normal distribution, responsible for the number of I/O operations a single process will have. It takes 1 CPU tick from the *procTicks* argument to start an I/O operation. If you set it so that the random generator produces more I/O operations than ticks in a process, the number of I/O operations will match the *procTick* argument.")
	stddev := flag.Int("ioStd", 2, "(Read *ioMean* before) The standard deviation from ioMean in a normal distribution.")
	maxIo := flag.Int("maxIoTick", 10, "For each process, it's I/O operations will take an evenly distributed random value between 1 and *maxIoTick* (inclusive) to finish.")
	seed := flag.Int64("seed", 0, "The random seed that the input generation will work with. If not set, it will be different on every run.")

	flag.Parse()
	if !isFlagPassed("seed") {
		*seed = time.Now().UnixNano()
	}
	options := opts{
		procTicks: *procTicks,
		procCount: *procCount,
		mean:      *mean,
		stddev:    *stddev,
		maxIo:     *maxIo,
		seed:      *seed,
	}
	return options
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func ioOpsCountPerProcNormalDist(maxIoOpsCount int, stddev float64, mean float64, procCount int, r *rand.Rand) []int {

	result := make([]int, procCount)
	min := 0

	for i := 0; i < procCount; i++ {
		value := int(r.NormFloat64()*stddev + mean)

		if value < min {
			value = min
		} else if value > maxIoOpsCount {
			value = maxIoOpsCount
		}
		result[i] = value
	}

	return result
}

func genIoOps(ioOpCount, procTickCount, maxIoOpTicks int, r *rand.Rand) []simulation.IoOp {
	ioOps := make([]simulation.IoOp, ioOpCount)
	if ioOpCount > procTickCount {
		ioOpCount = procTickCount
	}
	startTicks := r.Perm(procTickCount)[:ioOpCount]
	sort.Ints(startTicks)

	for i := 0; i < ioOpCount; i++ {
		ticksLeft := r.Intn(maxIoOpTicks) + 1
		ioOps[i] = simulation.IoOp{StartsAfter: startTicks[i], TicksLeft: ticksLeft}
	}

	return ioOps
}
