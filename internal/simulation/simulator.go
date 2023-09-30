package simulation

import (
	"sort"
)

func SimulateScheduling(procs []Proc, SchedAlg func([]Proc, int) (int, bool)) SimResult {
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
			return endSimulation(procs, result, tick)
		}
		procs = readyUpProcs(procs, tick)

		// pick next ready proc to run
		choosenProcIdx, found := SchedAlg(procs, procIdx)
		if !found {
			tick++
			tickForIoOps(ioTasksRunning)
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

		tickForWaitingProcs(procs, &result)
		tickForIoOps(ioTasksRunning)
		proc, ioTasksRunning = blockProcOnIo(proc, ioTasksRunning)
		proc.ticksLeft--

		procToBeTerminated := proc.ticksLeft <= 0 && proc.state == Running
		if procToBeTerminated {
			proc.state = Terminated
		}

		tick++
	}
}

func blockProcOnIo(proc *Proc, ioTasksRunning []IoTaskRunning) (*Proc, []IoTaskRunning) {
	ioOpToStart, isIoOpReady := proc.getReadyIoOp()
	taskToStart := IoTaskRunning{ioOp: ioOpToStart, ownerProc: proc}
	if isIoOpReady {
		ioTasksRunning = append(ioTasksRunning, taskToStart)
		proc.state = Blocked
	}
	return proc, ioTasksRunning
}

func tickForWaitingProcs(procs []Proc, result *SimResult) {
	for i, p := range procs {
		if p.state == Ready {
			result.procResults[i].idleTicks++
		}
	}
}

func endSimulation(procs []Proc, result SimResult, tick int) SimResult {
	result.totalTicks = tick
	for i := range result.procResults {
		result.procResults[i].totalTicks = procs[i].ticksDone()
	}
	return result
}

func tickForIoOps(ioTasks []IoTaskRunning) {
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
		return tasks
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
