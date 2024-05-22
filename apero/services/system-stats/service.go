package systemStats

import (
	"sync"
	"time"
	"webflo-dev/apero/events"
	"webflo-dev/apero/services"
)

type service struct {
	base services.Service

	stats *SystemStats

	eventAll    events.Event[SystemStats]
	eventCpu    events.Event[CpuStats]
	eventMemory events.Event[MemoryStats]
	eventNvidia events.Event[NvidiaStats]
}

func newService() *service {
	service := &service{
		base:        services.NewService(),
		stats:       &SystemStats{},
		eventAll:    events.New[SystemStats](),
		eventCpu:    events.New[CpuStats](),
		eventMemory: events.New[MemoryStats](),
		eventNvidia: events.New[NvidiaStats](),
	}

	return service
}

func (s *service) start() {
	s.base.Start(nil, s.loop)
}

func (s *service) loop() services.LoopBehavior {

	var wg sync.WaitGroup
	for range time.Tick(time.Second) {
		done := func() {
			wg.Done()
		}

		wg.Add(3)

		s.runCpu(done)
		s.runMemory(done)
		s.runNvidia(done)

		wg.Wait()

		s.eventAll.Publish(*s.stats)
	}

	return services.LoopBehaviorStop
}

func (s *service) runCpu(done func()) {
	go func() {
		if cpu, _ := GetCpuStats(); cpu != nil {
			s.stats.Cpu = cpu
		}

		s.eventCpu.Publish(*s.stats.Cpu)

		done()
	}()
}

func (s *service) runMemory(done func()) {
	go func() {
		if memory, _ := GetMemoryStats(); memory != nil {
			s.stats.Memory = memory
		}

		s.eventMemory.Publish(*s.stats.Memory)

		done()
	}()
}

func (s *service) runNvidia(done func()) {
	go func() {
		if nvidia, _ := GetNvidiaStats(); nvidia != nil {
			s.stats.Nvidia = nvidia
		}

		s.eventNvidia.Publish(*s.stats.Nvidia)

		done()
	}()
}
