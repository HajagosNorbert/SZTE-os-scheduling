package simulation

import (
	"sort"
)

func FirstComeFirstServe(procs []Proc) SimResult {
	tick := 0
	nextNewProc := 0
	var proc Proc

	sort.Slice(procs, func(i, j int) bool {
		return procs[i].spawnedAt < procs[j].spawnedAt
	})

	for {
		// change state where needed New -> Ready
		for i, p := range procs[nextNewProc:] {
			if p.state == New && p.spawnedAt <= tick {
				p.state = Ready
			} else {
				// there will be no new proc which spawned earlyer, since they are sorted by spawnedAt
				// we save where we stopped
				nextNewProc = i
				break
			}
		}

		// pick next proc to operate on
		idx, found := firstReady(procs)
		if !found {
			tick++
			break
		}
		proc = procs[idx]
		proc.state = Running
		// do the work

		// check if finished
		//check for blocking


		//at the end
		tick++
	}

	procResults := make([]ProcResult, len(procs))
	res := SimResult{procResults: procResults}
	return res
}

func firstReady(procs []Proc) (int, bool) {
	for i, p := range procs {
		if p.state == Ready {
			return i, true
		}
	}
	return -1, false
}

func SortestJobFirst(procs []Proc) SimResult {
	panic("algorithm not yet implemented!")
}
