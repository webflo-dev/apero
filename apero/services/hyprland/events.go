package hyprland

import (
	"strings"
)

type hyprlandEventService struct {
	hyprlandIpcService
	listening bool
	// subscribers map[EventType][]HyprlandEventHandler
	subscribers map[EventType][]any
}

type SyncMethod func()
type PreProcessor func([]string) []string

type ProcessMap struct {
	methodName   string
	preprocessor []PreProcessor
	syncMethods  []SyncMethod
}

var eventService = newHyprlandEventService()

func newHyprlandEventService() *hyprlandEventService {
	service := &hyprlandEventService{
		listening: false,
		// subscribers: make(map[EventType][]HyprlandEventHandler),
		subscribers: make(map[EventType][]any),
	}

	return service
}

func (s *hyprlandEventService) processEvent(msg EventData) {
	eventData := eventMap[msg.Type]

	rawValues := strings.Split(msg.Data, ",")
	if eventData.preprocessor != nil {
		for _, preprocessor := range eventData.preprocessor {
			rawValues = preprocessor(rawValues)
		}
	}

	if eventData.syncMethods != nil {
		for _, syncMethod := range eventData.syncMethods {
			syncMethod()
		}
	}
	for _, subscriber := range s.subscribers[msg.Type] {
		callHandler(subscriber, eventData.methodName, rawValues)
	}
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

var eventMap = map[EventType]ProcessMap{
	// EventWorkspace:          {"Workspace", nil, nil},

	EventWorkspacev2: {"WorkspaceV2",
		nil,
		[]SyncMethod{hyprlCtl.syncWorkspaces}},

	EventFocusedMonitor: {"FocusedMonitor", nil, nil},

	// EventActiveWindow:       {"ActiveWindow", nil, nil},

	EventActiveWindowv2: {"ActiveWindowV2",
		[]PreProcessor{hyprlCtl.preprocessClientAddress},
		[]SyncMethod{hyprlCtl.syncActiveClient}},

	EventFullscreen: {"Fullscreen",
		nil,
		[]SyncMethod{hyprlCtl.syncWorkspaces, hyprlCtl.syncClients}},

	EventMonitorRemoved: {"MonitorRemoved", nil, nil},

	// EventMonitorAdded:   {"MonitorAdded", nil, nil},

	EventMonitorAddedv2: {"MonitorAddedV2", nil, nil},

	// EventCreateWorkspace:    {"CreateWorkspace", nil, nil},

	EventCreateWorkspacev2: {"CreateWorkspaceV2",
		nil,
		[]SyncMethod{hyprlCtl.syncWorkspaces}},

	// EventDestroyWorkspace:   {"DestroyWorkspace", nil, nil},

	EventDestroyWorkspacev2: {"DestroyWorkspaceV2",
		nil,
		[]SyncMethod{hyprlCtl.syncWorkspaces}},

	// EventMoveWorkspace:      {"MoveWorkspace", nil, []SyncMethod{hyprlCtl.syncWorkspaces}},

	EventMoveWorkspacev2: {"MoveWorkspaceV2",
		nil,
		[]SyncMethod{hyprlCtl.syncWorkspaces}},

	EventRenameWorkspace: {"RenameWorkspace",
		nil,
		[]SyncMethod{hyprlCtl.syncWorkspaces}},

	EventActiveSpecial: {"ActiveSepcial", nil, nil},

	EventActiveLayout: {"ActiveLayout", nil, nil},

	EventOpenWindow: {"OpenWindow",
		[]PreProcessor{hyprlCtl.preprocessClientAddress},
		[]SyncMethod{hyprlCtl.syncWorkspaces, hyprlCtl.syncClients}},

	EventCloseWindow: {"CloseWindow",
		[]PreProcessor{hyprlCtl.preprocessClientAddress},
		[]SyncMethod{hyprlCtl.syncWorkspaces, hyprlCtl.syncClients}},
	// EventMoveWindow:         {"MoveWindow", nil, nil},

	EventMoveWindowv2: {"MoveWindowV2",
		[]PreProcessor{hyprlCtl.preprocessClientAddress},
		[]SyncMethod{hyprlCtl.syncWorkspaces, hyprlCtl.syncClients}},

	EventOpenLayer: {"OpenLayer", nil, nil},

	EventCloseLayer: {"CloseLayer", nil, nil},

	EventSubMap: {"SubMap", nil, nil},

	EventChangeFloatingMode: {"ChangeFloatingMode",
		[]PreProcessor{hyprlCtl.preprocessClientAddress},
		[]SyncMethod{hyprlCtl.syncClients}},

	EventUrgent: {"Urgent",
		[]PreProcessor{hyprlCtl.preprocessClientAddress},
		nil},

	EventMinimize: {"Minimize",
		[]PreProcessor{hyprlCtl.preprocessClientAddress},
		nil},

	EventScreencast: {"Screencast", nil, nil},

	EventWindowtitle: {"WindowTitle",
		[]PreProcessor{hyprlCtl.preprocessClientAddress},
		[]SyncMethod{hyprlCtl.syncClients}},

	EventIgnoreGroupLock: {"IgnoreGrouplock", nil, nil},

	EventLockGroups: {"LockGroups", nil, nil},

	EventConfigreloaded: {"ConfigReloaded", nil, nil},

	EventPin: {"Pin",
		[]PreProcessor{hyprlCtl.preprocessClientAddress},
		[]SyncMethod{hyprlCtl.syncClients}},
}
