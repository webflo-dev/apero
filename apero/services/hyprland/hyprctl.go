package hyprland

import "fmt"

type hyprlandCtlService struct {
	workspaces      []Workspace
	activeWorkspace Workspace
	clients         []Client
	activeClient    Client
}

var hyprlCtl = newHyprlandCtlService()

func newHyprlandCtlService() *hyprlandCtlService {
	service := &hyprlandCtlService{}
	service.syncWorkspaces()
	service.syncClients()
	return service
}

func (service *hyprlandCtlService) preprocessClientAddress(rawValues []string) []string {
	rawValues[0] = fmt.Sprintf("0x%s", rawValues[0])
	return rawValues
}

func (service *hyprlandCtlService) syncWorkspaces() {
	writeCmd("j/workspaces", &service.workspaces)
	service.syncActiveWorkspace()
}

func (service *hyprlandCtlService) syncActiveWorkspace() {
	writeCmd("j/activeworkspace", &service.activeWorkspace)
}

func (service *hyprlandCtlService) syncClients() {
	writeCmd("j/clients %s", &service.clients)
	service.syncActiveClient()
}

func (service *hyprlandCtlService) syncActiveClient() {
	writeCmd("j/activewindow", &service.activeClient)
}
