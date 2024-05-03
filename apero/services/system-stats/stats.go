package systemStats

type SystemStats struct {
	Cpu    *CpuStats
	Memory *MemoryStats
	Nvidia *NvidiaStats
}

func runCpu(stats *SystemStats, done func()) {
	go func() {
		if cpu, _ := GetCpuStats(); cpu != nil {
			stats.Cpu = cpu
		}

		if eventMethod, found := eventMethods[EventCpu]; found {
			eventMethod.call(EventCpu, []any{stats.Cpu})
		}

		done()
	}()
}

func runMemory(stats *SystemStats, done func()) {
	go func() {
		if memory, _ := GetMemoryStats(); memory != nil {
			stats.Memory = memory
		}

		if eventMethod, found := eventMethods[EventMemory]; found {
			eventMethod.call(EventMemory, []any{stats.Memory})
		}

		done()
	}()
}

func runNvidia(stats *SystemStats, done func()) {
	go func() {
		if nvidia, _ := GetNvidiaStats(); nvidia != nil {
			stats.Nvidia = nvidia
		}

		if eventMethod, found := eventMethods[EventNvidia]; found {
			eventMethod.call(EventNvidia, []any{stats.Nvidia})
		}

		done()
	}()
}
