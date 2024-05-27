package hyprland

import (
	"fmt"
	"strings"
	"webflo-dev/apero/events"
)

var _service = newService()

func StartService() {
	_service.syncWorkspaces()
	_service.syncClients()

	_service.start()
}

//
// Event registration
//

func OnWorkspace(f func(payload PayloadWorkspace)) (events.ID, error) {
	return registerEvt(EventWorkspace, f)
}
func OnWorkspacev2(f func(payload PayloadWorkspaceV2)) (events.ID, error) {
	return registerEvt(EventWorkspaceV2, f)
}
func OnFocusedMonitor(f func(payload PayloadFocusedMonitor)) (events.ID, error) {
	return registerEvt(EventFocusedMonitor, f)
}
func OnActiveWindow(f func(payload PayloadActiveWindow)) (events.ID, error) {
	return registerEvt(EventActiveWindow, f)
}
func OnActiveWindowv2(f func(payload PayloadActiveWindowV2)) (events.ID, error) {
	return registerEvt(EventActiveWindowV2, f)
}
func OnFullscreen(f func(payload PayloadFullscreen)) (events.ID, error) {
	return registerEvt(EventFullscreen, f)
}
func OnMonitorRemoved(f func(payload PayloadMonitorRemoved)) (events.ID, error) {
	return registerEvt(EventMonitorRemoved, f)
}
func OnMonitorAdded(f func(payload PayloadMonitorAdded)) (events.ID, error) {
	return registerEvt(EventMonitorAdded, f)
}
func OnMonitorAddedv2(f func(payload PayloadMonitorAddedV2)) (events.ID, error) {
	return registerEvt(EventMonitorAddedV2, f)
}
func OnCreateWorkspace(f func(payload PayloadCreateWorkspace)) (events.ID, error) {
	return registerEvt(EventCreateWorkspace, f)
}
func OnCreateWorkspacev2(f func(payload PayloadCreateWorkspaceV2)) (events.ID, error) {
	return registerEvt(EventCreateWorkspaceV2, f)
}
func OnDestroyWorkspace(f func(payload PayloadDestroyWorkspace)) (events.ID, error) {
	return registerEvt(EventDestroyWorkspace, f)
}
func OnDestroyWorkspacev2(f func(payload PayloadDestroyWorkspaceV2)) (events.ID, error) {
	return registerEvt(EventDestroyWorkspaceV2, f)
}
func OnMoveWorkspace(f func(payload PayloadMoveWorkspace)) (events.ID, error) {
	return registerEvt(EventMoveWorkspace, f)
}
func OnMoveWorkspacev2(f func(payload PayloadMoveWorkspaceV2)) (events.ID, error) {
	return registerEvt(EventMoveWorkspaceV2, f)
}
func OnRenameWorkspace(f func(payload PayloadRenameWorkspace)) (events.ID, error) {
	return registerEvt(EventRenameWorkspace, f)
}
func OnActiveSpecial(f func(payload PayloadActiveSpecial)) (events.ID, error) {
	return registerEvt(EventActiveSpecial, f)
}
func OnActiveLayout(f func(payload PayloadActiveLayout)) (events.ID, error) {
	return registerEvt(EventActiveLayout, f)
}
func OnOpenWindow(f func(payload PayloadOpenWindow)) (events.ID, error) {
	return registerEvt(EventOpenWindow, f)
}
func OnCloseWindow(f func(payload PayloadCloseWindow)) (events.ID, error) {
	return registerEvt(EventCloseWindow, f)
}
func OnMoveWindow(f func(payload PayloadMoveWindow)) (events.ID, error) {
	return registerEvt(EventMoveWindow, f)
}
func OnMoveWindowv2(f func(payload PayloadMoveWindowV2)) (events.ID, error) {
	return registerEvt(EventMoveWindowV2, f)
}
func OnOpenLayer(f func(payload PayloadOpenLayer)) (events.ID, error) {
	return registerEvt(EventOpenLayer, f)
}
func OnCloseLayer(f func(payload PayloadCloseLayer)) (events.ID, error) {
	return registerEvt(EventCloseLayer, f)
}
func OnSubMap(f func(payload PayloadSubMap)) (events.ID, error) {
	return registerEvt(EventSubMap, f)
}
func OnChangeFloatingMode(f func(payload PayloadChangeFloatingMode)) (events.ID, error) {
	return registerEvt(EventChangeFloatingMode, f)
}
func OnUrgent(f func(payload PayloadUrgent)) (events.ID, error) {
	return registerEvt(EventUrgent, f)
}
func OnMinimize(f func(payload PayloadMinimize)) (events.ID, error) {
	return registerEvt(EventMinimize, f)
}
func OnScreencast(f func(payload PayloadScreencast)) (events.ID, error) {
	return registerEvt(EventScreencast, f)
}
func OnWindowTitle(f func(payload PayloadWindowTitle)) (events.ID, error) {
	return registerEvt(EventWindowTitle, f)
}
func OnIgnoreGroupLock(f func(payload PayloadIgnoreGroupLock)) (events.ID, error) {
	return registerEvt(EventIgnoreGroupLock, f)
}
func OnLockGroups(f func(payload PayloadLockGroups)) (events.ID, error) {
	return registerEvt(EventLockGroups, f)
}
func OnConfigReloaded(f func()) (events.ID, error) {
	return registerEvt(EventConfigReloaded, func(_ any) { f() })
}
func OnPin(f func(payload PayloadPin)) (events.ID, error) {
	return registerEvt(EventPin, f)
}

//
// API methods
//

func Dispatch(dispatch string, args ...any) error {
	cmd := []string{"dispatch", dispatch}
	for _, arg := range args {
		cmd = append(cmd, fmt.Sprintf("%v", arg))
	}
	return writeCmd(strings.Join(cmd, " "), nil)
}

func Workspaces() []Workspace {
	return _service.workspaces
}
func ActiveWorkspace() Workspace {
	return _service.activeWorkspace
}

func Clients() []Client {
	return _service.clients
}

func ActiveClient() Client {
	return _service.activeClient
}

// func Monitors() []Monitor {
// 	return _service.monitors
// }

// func ActiveMonitor() Monitor {
// 	return _service.activeMonitor
// }

func Binds() ([]Bind, error) {
	var binds []Bind
	err := writeCmd("j/binds", &binds)
	return binds, err
}

func ConfigErrors() ([]string, error) {
	var errors []string
	err := writeCmd("j/configerrors", &errors)
	return errors, err
}

func CursorPosition() (CursorPos, error) {
	var pos CursorPos
	err := writeCmd("j/configerrors", &pos)
	return pos, err
}

func GetDevices() (Devices, error) {
	var devices Devices
	err := writeCmd("j/devices", &devices)
	return devices, err
}

func GetInstances() ([]Instance, error) {
	var instances []Instance
	err := writeCmd("j/instances", &instances)
	return instances, err
}

func Reload() error {
	return writeCmd("reload", nil)
}

func Keyword(keyword string, args ...any) error {
	cmd := []string{"keyword", keyword}
	for _, arg := range args {
		cmd = append(cmd, fmt.Sprintf("%v", arg))
	}
	return writeCmd(strings.Join(cmd, " "), nil)
}

func Layers() ([]Layer, error) {
	var values map[string]interface{}
	err := writeCmd("j/layers", &values)
	if err != nil {
		return nil, err
	}

	var allLayers []Layer

	for monitorName, perMonitor := range values {
		levels := (perMonitor).(map[string]interface{})["levels"]
		for level, rawLayers := range levels.(map[string]interface{}) {
			layers := toStruct[[]Layer](rawLayers)

			for _, layer := range layers {
				layer.MonitorName = monitorName
				layer.Layer = LayerType(level)
				allLayers = append(allLayers, layer)
			}
		}
	}
	return allLayers, nil
}

func Layouts() ([]string, error) {
	var layouts []string
	err := writeCmd("j/layouts", &layouts)
	return layouts, err
}
