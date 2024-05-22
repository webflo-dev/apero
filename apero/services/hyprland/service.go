package hyprland

import (
	"net"
	"strings"
	"webflo-dev/apero/services"
)

type hyprlandService struct {
	base services.Service
	conn net.Conn

	workspaces      []Workspace
	activeWorkspace Workspace
	clients         []Client
	activeClient    Client
}

func newService() *hyprlandService {
	service := &hyprlandService{
		base: services.NewService(),
	}

	return service
}

func (s *hyprlandService) stop() {
	s.base.Stop()
}

func (s *hyprlandService) close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *hyprlandService) start() {
	s.conn = createEventsConnection()

	s.base.Start(s.close, s.loop)
	logger.Println("listening for hyprland events")
}

func (s *hyprlandService) loop() services.LoopBehavior {
	msg, err := readEvent(s.conn)
	if err != nil {
		logger.Println("Error receiving message", err)
		return services.LoopBehaviorContinue
	}

	for _, data := range msg {
		// log.Printf("%+v\n", data)
		s.processEvent(data)
	}

	return services.LoopBehaviorContinue
}

func (s *hyprlandService) processEvent(msg EventMessage) {
	values := MsgValues(strings.Split(msg.Data, ","))

	switch msg.Type {
	case EventActiveWindow, EventActiveWindowV2:
		s.syncActiveClient()

	case
		EventFullscreen,
		EventOpenWindow,
		EventCloseWindow,
		EventMoveWindow,
		EventMoveWindowV2:
		s.syncWorkspaces()
		s.syncClients()

	case
		EventCreateWorkspace,
		EventCreateWorkspaceV2,
		EventDestroyWorkspace,
		EventDestroyWorkspaceV2,
		EventMoveWorkspace,
		EventMoveWorkspaceV2,
		EventRenameWorkspace:
		s.syncWorkspaces()

	case
		EventChangeFloatingMode,
		EventWindowTitle,
		EventPin:
		s.syncClients()
	}

	publishEvt(msg.Type, values)
}
