package systray

import (
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/prop"
)

type dbusServer struct {
	started bool
	conn    *dbus.Conn
	props   *prop.Properties
	items   map[string]string
}

type SystrayServer interface {
	Start() error
	Stop() error
	IsStarted() bool
}

func StartNewServer() (s SystrayServer, err error) {
	s = &dbusServer{
		started: false,
		items:   make(map[string]string),
	}
	err = s.Start()
	return
}

func (s *dbusServer) Start() (err error) {
	if s.started {
		return nil
	}

	s.conn, err = dbus.ConnectSessionBus()
	if err != nil {
		logger.Println("Systray server is disabled (failed to connect to session bus)", err)
		return
	}

	s.registerWatcher()

	s.started = true

	go func() {
		defer s.conn.Close()

		for {
			if s.started == false {
				break
			}
		}
	}()

	return
}

func (s *dbusServer) Stop() error {
	s.started = false
	return nil
}

func (s *dbusServer) IsStarted() bool {
	return s.started
}
