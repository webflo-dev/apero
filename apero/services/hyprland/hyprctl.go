package hyprland

type hyprlandCtlService struct {
	workspaces      []Workspace
	activeWorkspace Workspace
}

var hyprlCtl = newHyprlandCtlService()

func newHyprlandCtlService() *hyprlandCtlService {
	service := &hyprlandCtlService{}
	service.syncWorkspaces()
	return service
}

func (service *hyprlandCtlService) syncWorkspaces() {
	writeCmd("j/workspaces", &service.workspaces)
	writeCmd("j/activeworkspace", &service.activeWorkspace)
}
