package simulation

/*
An IoOp's ticksLeft prop will first be decremented when the tick counter is at startsAfter + 1
Could say that it takes up 1 tick to start the ioOp which is part of the proc's totalTicks
*/
type IoOp struct {
	//zero based, relative to the proc's ticks, not to the global ticks
	StartsAfter int
	TicksLeft   int
}

// There needs to be at least as meany ticksLeft as len(ioOps), since starting an ioOp takes 1 tick.
// ticksLeft and totalTicks should have the same value when instantiated
type Proc struct {
	//zero based
	SpawnedAt  int
	TicksLeft  int
	TotalTicks int
	Priority   int
	UserId     int
	State      ProcState
	IoOps      []IoOp
}

type IoTaskRunning struct {
	ioOp      *IoOp
	ownerProc *Proc
}

type IoOpState int
type ProcState = string

const (
	New        ProcState = "new"
	Ready      ProcState = "ready"
	Running    ProcState = "running"
	Blocked    ProcState = "blocked"
	Terminated ProcState = "terminated"
)

type Algorithm = string

const (
	AlgFcfs       Algorithm = "fcfs"
	AlgSjr        Algorithm = "sjr"
	AlgLottery    Algorithm = "lottery"
	AlgSrr        Algorithm = "srr"
	AlgRoundRobin Algorithm = "rr"
)

type SimResult struct {
	algName     string
	idleTicks   int
	totalTicks  int
	procResults []ProcResult
}

type ProcResult struct {
	ctxSwitchCount int
	idleTicks      int
	totalTicks     int
}
