package simulation

import (
	"sort"
)

func FirstComeFirstServe(procs []Proc) SimResult {
	tick := 0
	nextNewProc := 0
	var proc Proc

	sort.Slice(procs, func(i, j int) bool {
		return procs[i].startAt < procs[j].startAt
	})

	for {
		// change state where needed New -> Ready
		for i, p := range procs[nextNewProc:] {
			if p.state == New && p.startAt <= tick {
				p.state = Ready
			} else {
				nextNewProc = i
				break
			}
		}

		// pick next proc to operate on
		idx, found := firstReady(procs)
		if !found {
			break
		}
		proc = procs[idx]
		proc.state = Running

        //way down lower
        tick++;
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
