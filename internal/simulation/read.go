package simulation

//TODO: fn mondja meg az input proc-ok IO-ja milyen valószínűséggel jelen meg adott időben
func ReadProcs() []Proc {
	firstIoOp := []IoOp{{startsAfter: 1, ticksLeft: 2}}
	secIoOp := []IoOp{{startsAfter: 2, ticksLeft: 2}}

	procs := []Proc{
		{ticksLeft: 3, ioOps: firstIoOp},
		{ticksLeft: 3, spawnedAt: 1, ioOps: secIoOp},
	}

	return procs
}
