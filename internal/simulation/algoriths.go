package simulation

import (
	"sort"
)

func FirstComeFirstServe(procs []Proc) SimResult {
	tick := 0
	nextNewProc := 0
	var ioTasksRunning []IoTaskRunning
	var proc Proc
	procResults := make([]ProcResult, len(procs))
	result := SimResult{procResults: procResults}

	sort.Slice(procs, func(i, j int) bool {
		return procs[i].spawnedAt < procs[j].spawnedAt
	})

	for {
		if allProcsTerminated(procs) {
			result.totalTicks = tick
			for i, pRes := range procResults{
				pRes.totalTicks = procs[i].ticksDone()
				pRes.idleTicks = procs[i].
			}
			return result
		}

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

		// pick next ready proc to run
		idx, found := firstReady(procs)
		if !found {
			tick++
			tickIoOps(ioTasksRunning)
			result.idleTicks++
			continue //??
		}
		proc = procs[idx]
		proc.state = Running

		// increment idleTicks on Ready procs
		for i, p := range procs {
			if p.state == Ready {
				result.procResults[i].idleTicks++ 
			}
		}

		//check for blocking
		ioOpToStart, isIoOpReady := proc.getReadyIoOp()
		if isIoOpReady {
			ioTasksRunning = append(ioTasksRunning, IoTaskRunning{ioOp: &ioOpToStart, ownerProc: &proc})
			proc.state = Blocked
		}
		// do the work (CPU or initiating IO)
		proc.ticksLeft--

		// check if terminated
		if proc.ticksLeft == 0 && proc.state == Running {
			proc.state = Terminated
		}

		//at the end
		tick++
		tickIoOps(ioTasksRunning)
	}
}

func tickIoOps(ioTasks []IoTaskRunning) {
	for i, task := range ioTasks {
		task.ioOp.ticksLeft--
		if task.ioOp.ticksLeft == 0 {
			if task.ownerProc.ticksLeft > 0 {
				task.ownerProc.state = Ready
			} else {
				task.ownerProc.state = Terminated
			}
			ioTasks = removeTaskAt(ioTasks, i)
		}
	}
}
func removeTaskAt(tasks []IoTaskRunning, i int) []IoTaskRunning {
	tasks[i] = tasks[len(tasks)-1]
	return tasks[:len(tasks)-1]
}

func allProcsTerminated(procs []Proc) bool {
	for _, p := range procs {
		if p.state != Terminated {
			return false
		}
	}
	return true
}

func (p Proc) ticksDone() int {
	return p.totalTicks - p.ticksLeft
}

func (p *Proc) getReadyIoOp() (IoOp, bool) {
	for _, ioOp := range p.ioOps {
		if ioOp.ticksLeft > 0 && ioOp.startsAfter >= p.ticksDone() {
			return ioOp, true
		}
	}
	return IoOp{}, false
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
