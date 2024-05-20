package systray

import (
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/prop"
)

type ItemList map[string]*SystrayItem

type dbusServer struct {
	started  bool
	conn     *dbus.Conn
	props    *prop.Properties
	observer ServerObserver
	items    ItemList
}

type ServerObserver interface {
	NewItem(name string, item *SystrayItem)
	ItemRemoved(name string)
	ItemUpdated(name string)
}

type SystrayServer interface {
	Start() error
	Stop() error
	IsStarted() bool
	GetItems() ItemList
}

func StartNewServer(observer ServerObserver) (s SystrayServer, err error) {
	s = &dbusServer{
		started:  false,
		items:    make(ItemList),
		observer: observer,
	}

	err = s.Start()

	return
}

func (s *dbusServer) Start() (err error) {
	if s.started {
		return nil
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

func (s *dbusServer) GetItems() ItemList {
	return s.items
}
