package simulation

// TODO: fn mondja meg az input proc-ok IO-ja milyen valószínűséggel jelen meg adott időben
func ReadProcs() []Proc {
	// firstIoOp := []IoOp{{startsAfter: 2, ticksLeft: 3}}
	// secIoOp := []IoOp{{startsAfter: 2, ticksLeft: 2}}
	//
	totalTicks1 := 9
	totalTicks2 := 5
	procs := []Proc{
		{TicksLeft: totalTicks1, TotalTicks: totalTicks1},
		{TicksLeft: totalTicks2, TotalTicks: totalTicks2, SpawnedAt: 1},
	}
	// firstIoOp := []IoOp{{startsAfter: 2, ticksLeft: 3}}
	// secIoOp := []IoOp{{startsAfter: 2, ticksLeft: 2}}
	//
	// totalTicks := 3
	// procs := []Proc{
	// 	{ticksLeft: totalTicks, totalTicks: totalTicks, ioOps: firstIoOp},
	// 	{ticksLeft: totalTicks, totalTicks: totalTicks, spawnedAt: 1, ioOps: secIoOp},
	// }
	return procs
}


