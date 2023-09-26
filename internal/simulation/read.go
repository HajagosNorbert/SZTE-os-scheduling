package simulation

//TODO: fn mondja meg az input proc-ok IO-ja milyen valószínűséggel jelen meg adott időben
func ReadProcs() []Proc {
	firstIoOp := []IoOp{{startsAt: 1, ticksLeft: 2}}
	secIoOp := []IoOp{{startsAt: 2, ticksLeft: 2}}

	procs := []Proc{
		{ticksLeft: 3, ioOps: firstIoOp},
		{ticksLeft: 3, spawnedAt: 1, ioOps: secIoOp},
	}

	return procs
}
