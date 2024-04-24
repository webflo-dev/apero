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
		done()
	}()
}

func runMemory(stats *SystemStats, done func()) {
	go func() {
		if memory, _ := GetMemoryStats(); memory != nil {
			stats.Memory = memory
		}
		done()
	}()
}

func runNvidia(stats *SystemStats, done func()) {
	go func() {
		if nvidia, _ := GetNvidiaStats(); nvidia != nil {
			stats.Nvidia = nvidia
		}
		done()
	}()
}
