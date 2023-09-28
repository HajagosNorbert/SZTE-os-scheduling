package simulation

/*
An IoOp's ticksLeft prop will first be decremented when the tick counter is at startsAfter + 1
Could say that it costs 1 tick to start the ioOp
*/
type IoOp struct {
	startsAfter int
	ticksLeft   int
}

//There needs to be at least as meany ticksLeft as len(ioOps), since starting an ioOp takes 1 tick
type Proc struct {
	spawnedAt  int
	ticksLeft  int
	totalTicks int
	state      ProcState
	ioOps      []IoOp
}

type IoTaskRunning struct {
	ioOp      *IoOp
	ownerProc *Proc
}

type IoOpState int
type ProcState int

const (
	New ProcState = iota
	Ready
	Running
	Terminated
	Blocked
)

type SimResult struct {
	idleTicks   int
	totalTicks  int
	procResults []ProcResult
}

type ProcResult struct {
	ctxSwitchCount int
	idleTicks      int
	totalTicks     int
}
