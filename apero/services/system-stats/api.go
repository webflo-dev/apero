package systemStats

import (
	"webflo-dev/apero/events"
)

type SystemStats struct {
	Cpu    *CpuStats
	Memory *MemoryStats
	Nvidia *NvidiaStats
}

var _service = newService()

func StartService() {
	_service.start()
}

func OnStats(id string, f func(stats SystemStats)) {
	_service.eventAll.RegisterHandler(id, events.HandlerFunc[SystemStats](f))
}
func OnCpuStats(id string, f func(stats CpuStats)) {
	_service.eventCpu.RegisterHandler(id, events.HandlerFunc[CpuStats](f))
}
func OnMemoryStats(id string, f func(stats MemoryStats)) {
	_service.eventMemory.RegisterHandler(id, events.HandlerFunc[MemoryStats](f))
}
func OnNvidiaStats(id string, f func(stats NvidiaStats)) {
	_service.eventNvidia.RegisterHandler(id, events.HandlerFunc[NvidiaStats](f))
}
