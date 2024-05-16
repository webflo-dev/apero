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
)

func (s *dbusServer) registerWatcher() (err error) {
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

	return
}

func (s *dbusServer) RegisterStatusNotifierItem(service string) {
	s.items[service] = service

	var services = []string{}
	for s := range s.items {
		services = append(services, s)
	}

	err := s.props.Set(watcherIface, "RegisteredStatusNotifierItems", dbus.MakeVariant(services))
	if err != nil {
		logger.Println("[dbus] RegisterStatusNotifierItem:", err)
	}
}

func (s *dbusServer) RegisterStatusNotifierHost(service string) {
	logger.Println("[dbus] RegisterStatusNotifierHost:", service)
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
