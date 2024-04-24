package hyprland

import (
	"strings"
)

type hyprlandEventService struct {
	hyprlandIpcService
	listening   bool
	subscribers map[EventType][]HyprlandEventHandler
}

type SyncMethod func()

type SyncMap struct {
	methodName  string
	syncMethods []SyncMethod
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

var eventMap = map[EventType]SyncMap{
	EventWorkspace:          {"Workspace", []SyncMethod{hyprlCtl.syncWorkspaces}},
	EventWorkspacev2:        {"WorkspaceV2", []SyncMethod{hyprlCtl.syncWorkspaces}},
	EventFocusedMonitor:     {"FocusedMonitor", nil},
	EventActiveWindow:       {"ActiveWindow", nil},
	EventActiveWindowv2:     {"ActiveWindowV2", nil},
	EventFullscreen:         {"Fullscreen", []SyncMethod{hyprlCtl.syncWorkspaces}},
	EventMonitorRemoved:     {"MonitorRemoved", nil},
	EventMonitorAdded:       {"MonitorAdded", nil},
	EventMonitorAddedv2:     {"MonitorAddedV2", nil},
	EventCreateWorkspace:    {"CreateWorkspace", []SyncMethod{hyprlCtl.syncWorkspaces}},
	EventCreateWorkspacev2:  {"CreateWorkspaceV2", []SyncMethod{hyprlCtl.syncWorkspaces}},
	EventDestroyWorkspace:   {"DestroyWorkspace", []SyncMethod{hyprlCtl.syncWorkspaces}},
	EventDestroyWorkspacev2: {"DestroyWorkspaceV2", []SyncMethod{hyprlCtl.syncWorkspaces}},
	EventMoveWorkspace:      {"MoveWorkspace", []SyncMethod{hyprlCtl.syncWorkspaces}},
	EventMoveWorkspacev2:    {"MoveWorkspaceV2", []SyncMethod{hyprlCtl.syncWorkspaces}},
	EventRenameWorkspace:    {"RenameWorkspace", []SyncMethod{hyprlCtl.syncWorkspaces}},
	EventActiveSpecial:      {"ActiveSepcial", nil},
	EventActiveLayout:       {"ActiveLayout", nil},
	EventOpenWindow:         {"OpenWindow", []SyncMethod{hyprlCtl.syncWorkspaces}},
	EventCloseWindow:        {"CloseWindow", []SyncMethod{hyprlCtl.syncWorkspaces}},
	EventMoveWindow:         {"MoveWindow", []SyncMethod{hyprlCtl.syncWorkspaces}},
	EventMoveWindowv2:       {"MoveWindowV2", []SyncMethod{hyprlCtl.syncWorkspaces}},
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
