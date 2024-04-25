package hyprland

import (
	"fmt"
	"log"
	"webflo-dev/apero/logger"
)

func (service *hyprlandEventService) listen() {
	if service.listening {
		return
	}

	go func() {
		service.createEventsConnection()
		defer service.closeConnection()

		logger.AppLogger.Println("listening for hyprland events")

		for {
			msg, err := service.readEvent()
			if err != nil {
				logger.AppLogger.Println("Error receiving message", err)
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

func ActiveClient() Client {
	return hyprlCtl.activeClient
}

func WatchEvents(handler HyprlandEventHandler, events ...EventType) {
	if eventService.listening == false {
		eventService.listen()
	}

	for _, eventType := range events {
		eventService.subscribers[eventType] = append(eventService.subscribers[eventType], handler)
	}
}

type HyprlandEventHandler interface {
	// Workspace(workspaceName string)
	WorkspaceV2(workspaceId int, workspaceName string)
	FocusedMonitor(monitorName string, workspaceName string)
	// ActiveWindow(windowClass string, windowTitle string)
	ActiveWindowV2(windowAddress string)
	Fullscreen(fullscreen bool)
	MonitorRemoved(monitorName string)
	// MonitorAdded(monitorName string)
	MonitorAddedV2(monitorId int, monitorName string, monitorDescription string)
	// CreateWorkspace(workspaceName string)
	CreateWorkspaceV2(workspaceId int, workspaceName string)
	// DestroyWorkspace(workspaceName string)
	DestroyWorkspaceV2(workspaceId int, workspaceName string)
	// MoveWorkspace(workspaceName string, monitorName string)
	MoveWorkspaceV2(workspaceId int, workspaceName string, monitorName string)
	RenameWorkspace(workspaceId int, newWorkspaceName string)
	ActiveSepcial(workspaceName string, monitorName string)
	ActiveLayout(keyboardName string, layoutName string)
	OpenWindow(windowAddress string, workspaceName string, windowClass string, windowTitle string)
	CloseWindow(windowAddress string)
	// MoveWindow(windowAddress string, workspaceName string)
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
