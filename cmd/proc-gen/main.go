package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"time"

	"github.com/HajagosNorbert/SZTE-os-scheduling/internal/simulation"
)

func main() {
	procTickCount := 100
	procCount := 5
	stddev := 1.0
	mean := 5.0

	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))

	ioOpsPerProc := ioOpsCountPerProcNormalDist(procTickCount, stddev, mean, procCount, r)
	maxIoOpTick := 10

	procs := make([]simulation.Proc, procCount)
	for i, ioOpCount := range ioOpsPerProc {
		procs[i] = simulation.Proc{TicksLeft: procTickCount, TotalTicks: procTickCount, SpawnedAt: 0}
		procs[i].IoOps = genIoOps(ioOpCount, procTickCount, maxIoOpTick, r)
	}

	procsJson, err := json.MarshalIndent(procs, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf(string(procsJson))
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
		ticksLeft := r.Intn(maxIoOpTicks-1) + 1
		ioOps[i] = simulation.IoOp{StartsAfter: startTicks[i], TicksLeft: ticksLeft}
	}

	return ioOps
}
