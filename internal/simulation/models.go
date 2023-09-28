package simulation

/*
An IoOp's ticksLeft prop will first be decremented when the tick counter is at startsAfter + 1
Could say that it takes up 1 tick to start the ioOp which is part of the totalTicks
*/
type IoOp struct {
	//zero based
	startsAfter int
	ticksLeft   int
}

//There needs to be at least as meany ticksLeft as len(ioOps), since starting an ioOp takes 1 tick.
//ticksLeft and totalTicks should have the same value when instantiated
type Proc struct {
	//zero based
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
	Blocked
	Terminated
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
