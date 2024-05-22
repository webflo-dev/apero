package hyprland

import (
	"fmt"
	"strings"
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

func OnWorkspace(id string, f func(payload PayloadWorkspace)) {
	registerEvt(EventWorkspace, id, f)
}
func OnWorkspacev2(id string, f func(payload PayloadWorkspaceV2)) {
	registerEvt(EventWorkspaceV2, id, f)
}
func OnFocusedMonitor(id string, f func(payload PayloadFocusedMonitor)) {
	registerEvt(EventFocusedMonitor, id, f)
}
func OnActiveWindow(id string, f func(payload PayloadActiveWindow)) {
	registerEvt(EventActiveWindow, id, f)
}
func OnActiveWindowv2(id string, f func(payload PayloadActiveWindowV2)) {
	registerEvt(EventActiveWindowV2, id, f)
}
func OnFullscreen(id string, f func(payload PayloadFullscreen)) {
	registerEvt(EventFullscreen, id, f)
}
func OnMonitorRemoved(id string, f func(payload PayloadMonitorRemoved)) {
	registerEvt(EventMonitorRemoved, id, f)
}
func OnMonitorAdded(id string, f func(payload PayloadMonitorAdded)) {
	registerEvt(EventMonitorAdded, id, f)
}
func OnMonitorAddedv2(id string, f func(payload PayloadMonitorAddedV2)) {
	registerEvt(EventMonitorAddedV2, id, f)
}
func OnCreateWorkspace(id string, f func(payload PayloadCreateWorkspace)) {
	registerEvt(EventCreateWorkspace, id, f)
}
func OnCreateWorkspacev2(id string, f func(payload PayloadCreateWorkspaceV2)) {
	registerEvt(EventCreateWorkspaceV2, id, f)
}
func OnDestroyWorkspace(id string, f func(payload PayloadDestroyWorkspace)) {
	registerEvt(EventDestroyWorkspace, id, f)
}
func OnDestroyWorkspacev2(id string, f func(payload PayloadDestroyWorkspaceV2)) {
	registerEvt(EventDestroyWorkspaceV2, id, f)
}
func OnMoveWorkspace(id string, f func(payload PayloadMoveWorkspace)) {
	registerEvt(EventMoveWorkspace, id, f)
}
func OnMoveWorkspacev2(id string, f func(payload PayloadMoveWorkspaceV2)) {
	registerEvt(EventMoveWorkspaceV2, id, f)
}
func OnRenameWorkspace(id string, f func(payload PayloadRenameWorkspace)) {
	registerEvt(EventRenameWorkspace, id, f)
}
func OnActiveSpecial(id string, f func(payload PayloadActiveSpecial)) {
	registerEvt(EventActiveSpecial, id, f)
}
func OnActiveLayout(id string, f func(payload PayloadActiveLayout)) {
	registerEvt(EventActiveLayout, id, f)
}
func OnOpenWindow(id string, f func(payload PayloadOpenWindow)) {
	registerEvt(EventOpenWindow, id, f)
}
func OnCloseWindow(id string, f func(payload PayloadCloseWindow)) {
	registerEvt(EventCloseWindow, id, f)
}
func OnMoveWindow(id string, f func(payload PayloadMoveWindow)) {
	registerEvt(EventMoveWindow, id, f)
}
func OnMoveWindowv2(id string, f func(payload PayloadMoveWindowV2)) {
	registerEvt(EventMoveWindowV2, id, f)
}
func OnOpenLayer(id string, f func(payload PayloadOpenLayer)) {
	registerEvt(EventOpenLayer, id, f)
}
func OnCloseLayer(id string, f func(payload PayloadCloseLayer)) {
	registerEvt(EventCloseLayer, id, f)
}
func OnSubMap(id string, f func(payload PayloadSubMap)) {
	registerEvt(EventSubMap, id, f)
}
func OnChangeFloatingMode(id string, f func(payload PayloadChangeFloatingMode)) {
	registerEvt(EventChangeFloatingMode, id, f)
}
func OnUrgent(id string, f func(payload PayloadUrgent)) {
	registerEvt(EventUrgent, id, f)
}
func OnMinimize(id string, f func(payload PayloadMinimize)) {
	registerEvt(EventMinimize, id, f)
}
func OnScreencast(id string, f func(payload PayloadScreencast)) {
	registerEvt(EventScreencast, id, f)
}
func OnWindowTitle(id string, f func(payload PayloadWindowTitle)) {
	registerEvt(EventWindowTitle, id, f)
}
func OnIgnoreGroupLock(id string, f func(payload PayloadIgnoreGroupLock)) {
	registerEvt(EventIgnoreGroupLock, id, f)
}
func OnLockGroups(id string, f func(payload PayloadLockGroups)) {
	registerEvt(EventLockGroups, id, f)
}
func OnConfigReloaded(id string, f func()) {
	registerEvt(EventConfigReloaded, id, func(_ any) { f() })
}
func OnPin(id string, f func(payload PayloadPin)) {
	registerEvt(EventPin, id, f)
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
