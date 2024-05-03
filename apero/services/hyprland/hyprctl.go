package hyprland

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

func (service *hyprlandCtlService) syncWorkspaces() {
	writeCmd("j/workspaces", &service.workspaces)
	service.syncActiveWorkspace()
}

func (service *hyprlandCtlService) syncActiveWorkspace() {
	writeCmd("j/activeworkspace", &service.activeWorkspace)
}

func (service *hyprlandCtlService) syncClients() {
	writeCmd("j/clients", &service.clients)
	service.syncActiveClient()
}

func (service *hyprlandCtlService) syncActiveClient() {
	writeCmd("j/activewindow", &service.activeClient)
}
