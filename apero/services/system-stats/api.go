package systemStats

import (
	"log"
	"reflect"
	"strings"
	"sync"
	"time"
)

type EventType string

const (
	EventAll    EventType = "updateall"
	EventCpu    EventType = "updatecpu"
	EventMemory EventType = "updatememory"
	EventNvidia EventType = "updatenvidia"
)

type SystemStatsEventHandler interface {
	UpdateAll(stats *SystemStats)
	UpdateCpu(cpu *CpuStats)
	UpdateMemory(memory *MemoryStats)
	UpdateNvidia(nvidia *NvidiaStats)
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

func StartService() {
	service.listen()
}

func (service *systemStatsService) listen() {
	if service.listening {
		return
	}

	service.listening = true

	log.Printf("methods? %+v\n", eventMethods)

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

			if eventMethod, found := eventMethods[EventAll]; found {
				eventMethod.call(EventAll, []any{stats})
			}
		}
	}()
}

// / Registering events
type valueConvertor = func(value any) reflect.Value

type eventMethod struct {
	name   string
	values []valueConvertor
}

var eventMethods = newEventMethods()

func newEventMethods() map[EventType]*eventMethod {
	iface := reflect.TypeOf(struct{ SystemStatsEventHandler }{})
	lenMethods := iface.NumMethod()

	eventMethods := make(map[EventType]*eventMethod, lenMethods)

	for i := 0; i < lenMethods; i++ {
		method := iface.Method(i)
		lenParams := method.Type.NumIn()

		eventMethod := &eventMethod{
			name:   method.Name,
			values: make([]valueConvertor, lenParams-1),
		}

		for j := 1; j < lenParams; j++ {
			t := method.Type.In(j)
			switch t.Kind() {
			// case reflect.Int:
			// 	eventMethod.values[j-1] = toInt
			// 	break
			// case reflect.Bool:
			// 	eventMethod.values[j-1] = toBool
			// 	break
			default:
				eventMethod.values[j-1] = func(value any) reflect.Value {
					return reflect.ValueOf(value)
				}
				break
			}
		}
		eventMethods[EventType(strings.ToLower(method.Name))] = eventMethod
	}
	return eventMethods
}

func (m *eventMethod) call(eventType EventType, values []any) {
	in := make([]reflect.Value, len(m.values)+1)
	for i, value := range values {
		in[i+1] = m.values[i](value)
	}

	for _, subscriber := range eventSubscribers[eventType] {
		in[0] = reflect.ValueOf(subscriber.handle)
		subscriber.callback.Call(in)
	}
}

type subscriber struct {
	callback reflect.Value
	handle   any
}

type subscribers map[EventType][]subscriber

var eventSubscribers = make(subscribers)

func RegisterForEvents(handler any) {
	objType := reflect.TypeOf(handler)
	for i := 0; i < objType.NumMethod(); i++ {
		method := objType.Method(i)
		eventType := EventType(strings.ToLower(method.Name))
		eventSubscribers[eventType] = append(eventSubscribers[eventType], subscriber{method.Func, handler})
	}
}
