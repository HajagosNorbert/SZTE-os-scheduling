package simulation

func ReadProcs() []Proc {
	firstIoOp := []IoOp{{startsAt: 1, ticksLeft: 2}}
	secIoOp := []IoOp{{startsAt: 2, ticksLeft: 2}}

	procs := []Proc{
		{ticksLeft: 3, ioOps: firstIoOp},
		{ticksLeft: 3, startAt: 1, ioOps: secIoOp},
	}

	return procs
}
