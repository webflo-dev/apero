package systemStats

import (
	"sync"
	"time"
)

type SystemStatsEventHandler interface {
	Notify(stats *SystemStats)
}

type systemStatsService struct {
	listening   bool
	subscribers []SystemStatsEventHandler
}

var service = newSystemStatsService()

func newSystemStatsService() *systemStatsService {
	service := &systemStatsService{
		listening: false,
	}
	return service
}

func (service *systemStatsService) listen() {
	if service.listening {
		return
	}

	go func() {
		stats := &SystemStats{}

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

			for _, subscriber := range service.subscribers {
				subscriber.Notify(stats)
			}
		}
	}()
}

func WatchSystemStats(handler SystemStatsEventHandler) {
	if service.listening == false {
		service.listen()
	}

	service.subscribers = append(service.subscribers, handler)
}
