package simulation

func FirstComeFirstServe(procs []Proc, currProcIdx int) (int, bool) {

	if isValidRunningProcIdx(procs, currProcIdx) {
		return currProcIdx, true
	}
	for i, p := range procs {
		if p.state == Ready {
			return i, true
		}
	}
	return -1, false
}

// Only considers ticksLeft of the proc, not the ioOps
func ShortestJobRemaining(procs []Proc, currProcIdx int) (int, bool) {
	invalidProcIdx := -1
	var procIdx int

	if isValidRunningProcIdx(procs, currProcIdx) {
		procIdx = currProcIdx
	} else if firstReadyIdx, found := FirstComeFirstServe(procs, invalidProcIdx); found {
		procIdx = firstReadyIdx
	} else {
		return -1, false
	}

	minTicksLeft := procs[procIdx].ticksLeft

	for i := procIdx + 1; i < len(procs); i++ {
		p := procs[i]
		if p.state == Ready && p.ticksLeft < minTicksLeft {
			minTicksLeft = p.ticksLeft
			procIdx = i
		}
	}
	return procIdx, true
}

func isValidRunningProcIdx(procs []Proc, idx int) bool {
	return 0 <= idx && idx < len(procs) && procs[idx].state == Running
}
