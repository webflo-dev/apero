package systray

import (
	gdb "webflo-dev/apero/dbus"
	"webflo-dev/apero/services"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
)

const (
	path  = "/StatusNotifierWatcher"
	iface = "org.kde.StatusNotifierWatcher"

	signalWatcherItemRegistered   = iface + ".StatusNotifierItemRegistered"
	signalWatcherItemUnregistered = iface + ".StatusNotifierItemUnregistered"
)

var sigArrayStr, _ = dbus.ParseSignature("as")

type ItemList map[string]*SystrayItem

type server struct {
	base services.Service
	bus  *gdb.Bus

	props *prop.Properties
	// observer ServerObserver
	items ItemList
}

type spec struct {
	*server
}

func newServer() *server {
	return &server{
		base:  services.NewService(),
		items: make(ItemList),
	}
}
func (s *server) start() bool {
	if bus, err := gdb.ConnectToSessionBus(logger); err != nil {
		return false
	} else {
		s.bus = bus
	}

	spec := &spec{
		server: s,
	}

	if !s.bus.ExportAll(spec, path, iface) {
		return false
	}

	if props, ok := s.bus.ExportProps(path, iface, map[string]*prop.Prop{
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
	}); ok {
		s.props = props
	} else {
		return false
	}

	if !s.bus.RequestName(iface) {
		return false
	}

	if !s.bus.ExportIntrospectable(path, introspectData) {
		return false
	}

	if !s.bus.WatchSignals(s.watchForRemovedItem, map[string][]dbus.MatchOption{
		"freedesktop": {
			dbus.WithMatchInterface("org.freedesktop.DBus"),
			dbus.WithMatchObjectPath("/org/freedesktop/DBus"),
		},
		"item": {
			dbus.WithMatchInterface("org.kde.StatusNotifierItem"),
		},
	}) {
		return false
	}

	return s.base.Start(s.bus.Close, nil)
}

func (s *spec) RegisterStatusNotifierItem(service string, sender dbus.Sender) {
	logger.Printf("[dbus] RegisterStatusNotifierItem: service(%s) sender(%s)\n", service, sender)

	senderName := string(sender)

	if _, ok := s.server.items[senderName]; ok {
		logger.Println("already registered:", senderName)
		return
	}

	item := newSysTrayItem(senderName, service)
	s.server.items[senderName] = item
	s.server.updatePropItems()

	s.server.bus.EmitSignal(path, signalWatcherItemRegistered, senderName)

	// go func() {
	// 	s.observer.NewItem(name, item)
	// }()

}

func (s *spec) RegisterStatusNotifierHost(service string) {
	logger.Println("[dbus] RegisterStatusNotifierHost:", service)
	// emit signal <StatusNotifierHostRegistered>
}

var introspectData = introspect.Interface{
	Name: iface,

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

func (s *server) watchForRemovedItem(v *dbus.Signal) gdb.WatchBehavior {
	switch v.Name {
	case "org.freedesktop.DBus.NameOwnerChanged":
		name := v.Body[0].(string)
		newOwner := v.Body[2].(string)
		if newOwner == "" {
			if item, ok := s.items[name]; ok {
				delete(s.items, name)

				item.unregister()

				s.bus.EmitSignal(path, signalWatcherItemUnregistered, name)

				// go func() {
				// 	s.observer.ItemRemoved(name)
				// }()
			}
		}
	}

	return gdb.WatchBehaviorContinue
}

func (s *server) updatePropItems() {
	items := make([]string, 0, len(s.items))
	for s := range s.items {
		items = append(items, string(s))
	}

	err := s.props.Set(iface, "RegisteredStatusNotifierItems", dbus.MakeVariantWithSignature(items, sigArrayStr))
	if err != nil {
		logger.Println("RegisterStatusNotifierItem:", err)
	}
}
