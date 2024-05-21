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

func Register[T Subscriber](handle T, events ...EventType) {
	for _, event := range events {
		_service.subscribers[event] = append(_service.subscribers[event], handle)
	}
}

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
