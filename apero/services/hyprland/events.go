package hyprland

import (
	"log"
	"net"
	"reflect"
	"strconv"
	"strings"
	"webflo-dev/apero/logger"
)

type hyprlandEventService struct {
	connection  net.Conn
	subscribers map[EventType][]EventHandler
}

const (
	bufferSize     = 8192
	eventSeperator = ">>"
)

var eventService *hyprlandEventService

func WatchEvents(ev EventHandler, events ...EventType) {
	if eventService == nil {
		eventService = &hyprlandEventService{
			connection:  createEventsConnection(),
			subscribers: make(map[EventType][]EventHandler),
		}
		syncWorkspaces()
	}

	for _, eventType := range events {
		eventService.subscribers[eventType] = append(eventService.subscribers[eventType], ev)
	}

	go func() {
		defer eventService.closeConnection()

		for {
			msg, err := eventService.receive()
			if err != nil {
				logger.AppLogger.Println("Error receiving message", err)
				return
			}

			for _, data := range msg {
				log.Printf("%+v\n", data)
				eventService.processEvent(data)
			}
		}
	}()
}

var workspaces = []Workspace{}

type SyncMethod func()

type SyncMap struct {
	methodName  string
	syncMethods []SyncMethod
}

var eventMap = map[EventType]SyncMap{
	EventWorkspace:          {"Workspace", nil},
	EventWorkspacev2:        {"WorkspaceV2", nil},
	EventFocusedMonitor:     {"FocusedMonitor", nil},
	EventActiveWindow:       {"ActiveWindow", nil},
	EventActiveWindowv2:     {"ActiveWindowV2", nil},
	EventFullscreen:         {"Fullscreen", []SyncMethod{syncWorkspaces}},
	EventMonitorRemoved:     {"MonitorRemoved", nil},
	EventMonitorAdded:       {"MonitorAdded", nil},
	EventMonitorAddedv2:     {"MonitorAddedV2", nil},
	EventCreateWorkspace:    {"CreateWorkspace", []SyncMethod{syncWorkspaces}},
	EventCreateWorkspacev2:  {"CreateWorkspaceV2", []SyncMethod{syncWorkspaces}},
	EventDestroyWorkspace:   {"DestroyWorkspace", []SyncMethod{syncWorkspaces}},
	EventDestroyWorkspacev2: {"DestroyWorkspaceV2", []SyncMethod{syncWorkspaces}},
	EventMoveWorkspace:      {"MoveWorkspace", []SyncMethod{syncWorkspaces}},
	EventMoveWorkspacev2:    {"MoveWorkspaceV2", []SyncMethod{syncWorkspaces}},
	EventRenameWorkspace:    {"RenameWorkspace", []SyncMethod{syncWorkspaces}},
	EventActiveSpecial:      {"ActiveSepcial", nil},
	EventActiveLayout:       {"ActiveLayout", nil},
	EventOpenWindow:         {"OpenWindow", []SyncMethod{syncWorkspaces}},
	EventCloseWindow:        {"CloseWindow", []SyncMethod{syncWorkspaces}},
	EventMoveWindow:         {"MoveWindow", []SyncMethod{syncWorkspaces}},
	EventMoveWindowv2:       {"MoveWindowV2", []SyncMethod{syncWorkspaces}},
	EventOpenLayer:          {"OpenLayer", nil},
	EventCloseLayer:         {"CloseLayer", nil},
	EventSubMap:             {"SubMap", nil},
	EventChangeFloatingMode: {"ChangeFloatingMode", nil},
	EventUrgent:             {"Urgent", nil},
	EventMinimize:           {"Minimize", nil},
	EventScreencast:         {"Screencast", nil},
	EventWindowtitle:        {"WindowTitle", nil},
	EventIgnoreGroupLock:    {"IgnoreGrouplock", nil},
	EventLockGroups:         {"LockGroups", nil},
	EventConfigreloaded:     {"ConfigReloaded", nil},
	EventPin:                {"Pin", nil},
}

func (s *hyprlandEventService) processEvent(msg EventData) {
	rawValues := strings.Split(msg.Data, ",")
	eventData := eventMap[msg.Type]

	if eventData.syncMethods != nil {
		for _, syncMethod := range eventData.syncMethods {
			syncMethod()
		}
	}
	for _, subscriber := range s.subscribers[msg.Type] {
		CallHandler(subscriber, eventData.methodName, rawValues)
	}
}

func CallHandler(handler interface{}, target string, rawValues []string) {
	method := reflect.ValueOf(handler).MethodByName(target)
	methodType := method.Type()

	in := make([]reflect.Value, methodType.NumIn())

	for i := 0; i < method.Type().NumIn(); i++ {
		value := rawValues[i]
		t := methodType.In(i)
		switch t.Kind() {
		case reflect.Bool:
			in[i] = reflect.ValueOf(value == "1")
			break
		case reflect.Int:
			intValue, _ := strconv.Atoi(value)
			in[i] = reflect.ValueOf(intValue)
		default:
			in[i] = reflect.ValueOf(value)
			break
		}
	}
	method.Call(in)
}

func syncWorkspaces() {
	workspaces, _ = Workspaces()
}

type EventType string

type EventData struct {
	Type EventType
	Data string
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
	EventWindowtitle        EventType = "windowtitle"
	EventIgnoreGroupLock    EventType = "ignoregrouplock"
	EventLockGroups         EventType = "lockgroups"
	EventConfigreloaded     EventType = "configreloaded"
	EventPin                EventType = "pin"
)

type EventHandler interface {
	Workspace(workspaceName string)
	WorkspaceV2(workspaceId int, workspaceName string)
	FocusedMonitor(monitorName string, workspaceName string)
	ActiveWindow(windoowClass string, windowTitle string)
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
	ActiveSepcial(workspaceName string, monitorName string)
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
	IgnoreGrouplock(ignoringGrouplock bool)
	LockGroups(lockingGroups bool)
	ConfigReloaded()
	Pin(windowAddress string, pinned bool)
}
