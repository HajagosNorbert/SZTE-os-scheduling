package simulation

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
)

func FirstComeFirstServe(procs []Proc, currProcIdx int) (int, bool) {

	if isValidRunningProcIdx(procs, currProcIdx) {
		return currProcIdx, true
	}
	for i, p := range procs {
		if p.State == Ready {
			return i, true
		}
	}
	return -1, false
}

func MakeLottery() func([]Proc, int) (int, bool) {
	lottery := func(procs []Proc, currProcIdx int) (int, bool) {
		totalTickets := 0
		for i := 0; i < len(procs); i++ {
			if procs[i].State == Ready || procs[i].State == Running {
				totalTickets += procs[i].Priority
			}
		}
		if totalTickets == 0 {
			return -1, false
		}

		winnerTicket := rand.Intn(totalTickets) + 1
		for i := 0; i < len(procs); i++ {
			if procs[i].State == Ready || procs[i].State == Running {
				if winnerTicket <= procs[i].Priority {
					return i, true
				}
				winnerTicket -= procs[i].Priority
			}
		}
		return -1, false
	}

	return lottery
}

func MakeRoundRobin() func([]Proc, int) (int, bool) {
	ticksRemainingForCurrent := 0

	alg := func(procs []Proc, currProcIdx int) (int, bool) {
		var procIdx int

		if isValidRunningProcIdx(procs, currProcIdx) {
			procIdx = currProcIdx
			if ticksRemainingForCurrent > 0 {
				ticksRemainingForCurrent--
				return procIdx, true
			}
		} else {
			return FirstComeFirstServe(procs, currProcIdx)
		}

		for i := 1; i < len(procs); i++ {
			procIdxCandidate := (procIdx + i) % len(procs)
			if procs[procIdxCandidate].State == Ready {
				ticksRemainingForCurrent = int(procs[procIdxCandidate].TicksLeft / 8)

				assignedTicks := int(procs[procIdxCandidate].TicksLeft / 8)
				if assignedTicks < 1 {
					ticksRemainingForCurrent = procs[procIdxCandidate].TicksLeft
				} else {
					ticksRemainingForCurrent = assignedTicks
				}

				return procIdxCandidate, true
			}
		}
		return procIdx, true
	}
	return alg
}

type procInCycle struct {
	ticksLeft   int
	originalIdx int
	// TotalTicks  int
	// Priority    int
	// UserId      int
	// State       ProcState
	// IoOps       []IoOp
}

// Handle ticks left being the same
func MakeSmartRoundRobin() func([]Proc, int) (int, bool) {
	ticksRemainingForCurrent := 0
	currCycleProcIdx := -1
	smartTimeQuanum := 0
	var cycleProcs []procInCycle
	insideCycle := false
	alg := func(procs []Proc, currProcIdx int) (int, bool) {
		// if nem running már a currProcIdx:q

		if currProcIdx != -1 {
			if procs[currProcIdx].State == Blocked || procs[currProcIdx].State == Terminated {
				ticksRemainingForCurrent = 0
				if currCycleProcIdx < len(cycleProcs) -1 {
					fmt.Printf("currCycleProcIdx: %+v\n", currCycleProcIdx)
					fmt.Printf("cycleProcs: %+v\n", cycleProcs)
					fmt.Printf("%+v == %+v \n", cycleProcs[currCycleProcIdx], currProcIdx)
					cycleProcs[currCycleProcIdx] = cycleProcs[len(cycleProcs)-1]
					cycleProcs = cycleProcs[0 : len(cycleProcs)-1]
					currCycleProcIdx--
				}
			}
		}

		endOfCycle := insideCycle && ticksRemainingForCurrent == 0 && currCycleProcIdx <= len(cycleProcs)-1
		// fmt.Printf("%+v", ticksRemainingForCurrnt)
		//cycle never ends
		if endOfCycle {
			// fmt.Printf("%+v\n", currProcIdx)
			fmt.Printf("End of cycle\n")
			insideCycle = false
		} else {
			fmt.Printf("Not End of cycle\n")
		}
		if !insideCycle {
			//one new cycle
			fmt.Printf("new cycle with proc: %+v\n", procs)
			if isValidRunningProcIdx(procs, currProcIdx) {
				procs[currProcIdx].State = Ready
			}

			cycleProcs = cycleProcs[:0]
			readyCount := 0
			for i, proc := range procs {
				if proc.State == Ready {
					cycleProcs = append(cycleProcs, procInCycle{originalIdx: i, ticksLeft: procs[i].TicksLeft})
					readyCount++
				}
			}
			// cycleProcs = cycleProcs[:readyCount]
			fmt.Printf("cycleProcs: %+v\n", cycleProcs)
			sort.Slice(cycleProcs, func(i, j int) bool {
				return cycleProcs[i].ticksLeft < cycleProcs[j].ticksLeft
			})

			fmt.Printf("1 cycleProcs: %+v\n", cycleProcs)
			if len(cycleProcs) == 0 {
				// fmt.Printf("ended with procs:\n%+v\n\n", procs)
				fmt.Printf("ended with procs:\n%+v\n\n", procs)
				return -1, false
			}

			insideCycle = true
			if len(cycleProcs) == 1 {
				currCycleProcIdx = 0
				ticksRemainingForCurrent = cycleProcs[currCycleProcIdx].ticksLeft
				fmt.Printf("returning this proc: %+v\n\n", procs[cycleProcs[currCycleProcIdx].originalIdx])
				return cycleProcs[currCycleProcIdx].originalIdx, true
			} else {
				currCycleProcIdx = -1
			}

			fmt.Printf("2 cyclaProcs: %+v\n", cycleProcs)
			ticksLeftSum := 0
			for i := 0; i < len(cycleProcs)-1; i++ {
				// fmt.Printf("%+v\n", cycleProcs)
				ticksLeftSum += cycleProcs[i+1].ticksLeft - cycleProcs[i].ticksLeft
			}

			// fmt.Printf("%+v\n", ticksLeftSum)
			smartTimeQuanum = int(math.Round(float64(ticksLeftSum) / (float64(len(procs)) - 1.0)))
			if smartTimeQuanum == 0 {
				assignedTicks := int(cycleProcs[0].ticksLeft / 8)
				if assignedTicks < 1 {
					smartTimeQuanum = cycleProcs[0].ticksLeft
				} else {
					smartTimeQuanum = assignedTicks
				}
			}
			fmt.Printf("STQ: %+v\n", smartTimeQuanum)
			// println(currCycleProcIdx)
		}

		if ticksRemainingForCurrent > 0 {
			ticksRemainingForCurrent--
			fmt.Printf("returning this proc: %+v\n\n", procs[currProcIdx])
			return currProcIdx, true
		}
		//MI VAN HA ÚT KÖZBEN BLOKKOL / TERMINÁL? El kell távolítani a queue -ból azt.
		if currCycleProcIdx < len(cycleProcs)-1 {
			currCycleProcIdx++
			delta := smartTimeQuanum / 2
			if cycleProcs[currCycleProcIdx].ticksLeft <= smartTimeQuanum+delta {
				ticksRemainingForCurrent = cycleProcs[currCycleProcIdx].ticksLeft - 1
				fmt.Printf("new ticks remaining: %+v\n", ticksRemainingForCurrent)
				fmt.Printf("currCycleProcIdx: %+v\n", currCycleProcIdx)
			} else {
				ticksRemainingForCurrent = smartTimeQuanum - 1
			}
			// fmt.Printf("from if: %+v for %d, ", ticksRemainingForCurrent, cycleProcs[currCycleProcIdx].originalIdx)
			fmt.Printf("returning this proc: %+v\n\n", procs[cycleProcs[currCycleProcIdx].originalIdx])
			return cycleProcs[currCycleProcIdx].originalIdx, true
		}

		// fmt.Printf("nope\n")
		fmt.Printf("not returning any from this proc: %+v\n\n", procs)
		return -1, false
	}
	return alg
}

// Only considers ticksLeft of the proc, not the ioOps
func ShortestJobRemaining(procs []Proc, currProcIdx int) (int, bool) {
	var procIdx int

	if isValidRunningProcIdx(procs, currProcIdx) {
		procIdx = currProcIdx
	} else if firstReadyIdx, found := FirstComeFirstServe(procs, currProcIdx); found {
		procIdx = firstReadyIdx
	} else {
		return -1, false
	}

	minTicksLeft := procs[procIdx].TicksLeft

	for i := procIdx + 1; i < len(procs); i++ {
		p := procs[i]
		if p.State == Ready && p.TicksLeft < minTicksLeft {
			minTicksLeft = p.TicksLeft
			procIdx = i
		}
	}
	return procIdx, true
}

func isValidRunningProcIdx(procs []Proc, idx int) bool {
	return 0 <= idx && idx < len(procs) && procs[idx].State == Running
}
