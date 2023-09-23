package simulation

type IoOp struct {
	startsAt  int
	ticksLeft int
}

type Proc struct {
	startAt   int
	ticksLeft int
	state     ProcState
	ioOps     []IoOp
}

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
	interruptCount int
	idleTicks      int
	totalTicks     int
}
