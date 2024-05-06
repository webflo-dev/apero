package systemStats

import (
	"sync"
	"time"
)

type EventType string

const (
	EventAll    EventType = "UpdateAll"
	EventCpu    EventType = "UpdateCpu"
	EventMemory EventType = "UpdateMemory"
	EventNvidia EventType = "UpdateNvidia"
)

type SystemStats struct {
	Cpu    *CpuStats
	Memory *MemoryStats
	Nvidia *NvidiaStats
}

type Subscriber interface {
	UpdateAll(stats *SystemStats)
	UpdateCpu(cpu *CpuStats)
	UpdateMemory(memory *MemoryStats)
	UpdateNvidia(nvidia *NvidiaStats)
}

type service struct {
	started     bool
	subscribers map[EventType][]Subscriber
	stats       *SystemStats
}

var _service = newService()

func newService() *service {
	service := &service{
		started:     false,
		subscribers: make(map[EventType][]Subscriber),
		stats:       &SystemStats{},
	}
	return service
}

func StartService() {
	_service.start()
}

func (s *service) stop() {
	s.started = false
}

func (s *service) start() {
	if s.started {
		return
	}

	s.started = true

	go func() {
		defer s.stop()

		var wg sync.WaitGroup
		for range time.Tick(time.Second) {
			if s.started == false {
				return
			}

			done := func() {
				wg.Done()
			}

			wg.Add(3)

			s.runCpu(done)
			s.runMemory(done)
			s.runNvidia(done)

			wg.Wait()

			for _, subscriber := range s.subscribers[EventAll] {
				subscriber.UpdateAll(s.stats)
			}
		}
	}()
}

func Register[T Subscriber](handle T, events ...EventType) {
	for _, event := range events {
		_service.subscribers[event] = append(_service.subscribers[event], handle)
	}
}

func (s *service) runCpu(done func()) {
	go func() {
		if cpu, _ := GetCpuStats(); cpu != nil {
			s.stats.Cpu = cpu
		}

		for _, subscriber := range s.subscribers[EventCpu] {
			subscriber.UpdateCpu(s.stats.Cpu)
		}

		done()
	}()
}

func (s *service) runMemory(done func()) {
	go func() {
		if memory, _ := GetMemoryStats(); memory != nil {
			s.stats.Memory = memory
		}

		for _, subscriber := range s.subscribers[EventMemory] {
			subscriber.UpdateMemory(s.stats.Memory)
		}

		done()
	}()
}

func (s *service) runNvidia(done func()) {
	go func() {
		if nvidia, _ := GetNvidiaStats(); nvidia != nil {
			s.stats.Nvidia = nvidia
		}

		for _, subscriber := range s.subscribers[EventNvidia] {
			subscriber.UpdateNvidia(s.stats.Nvidia)
		}

		done()
	}()
}
