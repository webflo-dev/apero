package hyprland

func (s *service) syncWorkspaces() {
	writeCmd("j/workspaces", &s.workspaces)
	s.syncActiveWorkspace()
}

func (s *service) syncActiveWorkspace() {
	writeCmd("j/activeworkspace", &s.activeWorkspace)
}

func (s *service) syncClients() {
	writeCmd("j/clients", &s.clients)
	s.syncActiveClient()
}

func (s *service) syncActiveClient() {
	writeCmd("j/activewindow", &s.activeClient)
}

func (s *service) syncMonitors() {
	writeCmd("j/monitors", &s.monitors)
	for _, monitor := range s.monitors {
		if monitor.Focused {
			s.activeMonitor = monitor
			break
		}
	}
}
