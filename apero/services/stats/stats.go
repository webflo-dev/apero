package stats

import (
	"sync"
	"time"
)

type Stats struct {
	Cpu    *CpuStats
	Memory *MemoryStats
	Nvidia *NvidiaStats
}

func runCpu(stats *Stats, done func()) {
	go func() {
		if cpu, _ := GetCpuStats(); cpu != nil {
			stats.Cpu = cpu
		}
		done()
	}()
}

func runMemory(stats *Stats, done func()) {
	go func() {
		if memory, _ := GetMemoryStats(); memory != nil {
			stats.Memory = memory
		}
		done()
	}()
}

func runNvidia(stats *Stats, done func()) {
	go func() {
		if nvidia, _ := GetNvidiaStats(); nvidia != nil {
			stats.Nvidia = nvidia
		}
		done()
	}()
}

var statsChan chan *Stats

func WatchStats() <-chan *Stats {
	if statsChan != nil {
		return statsChan
	}

	statsChan := make(chan *Stats)

	go func() {
		stats := &Stats{}

		var wg sync.WaitGroup
		for range time.Tick(time.Second) {

			done := func() {
				wg.Done()
			}

			wg.Add(3)
			runCpu(stats, done)
			runMemory(stats, done)
			runNvidia(stats, done)

			wg.Wait()
			statsChan <- stats
		}
	}()

	return statsChan
}
