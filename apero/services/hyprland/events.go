package hyprland

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
