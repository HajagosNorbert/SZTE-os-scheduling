package simulation

import (
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

func MakeSmartRoundRobin() func([]Proc, int) (int, bool) {
	cycleIds := make([]int, 0)
	//index into the cycleIds to get an idx for the procs slice
	currCycleIdx := 0

	ticksForCurr := 0
	newCycle := true
	QUANTUM := 4
	//Smart Time Quantum
	stq := 0

	alg := func(procs []Proc, currProcIdx int) (int, bool) {
		//proceed to the next proc when the current has blocked or terminated
		if currProcIdx != -1 {
			if procs[currProcIdx].State == Terminated || procs[currProcIdx].State == Blocked {
				currCycleIdx++
				if currCycleIdx >= len(cycleIds) {
					newCycle = true
				} else {
					nextProcIdx := cycleIds[currCycleIdx]
					ticksForCurr = getTicksForNextProcInCycle(stq, procs[nextProcIdx].TicksLeft)
				}
			}
		}

		if newCycle {
			cycleIds = cycleIds[:0]
			cycleTickLefts := make([]int, 0)
			for i := 0; i < len(procs); i++ {
				if procs[i].State == Ready || procs[i].State == Running {
					cycleIds = append(cycleIds, i)
					cycleTickLefts = append(cycleTickLefts, procs[i].TicksLeft)
				}
			}
			if len(cycleIds) == 0 {
				return -1, false
			}
			currCycleIdx = 0
			stq = calculateSmartTimeQuantum(cycleTickLefts)
			ticksForCurr = getTicksForNextProcInCycle(stq, procs[cycleIds[0]].TicksLeft)
			ticksForCurr = QUANTUM
			newCycle = false
		}

		nextProcIdx := cycleIds[currCycleIdx]
		ticksForCurr--
		//proceed to the next proc when the alg decided
		//the current proc has ran for long enaugh
		currCycleProcDone := ticksForCurr == 0
		if currCycleProcDone {
			currCycleIdx++
			if currCycleIdx == len(cycleIds) {
				newCycle = true
			} else {
				ticksForCurr = getTicksForNextProcInCycle(stq, procs[cycleIds[currCycleIdx]].TicksLeft)
			}
		}

		return nextProcIdx, true
	}
	return alg
}

func getTicksForNextProcInCycle(smartTimeQuanum int, ticks int) int {
	delta := smartTimeQuanum / 2
	if ticks <= smartTimeQuanum+delta {
		return ticks
	} else {
		return smartTimeQuanum
	}

}

// ticksLeft is > 0
func calculateSmartTimeQuantum(ticksLeft []int) int {
	smartTimeQuanum := 0
	if len(ticksLeft) > 1 {
		sort.Ints(ticksLeft)
		ticksLeftSum := 0
		for i := 0; i < len(ticksLeft)-1; i++ {
			ticksLeftSum += ticksLeft[i+1] - ticksLeft[i]
		}
		smartTimeQuanum = int(math.Round(float64(ticksLeftSum) / (float64(len(ticksLeft)) - 1.0)))
	}

	if smartTimeQuanum == 0 {
		assignedTicks := int(ticksLeft[0] / 4)
		if assignedTicks < 1 {
			smartTimeQuanum = ticksLeft[0]
		} else {
			smartTimeQuanum = assignedTicks
		}
	}

	return smartTimeQuanum
}

func MakeSmartRoundRobinOld() func([]Proc, int) (int, bool) {
	ticksRemainingForCurrent := 0
	currCycleProcIdx := -1
	smartTimeQuanum := 0
	var cycleProcs []procInCycle
	insideCycle := false
	alg := func(procs []Proc, currProcIdx int) (int, bool) {

		if currProcIdx != -1 {
			if procs[currProcIdx].State == Blocked || procs[currProcIdx].State == Terminated {
				ticksRemainingForCurrent = 0
				if currCycleProcIdx < len(cycleProcs)-1 {
					cycleProcs[currCycleProcIdx] = cycleProcs[len(cycleProcs)-1]
					cycleProcs = cycleProcs[0 : len(cycleProcs)-1]
					currCycleProcIdx--
				}
			}
		}
		endOfCycle := insideCycle && ticksRemainingForCurrent == 0 && currCycleProcIdx <= len(cycleProcs)-1
		if endOfCycle {
			insideCycle = false
		}
		if !insideCycle {
			//one new cycle
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
			sort.Slice(cycleProcs, func(i, j int) bool {
				return cycleProcs[i].ticksLeft < cycleProcs[j].ticksLeft
			})

			if len(cycleProcs) == 0 {
				return -1, false
			}

			insideCycle = true
			if len(cycleProcs) == 1 {
				currCycleProcIdx = 0
				ticksRemainingForCurrent = cycleProcs[currCycleProcIdx].ticksLeft
				return cycleProcs[currCycleProcIdx].originalIdx, true
			} else {
				currCycleProcIdx = -1
			}

			ticksLeftSum := 0
			for i := 0; i < len(cycleProcs)-1; i++ {
				ticksLeftSum += cycleProcs[i+1].ticksLeft - cycleProcs[i].ticksLeft
			}

			smartTimeQuanum = int(math.Round(float64(ticksLeftSum) / (float64(len(procs)) - 1.0)))
			if smartTimeQuanum == 0 {
				assignedTicks := int(cycleProcs[0].ticksLeft / 8)
				if assignedTicks < 1 {
					smartTimeQuanum = cycleProcs[0].ticksLeft
				} else {
					smartTimeQuanum = assignedTicks
				}
			}
		}
		if ticksRemainingForCurrent > 0 {
			ticksRemainingForCurrent--
			return currProcIdx, true
		}
		if currCycleProcIdx < len(cycleProcs)-1 {
			currCycleProcIdx++
			delta := smartTimeQuanum / 2
			if cycleProcs[currCycleProcIdx].ticksLeft <= smartTimeQuanum+delta {
				ticksRemainingForCurrent = cycleProcs[currCycleProcIdx].ticksLeft
			} else {
				ticksRemainingForCurrent = smartTimeQuanum
			}
			return cycleProcs[currCycleProcIdx].originalIdx, true
		}

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
