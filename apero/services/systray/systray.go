package systray

import (
	"log"
	"webflo-dev/apero/logger"
	"webflo-dev/apero/services/systray/notifier"
	"webflo-dev/apero/services/systray/watcher"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
)

const (
	path     = "/StatusNotifierWatcher"
	hostPath = "/StatusNotifierHost"
)

type systray struct {
	conn *dbus.Conn
	// menu  *menu.Dbusmenu
	nodes map[dbus.Sender]*node
}

type node struct {
	ni *notifier.StatusNotifierItem
}

func NewSystrayService() *systray {
	t := &systray{nodes: make(map[dbus.Sender]*node)}

	conn, _ := dbus.ConnectSessionBus()
	t.conn = conn

	err := conn.ExportAll(struct{}{}, hostPath, "org.kde.StatusNotifierHost")
	if err != nil {
		logger.AppLogger.Println("Err", err)
		return t
	}

	// TODO this is create watcher (optional)
	err = conn.ExportAll(t, path, "org.kde.StatusNotifierWatcher")
	if err != nil {
		log.Println("Err2", err)
		return t
	}

	_, err = conn.RequestName("org.kde.StatusNotifierWatcher", dbus.NameFlagDoNotQueue)
	if err != nil {
		log.Println("Failed to claim notifier watcher name", err)
		return t
	}

	_, err = prop.Export(conn, path, createPropSpec())
	if err != nil {
		log.Printf("Failed to export notifier item properties to bus")
		return t
	}

	node := introspect.Node{
		Name: path,
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			prop.IntrospectData,
			IntrospectDataStatusNotifierWatcher,
		},
	}
	err = conn.Export(introspect.NewIntrospectable(&node), path,
		"org.freedesktop.DBus.Introspectable")
	if err != nil {
		log.Printf("Failed to export introspection %v", err)
		return t
	}
	// End TODO

	hostErr := t.RegisterStatusNotifierHost(conn.Names()[0], "")
	if hostErr != nil {
		logger.AppLogger.Println("Failed to register our host with the notifier watcher, maybe no watcher running?", hostErr)
		return t
	}

	watchErr := t.conn.AddMatchSignal(dbus.WithMatchInterface("org.freedesktop.DBus"), dbus.WithMatchObjectPath("/org/freedesktop/DBus"))
	_ = t.conn.AddMatchSignal(dbus.WithMatchInterface("org.kde.StatusNotifierItem"))
	if watchErr != nil {
		logger.AppLogger.Println("Failed to monitor systray name loss", watchErr)
		return t
	}

	c := make(chan *dbus.Signal, 10)
	t.conn.Signal(c)
	go func() {
		for v := range c {
			log.Println("Signal", v.Name, v.Body)
			switch v.Name {
			case "org.freedesktop.DBus.NameOwnerChanged":
				name := v.Body[0]
				newOwner := v.Body[2]
				if newOwner == "" {
					if item, ok := t.nodes[dbus.Sender(name.(string))]; ok {
						logger.AppLogger.Println("Remove item", item)
						// remove the item
					}
				}
			case "org.kde.StatusNotifierItem.NewIcon":
				item, ok := t.nodes[dbus.Sender(v.Sender)]
				if ok {
					logger.AppLogger.Println("Update icon of item", item)
					// update icon of the item
				}
			default:
				logger.AppLogger.Println("Also", v.Name)
				continue
			}
		}
	}()

	return t
}

func (t *systray) RegisterStatusNotifierHost(service string, sender dbus.Sender) (err *dbus.Error) {
	logger.AppLogger.Println("Register Host", service, sender)

	e := watcher.Emit(t.conn, &watcher.StatusNotifierWatcher_StatusNotifierHostRegisteredSignal{
		Path: dbus.ObjectPath(service),
		Body: &watcher.StatusNotifierWatcher_StatusNotifierHostRegisteredSignalBody{},
	})
	if e != nil {
		logger.AppLogger.Println("it was not emit the notification ", err)
	}
	return nil
}

func createPropSpec() map[string]map[string]*prop.Prop {
	return map[string]map[string]*prop.Prop{
		"org.kde.StatusNotifierWatcher": {
			"RegisteredStatusNotifierItems": {
				Value:    []string{},
				Writable: false,
				Emit:     prop.EmitTrue,
				Callback: nil,
			},
			"IsStatusNotifierHostRegistered": {
				Value:    true,
				Writable: false,
				Emit:     prop.EmitTrue,
				Callback: nil,
			},
			"ProtocolVersion": {
				Value:    int32(25),
				Writable: false,
				Emit:     prop.EmitTrue,
				Callback: nil,
			},
		},
	}
}

var (
	// Introspection for org.kde.StatusNotifierWatcher
	IntrospectDataStatusNotifierWatcher = introspect.Interface{
		Name: "org.kde.StatusNotifierWatcher",
		Methods: []introspect.Method{{Name: "RegisterStatusNotifierItem", Args: []introspect.Arg{
			{Name: "service", Type: "s", Direction: "in"},
		}},
			{Name: "RegisterStatusNotifierHost", Args: []introspect.Arg{
				{Name: "service", Type: "s", Direction: "in"},
			}},
		},
		Signals: []introspect.Signal{{Name: "StatusNotifierItemRegistered", Args: []introspect.Arg{
			{Name: "", Type: "s", Direction: ""},
		}},
			{Name: "StatusNotifierItemUnregistered", Args: []introspect.Arg{
				{Name: "", Type: "s", Direction: ""},
			}},
			{Name: "StatusNotifierHostRegistered"},
			{Name: "StatusNotifierHostUnregistered"},
		},
		Properties: []introspect.Property{{Name: "RegisteredStatusNotifierItems", Type: "as", Access: "read", Annotations: []introspect.Annotation{
			{Name: "org.qtproject.QtDBus.QtTypeName.Out0", Value: "QStringList"},
		}},
			{Name: "IsStatusNotifierHostRegistered", Type: "b", Access: "read"},
			{Name: "ProtocolVersion", Type: "i", Access: "read"},
		},
		Annotations: []introspect.Annotation{},
	}
)
