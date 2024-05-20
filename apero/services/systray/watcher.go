package systray

import (
	"errors"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
)

const (
	watcherPath  = "/StatusNotifierWatcher"
	watcherIface = "org.kde.StatusNotifierWatcher"

	signalWatcherItemRegistered   = watcherIface + ".StatusNotifierItemRegistered"
	signalWatcherItemUnregistered = watcherIface + ".StatusNotifierItemUnregistered"
)

var sigArrayStr, _ = dbus.ParseSignature("as")

func (s *dbusServer) registerWatcher() (err error) {
	s.conn, err = dbus.ConnectSessionBus()
	if err != nil {
		logger.Println("Systray server is disabled (failed to connect to session bus)", err)
		return
	}

	err = s.conn.ExportAll(s, watcherPath, watcherIface)
	if err != nil {
		logger.Println("Systray server is disabled (failed to export watcher to dbus)", err)
		return errors.New("Systray server is disabled (failed to export watcher to dbus)")
	}

	propsSpec := map[string]map[string]*prop.Prop{
		"org.kde.StatusNotifierWatcher": {
			"RegisteredStatusNotifierItems": {
				Value:    []string{},
				Writable: true,
				Emit:     prop.EmitFalse,
			},
			"IsStatusNotifierHostRegistered": {
				Value:    true,
				Writable: true,
				Emit:     prop.EmitConst,
			},
			"ProtocolVersion": {
				Value:    123,
				Writable: false,
				Emit:     prop.EmitConst,
			},
		},
	}

	s.props, err = prop.Export(s.conn, watcherPath, propsSpec)
	if err != nil {
		logger.Println("Systray server is disabled (failed to export props)", err)
		return errors.New("Systray server is disabled (failed to export props)")
	}

	reply, err := s.conn.RequestName(watcherIface, dbus.NameFlagDoNotQueue)
	if err != nil {
		logger.Println("Systray server is disabled (failed to request name for watcher on session dbus)", err)
		return errors.New("Systray server is disabled (failed to request name for watcher on session dbus)")
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		logger.Println("Systray server is disabled (watcher name already taken)", err)
		return errors.New("Systray server is disabled (watcher name already taken)")
	}

	node := &introspect.Node{
		Name: watcherPath,
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			prop.IntrospectData,
			introspectDataWatcher,
		},
	}
	err = s.conn.Export(introspect.NewIntrospectable(node), watcherPath, "org.freedesktop.DBus.Introspectable")
	if err != nil {
		logger.Println("Systray server is disabled (failed to be introspectable)", err)
		return errors.New("Systray server is disabled (failed to be introspectable)")
	}

	watchErr := s.conn.AddMatchSignal(dbus.WithMatchInterface("org.freedesktop.DBus"), dbus.WithMatchObjectPath("/org/freedesktop/DBus"))
	_ = s.conn.AddMatchSignal(dbus.WithMatchInterface("org.kde.StatusNotifierItem"))
	if watchErr != nil {
		logger.Println("Systray server is disabled (failed to monitor systray name loss)", err)
		return errors.New("Systray server is disabled (failed to monitor systray name loss)")
	}

	c := make(chan *dbus.Signal, 10)
	s.conn.Signal(c)

	go func() {
		for v := range c {
			switch v.Name {
			case "org.freedesktop.DBus.NameOwnerChanged":
				name := v.Body[0].(string)
				newOwner := v.Body[2].(string)
				if newOwner == "" {
					if _, ok := s.items[name]; ok {
						delete(s.items, name)

						err := s.conn.Emit(watcherPath, signalWatcherItemUnregistered)
						if err != nil {
							logger.Println("Failed to emit signal <StatusNotifierItemUnregistered>", err)
						}

						go func() {
							s.observer.ItemRemoved(name)
						}()

						logger.Println("[dbus] item lost:", name, s.items)
					}
				}
			}
		}
	}()

	return
}

func (s *dbusServer) RegisterStatusNotifierItem(service string, sender dbus.Sender) {
	logger.Printf("[dbus] RegisterStatusNotifierItem: service(%s) sender(%s)\n", service, sender)

	name := string(sender)

	if _, ok := s.items[name]; ok {
		logger.Println("[dbus] already registered:", name)
		return
	}

	item := newSysTrayItem(name, service)
	item.Register(s)
	s.items[name] = item

	err := s.conn.Emit(watcherPath, signalWatcherItemRegistered)
	if err != nil {
		logger.Println("Failed to emit signal <StatusNotifierItemRegistered>", err)
	}

	s.updatePropItems()

	go func() {
		s.observer.NewItem(name, item)
	}()

}

func (s *dbusServer) RegisterStatusNotifierHost(service string) {
	logger.Println("[dbus] RegisterStatusNotifierHost:", service)
	// emit signal <StatusNotifierHostRegistered>
}

func (s *dbusServer) updatePropItems() {
	items := make([]string, 0, len(s.items))
	for s := range s.items {
		items = append(items, string(s))
	}

	err := s.props.Set(watcherIface, "RegisteredStatusNotifierItems", dbus.MakeVariantWithSignature(items, sigArrayStr))
	if err != nil {
		logger.Println("[dbus] RegisterStatusNotifierItem:", err)
	}
}

func (s *dbusServer) ItemUpdated(name string) {
	s.observer.ItemUpdated(name)
}

var introspectDataWatcher = introspect.Interface{
	Name: watcherIface,

	Methods: []introspect.Method{
		{Name: "RegisterStatusNotifierItem", Args: []introspect.Arg{
			{Name: "service", Type: "s", Direction: "in"},
		}},
		{Name: "RegisterStatusNotifierHost", Args: []introspect.Arg{
			{Name: "service", Type: "s", Direction: "in"},
		}},
	},

	Properties: []introspect.Property{
		{Name: "RegisteredStatusNotifierItems", Type: "as", Access: "read"},
		{Name: "IsStatusNotifierHostRegistered", Type: "b", Access: "read"},
		{Name: "ProtocolVersion", Type: "i", Access: "read"},
	},

	Signals: []introspect.Signal{
		{Name: "StatusNotifierItemRegistered", Args: []introspect.Arg{
			{Name: "service", Type: "s", Direction: "out"},
		}},
		{Name: "StatusNotifierItemUnregistered", Args: []introspect.Arg{
			{Name: "service", Type: "s", Direction: "out"},
		}},
		{Name: "StatusNotifierHostRegistered"},
		{Name: "StatusNotifierHostUnregistered"},
	},
}
