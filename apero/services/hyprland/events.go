package hyprland

import (
	"fmt"
	"strconv"
	"strings"
)

func (s *service) processEvent(msg EventData) {
	values := strings.Split(msg.Data, ",")

	switch msg.Type {
	case EventWorkspace:
		s.syncWorkspaces()
		for _, subscriber := range s.subscribers[EventWorkspace] {
			subscriber.Workspace(values[0])
		}
		break
	case EventWorkspacev2:
		s.syncWorkspaces()
		for _, subscriber := range s.subscribers[EventWorkspacev2] {
			subscriber.WorkspaceV2(toInt(values[0]), values[1])
		}
		break
	case EventFocusedMonitor:
		for _, subscriber := range s.subscribers[EventFocusedMonitor] {
			subscriber.FocusedMonitor(values[0], values[1])
		}
		break
	case EventActiveWindow:
		s.syncActiveClient()
		for _, subscriber := range s.subscribers[EventActiveWindow] {
			subscriber.ActiveWindow(values[0], values[1])
		}
		break
	case EventActiveWindowv2:
		s.syncActiveClient()
		for _, subscriber := range s.subscribers[EventActiveWindowv2] {
			subscriber.ActiveWindowV2(toAddress(values[0]))
		}
		break
	case EventFullscreen:
		s.syncWorkspaces()
		s.syncClients()
		for _, subscriber := range s.subscribers[EventFullscreen] {
			subscriber.Fullscreen(toBool(values[0]))
		}
		break
	case EventMonitorRemoved:
		for _, subscriber := range s.subscribers[EventMonitorRemoved] {
			subscriber.MonitorRemoved(values[0])
		}
		break
	case EventMonitorAdded:
		for _, subscriber := range s.subscribers[EventMonitorAdded] {
			subscriber.MonitorAdded(values[0])
		}
		break
	case EventMonitorAddedv2:
		for _, subscriber := range s.subscribers[EventMonitorAddedv2] {
			subscriber.MonitorAddedV2(toInt(values[0]), values[1], values[2])
		}
		break
	case EventCreateWorkspace:
		s.syncWorkspaces()
		for _, subscriber := range s.subscribers[EventCreateWorkspace] {
			subscriber.CreateWorkspace(values[0])
		}
		break
	case EventCreateWorkspacev2:
		s.syncWorkspaces()
		for _, subscriber := range s.subscribers[EventCreateWorkspacev2] {
			subscriber.CreateWorkspaceV2(toInt(values[0]), values[1])
		}
		break
	case EventDestroyWorkspace:
		s.syncWorkspaces()
		for _, subscriber := range s.subscribers[EventDestroyWorkspace] {
			subscriber.DestroyWorkspace(values[0])
		}
		break
	case EventDestroyWorkspacev2:
		s.syncWorkspaces()
		for _, subscriber := range s.subscribers[EventDestroyWorkspacev2] {
			subscriber.DestroyWorkspaceV2(toInt(values[0]), values[1])
		}
		break
	case EventMoveWorkspace:
		s.syncWorkspaces()
		for _, subscriber := range s.subscribers[EventMoveWorkspace] {
			subscriber.MoveWorkspace(values[0], values[1])
		}
		break
	case EventMoveWorkspacev2:
		s.syncWorkspaces()
		for _, subscriber := range s.subscribers[EventMoveWorkspacev2] {
			subscriber.MoveWorkspaceV2(toInt(values[0]), values[1], values[2])
		}
		break
	case EventRenameWorkspace:
		s.syncWorkspaces()
		for _, subscriber := range s.subscribers[EventRenameWorkspace] {
			subscriber.RenameWorkspace(toInt(values[0]), values[1])
		}
		break
	case EventActiveSpecial:
		for _, subscriber := range s.subscribers[EventActiveSpecial] {
			subscriber.ActiveSpecial(values[0], values[1])
		}
		break
	case EventActiveLayout:
		for _, subscriber := range s.subscribers[EventActiveLayout] {
			subscriber.ActiveLayout(values[0], values[1])
		}
		break
	case EventOpenWindow:
		s.syncWorkspaces()
		s.syncClients()
		for _, subscriber := range s.subscribers[EventOpenWindow] {
			subscriber.OpenWindow(toAddress(values[0]), values[1], values[2], values[3])
		}
		break
	case EventCloseWindow:
		s.syncWorkspaces()
		s.syncClients()
		for _, subscriber := range s.subscribers[EventCloseWindow] {
			subscriber.CloseWindow(toAddress(values[0]))
		}
		break
	case EventMoveWindow:
		s.syncWorkspaces()
		s.syncClients()
		for _, subscriber := range s.subscribers[EventMoveWindow] {
			subscriber.MoveWindow(toAddress(values[0]), values[1])
		}
		break
	case EventMoveWindowv2:
		s.syncWorkspaces()
		s.syncClients()
		for _, subscriber := range s.subscribers[EventMoveWindowv2] {
			subscriber.MoveWindowV2(toAddress(values[0]), toInt(values[1]), values[2])
		}
		break
	case EventOpenLayer:
		for _, subscriber := range s.subscribers[EventOpenLayer] {
			subscriber.OpenLayer(values[0])
		}
		break
	case EventCloseLayer:
		for _, subscriber := range s.subscribers[EventCloseLayer] {
			subscriber.CloseLayer(values[0])
		}
		break
	case EventSubMap:
		for _, subscriber := range s.subscribers[EventSubMap] {
			subscriber.SubMap(values[0])
		}
		break
	case EventChangeFloatingMode:
		s.syncClients()
		for _, subscriber := range s.subscribers[EventChangeFloatingMode] {
			subscriber.ChangeFloatingMode(toAddress(values[0]), toBool(values[1]))
		}
		break
	case EventUrgent:
		for _, subscriber := range s.subscribers[EventUrgent] {
			subscriber.Urgent(toAddress(values[0]))
		}
		break
	case EventMinimize:
		for _, subscriber := range s.subscribers[EventMinimize] {
			subscriber.Minimize(toAddress(values[0]), toBool(values[1]))
		}
		break
	case EventScreencast:
		for _, subscriber := range s.subscribers[EventScreencast] {
			subscriber.Screencast(toBool(values[0]), toBool(values[1]))
		}
		break
	case EventWindowTitle:
		s.syncClients()
		for _, subscriber := range s.subscribers[EventWindowTitle] {
			subscriber.WindowTitle(toAddress(values[0]))
		}
		break
	case EventIgnoreGroupLock:
		for _, subscriber := range s.subscribers[EventIgnoreGroupLock] {
			subscriber.IgnoreGroupLock(toBool(values[0]))
		}
		break
	case EventLockGroups:
		for _, subscriber := range s.subscribers[EventLockGroups] {
			subscriber.LockGroups(toBool(values[0]))
		}
		break
	case EventConfigReloaded:
		for _, subscriber := range s.subscribers[EventConfigReloaded] {
			subscriber.ConfigReloaded()
		}
		break
	case EventPin:
		s.syncClients()
		for _, subscriber := range s.subscribers[EventPin] {
			subscriber.Pin(toAddress(values[0]), toBool(values[1]))
		}
		break
	}
}

type EventData struct {
	Type EventType
	Data string
}

type Subscriber interface {
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

type EventType string

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

func toInt(value string) int {
	ret, _ := strconv.Atoi(value)
	return ret
}

func toBool(value string) bool {
	return value == "1"
}

func toAddress(value string) string {
	return fmt.Sprintf("0x%s", value)
}
