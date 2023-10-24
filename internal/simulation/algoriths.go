package simulation

import "math/rand"

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
func RoundRobin(procs []Proc, currProcIdx int) (int, bool) {
	var procIdx int

	if isValidRunningProcIdx(procs, currProcIdx) {
		procIdx = currProcIdx
	} else {
		return FirstComeFirstServe(procs, currProcIdx)
	}

	for i := 1; i < len(procs); i++ {
		procIdxCandidate := (procIdx + i) % len(procs)
		if procs[procIdxCandidate].State == Ready {
			return procIdxCandidate, true
		}
	}
	return procIdx, true
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
