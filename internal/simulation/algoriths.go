package simulation

import (
	"fmt"
	"sort"
)

func FirstComeFirstServe(procs []Proc) SimResult {
	procIdx := -1
	var proc *Proc
	var ioTasksRunning []IoTaskRunning
	tick := 0
	procResults := make([]ProcResult, len(procs))
	result := SimResult{procResults: procResults}
	readyUpProcs := readyUpProcsFactory()

	sort.Slice(procs, func(i, j int) bool {
		return procs[i].spawnedAt < procs[j].spawnedAt
	})

	for {
		if allProcsTerminated(procs) {
			result.totalTicks = tick
			for i := range procResults {
				procResults[i].totalTicks = procs[i].ticksDone()
			}
			return result
		}
		fmt.Printf("")
		// fmt.Printf("%+v\n", procs)
		procs = readyUpProcs(procs, tick)

		// pick next ready proc to run
		choosenProcIdx, found := firstReady(procs, procIdx)
		if !found {
			tick++
			tickIoOps(ioTasksRunning)
			result.idleTicks++
			continue
		}

		contextSwitchHappened := choosenProcIdx != procIdx && proc != nil && proc.state == Running
		if contextSwitchHappened {
			println("this")
			proc.state = Ready
			result.procResults[procIdx].ctxSwitchCount++
		}
		procIdx = choosenProcIdx
		proc = &procs[choosenProcIdx]
		proc.state = Running

		// increment idleTicks on Ready procs
		for i, p := range procs {
			if p.state == Ready {
				result.procResults[i].idleTicks++
			}
		}

		tickIoOps(ioTasksRunning)
		//check for blocking
		ioOpToStart, isIoOpReady := proc.getReadyIoOp()
		taskToStart := IoTaskRunning{ioOp: ioOpToStart, ownerProc: proc}
		if isIoOpReady {
			ioTasksRunning = append(ioTasksRunning, taskToStart)
			proc.state = Blocked
		}
		proc.ticksLeft--

		// check if terminated
		if proc.ticksLeft <= 0 && proc.state == Running {
			proc.state = Terminated
		}

		tick++
	}
}

func tickIoOps(ioTasks []IoTaskRunning) {
	for i := range ioTasks {
		task := &ioTasks[i]
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
func readyUpProcsFactory() func([]Proc, int) []Proc {
	nextNewProc := 0

	return func(procs []Proc, tick int) []Proc {
		// change state where needed New -> Ready
		for i := range procs[nextNewProc:] {
			p := &procs[nextNewProc+i]
			if p.state == New && p.spawnedAt <= tick {
				p.state = Ready
			} else {
				// there will be no new proc which spawned earlyer, since they are sorted by spawnedAt
				// we save where we stopped
				nextNewProc = i
				break
			}
		}
		return procs
	}
}
func removeTaskAt(tasks []IoTaskRunning, i int) []IoTaskRunning {
	if len(tasks) == 0 {

	}
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

func (p *Proc) getReadyIoOp() (*IoOp, bool) {
	for i := range p.ioOps {
		ioOp := &p.ioOps[i]
		if ioOp.ticksLeft > 0 && ioOp.startsAfter <= p.ticksDone() {
			return ioOp, true
		}
	}
	return &IoOp{}, false
}

func firstReady(procs []Proc, currProcIdx int) (int, bool) {
	isValidRunningProcIdx := 0 <= currProcIdx && currProcIdx < len(procs) && procs[currProcIdx].state == Running

	if isValidRunningProcIdx {
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
func SortestJobFirst(procs []Proc) SimResult {
	procIdx := -1
	var proc *Proc
	var ioTasksRunning []IoTaskRunning
	tick := 0
	procResults := make([]ProcResult, len(procs))
	result := SimResult{procResults: procResults}
	readyUpProcs := readyUpProcsFactory()

	sort.Slice(procs, func(i, j int) bool {
		return procs[i].spawnedAt < procs[j].spawnedAt
	})

	for {
		if allProcsTerminated(procs) {
			result.totalTicks = tick
			for i := range procResults {
				procResults[i].totalTicks = procs[i].ticksDone()
			}
			return result
		}

		procs = readyUpProcs(procs, tick)

		// pick next ready proc to run
		choosenProcIdx, found := shortestReadyProc(procs, procIdx)
		if !found {
			tick++
			tickIoOps(ioTasksRunning)
			result.idleTicks++
			continue
		}

		contextSwitchHappened := choosenProcIdx != procIdx && proc != nil && proc.state == Running
		if contextSwitchHappened {
			proc.state = Ready
			result.procResults[procIdx].ctxSwitchCount++
		}
		procIdx = choosenProcIdx
		proc = &procs[choosenProcIdx]
		proc.state = Running

		// increment idleTicks on Ready procs
		for i, p := range procs {
			if p.state == Ready {
				result.procResults[i].idleTicks++
			}
		}

		tickIoOps(ioTasksRunning)
		//check for blocking
		ioOpToStart, isIoOpReady := proc.getReadyIoOp()
		taskToStart := IoTaskRunning{ioOp: ioOpToStart, ownerProc: proc}
		if isIoOpReady {
			ioTasksRunning = append(ioTasksRunning, taskToStart)
			proc.state = Blocked
		}
		proc.ticksLeft--

		// check if terminated
		if proc.ticksLeft <= 0 && proc.state == Running {
			proc.state = Terminated
		}

		tick++
	}
}

func shortestReadyProc(procs []Proc, currProcIdx int) (int, bool) {
	invalidProcIdx := -1
	var procIdx int
	isValidRunningProcIdx := 0 <= currProcIdx && currProcIdx < len(procs) && procs[currProcIdx].state == Running

	if isValidRunningProcIdx {
		procIdx = currProcIdx
	} else if firstReadyIdx, found := firstReady(procs, invalidProcIdx); found {
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
