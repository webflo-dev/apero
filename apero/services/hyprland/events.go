package hyprland

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type hyprlandEventService struct {
	hyprlandIpcService
	listening   bool
	subscribers map[EventType][]any
}

var eventService = newHyprlandEventService()

func newHyprlandEventService() *hyprlandEventService {
	service := &hyprlandEventService{
		listening:   false,
		subscribers: make(map[EventType][]any),
	}

	return service
}

func preprocessClientAddress(values []string) []string {
	values[0] = fmt.Sprintf("0x%s", values[0])
	return values
}

func (s *hyprlandEventService) processEvent(msg EventData) {
	values := strings.Split(msg.Data, ",")

	switch msg.Type {
	case EventWorkspace:
		hyprlCtl.syncWorkspaces()
		break
	case EventWorkspacev2:
		hyprlCtl.syncWorkspaces()
		break
	case EventFocusedMonitor:
		break
	case EventActiveWindow:
		preprocessClientAddress(values)
		hyprlCtl.syncActiveClient()
		break
	case EventActiveWindowv2:
		preprocessClientAddress(values)
		hyprlCtl.syncActiveClient()
		break
	case EventFullscreen:
		hyprlCtl.syncWorkspaces()
		hyprlCtl.syncClients()
		break
	case EventMonitorRemoved:
		break
	case EventMonitorAdded:
		break
	case EventMonitorAddedv2:
		break
	case EventCreateWorkspace:
		hyprlCtl.syncWorkspaces()
		break
	case EventCreateWorkspacev2:
		hyprlCtl.syncWorkspaces()
		break
	case EventDestroyWorkspace:
		hyprlCtl.syncWorkspaces()
		break
	case EventDestroyWorkspacev2:
		hyprlCtl.syncWorkspaces()
		break
	case EventMoveWorkspace:
		hyprlCtl.syncWorkspaces()
		break
	case EventMoveWorkspacev2:
		hyprlCtl.syncWorkspaces()
		break
	case EventRenameWorkspace:
		break
	case EventActiveSpecial:
		break
	case EventActiveLayout:
		break
	case EventOpenWindow:
		preprocessClientAddress(values)
		hyprlCtl.syncWorkspaces()
		hyprlCtl.syncClients()
		break
	case EventCloseWindow:
		preprocessClientAddress(values)
		hyprlCtl.syncWorkspaces()
		hyprlCtl.syncClients()
		break
	case EventMoveWindow:
		preprocessClientAddress(values)
		hyprlCtl.syncWorkspaces()
		hyprlCtl.syncClients()
		break
	case EventMoveWindowv2:
		preprocessClientAddress(values)
		hyprlCtl.syncWorkspaces()
		hyprlCtl.syncClients()
		break
	case EventOpenLayer:
		break
	case EventCloseLayer:
		break
	case EventSubMap:
		break
	case EventChangeFloatingMode:
		preprocessClientAddress(values)
		hyprlCtl.syncClients()
		break
	case EventUrgent:
		preprocessClientAddress(values)
		break
	case EventMinimize:
		preprocessClientAddress(values)
		break
	case EventScreencast:
		break
	case EventWindowTitle:
		preprocessClientAddress(values)
		hyprlCtl.syncClients()
		break
	case EventIgnoreGroupLock:
		break
	case EventLockGroups:
		break
	case EventConfigReloaded:
		break
	case EventPin:
		preprocessClientAddress(values)
		hyprlCtl.syncClients()
		break
	}

	eventMethod := eventMethods[msg.Type]
	eventMethod.call(msg.Type, values)
}

type EventType string

type EventData struct {
	Type EventType
	Data string
}

type HyprlandEventHandler interface {
	Workspace(workspaceName string)
	WorkspaceV2(workspaceId int, workspaceName string)
	FocusedMonitor(monitorName string, workspaceName string)
	ActiveWindow(windowClass string, windowTitle string)
	ActiveWindowV2(windowAddress string)
	Fullscreen(fullscreen bool)
	MonitorRemoved(monitorName string)
	MonitorAdded(monitorName string)
	MonitorAddedV2(monitorId int, monitorName string, monitorDescription string)
	CreateWorkspace(workspaceName string)
	CreateWorkspaceV2(workspaceId int, workspaceName string)
	DestroyWorkspace(workspaceName string)
	DestroyWorkspaceV2(workspaceId int, workspaceName string)
	MoveWorkspace(workspaceName string, monitorName string)
	MoveWorkspaceV2(workspaceId int, workspaceName string, monitorName string)
	RenameWorkspace(workspaceId int, newWorkspaceName string)
	ActiveSpecial(workspaceName string, monitorName string)
	ActiveLayout(keyboardName string, layoutName string)
	OpenWindow(windowAddress string, workspaceName string, windowClass string, windowTitle string)
	CloseWindow(windowAddress string)
	MoveWindow(windowAddress string, workspaceName string)
	MoveWindowV2(windowAddress string, workspaceId int, workspaceName string)
	OpenLayer(namespace string)
	CloseLayer(namespace string)
	SubMap(submapName string)
	ChangeFloatingMode(windowAddress string, floating bool)
	Urgent(windowAddress string)
	Minimize(windowAddress string, minimized bool)
	Screencast(screencasting bool, shareWindow bool)
	WindowTitle(windowAddress string)
	IgnoreGroupLock(ignoringGrouplock bool)
	LockGroups(lockingGroups bool)
	ConfigReloaded()
	Pin(windowAddress string, pinned bool)
}

const (
	EventWorkspace          EventType = "workspace"
	EventWorkspacev2        EventType = "workspacev2"
	EventFocusedMonitor     EventType = "focusedmon"
	EventActiveWindow       EventType = "activewindow"
	EventActiveWindowv2     EventType = "activewindowv2"
	EventFullscreen         EventType = "fullscreen"
	EventMonitorRemoved     EventType = "monitorremoved"
	EventMonitorAdded       EventType = "monitoradded"
	EventMonitorAddedv2     EventType = "monitoraddedv2"
	EventCreateWorkspace    EventType = "createworkspace"
	EventCreateWorkspacev2  EventType = "createworkspacev2"
	EventDestroyWorkspace   EventType = "destroyworkspace"
	EventDestroyWorkspacev2 EventType = "destroyworkspacev2"
	EventMoveWorkspace      EventType = "moveworkspace"
	EventMoveWorkspacev2    EventType = "moveworkspacev2"
	EventRenameWorkspace    EventType = "renameworkspace"
	EventActiveSpecial      EventType = "activespecial"
	EventActiveLayout       EventType = "activelayout"
	EventOpenWindow         EventType = "openwindow"
	EventCloseWindow        EventType = "closewindow"
	EventMoveWindow         EventType = "movewindow"
	EventMoveWindowv2       EventType = "movewindowv2"
	EventOpenLayer          EventType = "openlayer"
	EventCloseLayer         EventType = "closelayer"
	EventSubMap             EventType = "submap"
	EventChangeFloatingMode EventType = "changefloatingmode"
	EventUrgent             EventType = "urgent"
	EventMinimize           EventType = "minimize"
	EventScreencast         EventType = "screencast"
	EventWindowTitle        EventType = "windowtitle"
	EventIgnoreGroupLock    EventType = "ignoregrouplock"
	EventLockGroups         EventType = "lockgroups"
	EventConfigReloaded     EventType = "configreloaded"
	EventPin                EventType = "pin"
)

var eventMethods = newEventMethods()

func newEventMethods() map[EventType]*eventMethod {
	iface := reflect.TypeOf(struct{ HyprlandEventHandler }{})
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
			case reflect.Int:
				eventMethod.values[j-1] = toInt
				break
			case reflect.Bool:
				eventMethod.values[j-1] = toBool
				break
			default:
				eventMethod.values[j-1] = toString
				break
			}
		}
		eventMethods[EventType(strings.ToLower(method.Name))] = eventMethod
	}
	return eventMethods
}

func (m *eventMethod) call(eventType EventType, values []string) {
	in := make([]reflect.Value, len(m.values)+1)
	for i, value := range values {
		in[i+1] = m.values[i](value)
	}

	for _, subscriber := range eventSubscribers[eventType] {
		in[0] = reflect.ValueOf(subscriber.handle)
		subscriber.callback.Call(in)
	}
}

type eventMethod struct {
	name   string
	values []valueConvertor
}

type valueConvertor = func(value string) reflect.Value

func toInt(value string) reflect.Value {
	intValue, _ := strconv.Atoi(value)
	return reflect.ValueOf(intValue)
}

func toBool(value string) reflect.Value {
	return reflect.ValueOf(value == "1")
}

func toString(value string) reflect.Value {
	return reflect.ValueOf(value)
}
