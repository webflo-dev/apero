package hyprland

func (s *hyprlandService) syncWorkspaces() {
	writeCmd("j/workspaces", &s.workspaces)
	s.syncActiveWorkspace()
}

func (s *hyprlandService) syncActiveWorkspace() {
	writeCmd("j/activeworkspace", &s.activeWorkspace)
}

func (s *hyprlandService) syncClients() {
	writeCmd("j/clients", &s.clients)
	s.syncActiveClient()
}

func (s *hyprlandService) syncActiveClient() {
	writeCmd("j/activewindow", &s.activeClient)
}

// func (s *service) syncMonitors() {
// 	writeCmd("j/monitors", &s.monitors)
// 	for _, monitor := range s.monitors {
// 		if monitor.Focused {
// 			s.activeMonitor = monitor
// 			break
// 		}
// 	}
// }
