package hyprland

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

func (service *hyprlandEventService) listen() {
	if service.listening {
		return
	}

	go func() {
		service.createEventsConnection()
		defer service.closeConnection()

		logger.Println("listening for hyprland events")

		for {
			msg, err := service.readEvent()
			if err != nil {
				logger.Println("Error receiving message", err)
				return
			}

			for _, data := range msg {
				log.Printf("%+v\n", data)
				service.processEvent(data)
			}
		}
	}()
}

func Dispatch(pattern string, args ...any) error {
	cmd := fmt.Sprintf(pattern, args...)
	dispatchCmd := fmt.Sprintf("dispatch %s", cmd)
	return writeCmd(dispatchCmd, nil)
}

func Workspaces() []Workspace {
	return hyprlCtl.workspaces
}
func ActiveWorkspace() Workspace {
	return hyprlCtl.activeWorkspace
}

func Clients() []Client {
	return hyprlCtl.clients
}

func ActiveClient() Client {
	return hyprlCtl.activeClient
}

type subscriber struct {
	callback reflect.Value
	handle   any
}

type subscribers map[EventType][]subscriber

var eventSubscribers = make(subscribers)

func RegisterForEvents(handler any) {
	if eventService.listening == false {
		eventService.listen()
	}

	objType := reflect.TypeOf(handler)
	for i := 0; i < objType.NumMethod(); i++ {
		method := objType.Method(i)
		eventType := EventType(strings.ToLower(method.Name))
		eventSubscribers[eventType] = append(eventSubscribers[eventType], subscriber{method.Func, handler})
	}
}
