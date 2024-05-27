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

func OnStats(f func(stats SystemStats)) (events.ID, error) {
	return _service.eventAll.RegisterHandler(events.HandlerFunc[SystemStats](f))
}
func OnCpuStats(f func(stats CpuStats)) (events.ID, error) {
	return _service.eventCpu.RegisterHandler(events.HandlerFunc[CpuStats](f))
}
func OnMemoryStats(f func(stats MemoryStats)) (events.ID, error) {
	return _service.eventMemory.RegisterHandler(events.HandlerFunc[MemoryStats](f))
}
func OnNvidiaStats(f func(stats NvidiaStats)) (events.ID, error) {
	return _service.eventNvidia.RegisterHandler(events.HandlerFunc[NvidiaStats](f))
}
