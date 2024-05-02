package notifications

import (
	"errors"
	"log"
	"reflect"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
)

const (
	path                = "/org/freedesktop/Notifications"
	iface               = "org.freedesktop.Notifications"
	signalActionInvoked = iface + ".ActionInvoked"
)

type notificationServer struct {
	counter       uint32
	notifications map[uint32]Notification
	started       bool
	subscribers   []NotificationsEventHandler
	doNotDisturb  bool
}

var server = newNotificationServer()

func newNotificationServer() *notificationServer {
	return &notificationServer{
		counter:       0,
		started:       false,
		doNotDisturb:  false,
		notifications: make(map[uint32]Notification),
	}
}

func StartNotificationServer() error {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		logger.Println("Notification server is disabled (failed to connect to session bus)", err)
		return errors.New("Notification server is disabled (failed to connect to session bus)")
	}

	err = conn.ExportAll(server, path, iface)
	if err != nil {
		logger.Println("Notification server is disabled (failed to export to dbus)", err)
		return errors.New("Notification server is disabled (failed to export to dbus)")
	}

	reply, err := conn.RequestName(iface, dbus.NameFlagDoNotQueue)
	if err != nil {
		logger.Println("Notification server is disabled (failed to request name on session dbus)", err)
		return errors.New("Notification server is disabled (failed to request name on session dbus)")
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		logger.Println("Notification server is disabled (name already taken)", err)
		return errors.New("Notification server is disabled (name already taken)")
	}

	node := introspect.Node{
		Name: path,
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			prop.IntrospectData,
			introspectDataNotification,
		},
	}
	err = conn.Export(introspect.NewIntrospectable(&node), path, "org.freedesktop.DBus.Introspectable")
	if err != nil {
		logger.Println("Failed to export introspectable interface.", err)
	}

	go func() {
		defer conn.Close()

		server.started = true
		logger.Printf("Listening on iface=%s, path=%s ...\n", iface, path)

		select {}
	}()

	return nil
}

func (server notificationServer) GetCapabilities() []string {
	return []string{
		"action-icons",
		"actions",
		"body",
		"body-hyperlinks",
		"body-images",
		"body-markup",
		"icon-multi",
		"icon-static",
		"persistence",
		"sound",
	}
}

func (server notificationServer) GetServerInformation() (string, string, string, string) {
	return "apero", "webflo-dev", "0.1", "1.2"
}

func (server notificationServer) Notify(appName string, replacesId uint32, appIcon string, summary string, body string, actions []string, hints hints, expireTimeout int) uint32 {

	id := replacesId
	if id == 0 {
		server.counter++
		id = server.counter
	}

	n := newNotification(id, appName, appIcon, summary, body, actions, hints)

	server.notifications[n.id] = n

	// logger.Printf("Notification > %+v\n", n)

	if server.doNotDisturb == false {
		for _, subscriber := range server.subscribers {
			subscriber.NewNotification(n)
		}
	}

	return uint32(n.id)
}

func (server notificationServer) CloseNotification(id uint32) {
	delete(server.notifications, id)

	if server.doNotDisturb == false {
		for _, subscriber := range server.subscribers {
			subscriber.NotificationRemoved(id)
		}
	}
}

func (server notificationServer) InvokeAction(notificationId uint32, actionKey string) error {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		logger.Println("Failed to emit signal ActionInvoked (dbus connection)", err)
		return errors.New("Failed to emit signal ActionInvoked (dbus connection)")
	}
	err = conn.Emit(path, signalActionInvoked, notificationId, actionKey)
	if err != nil {
		logger.Println("Failed to emit signal ActionInvoked (emit signal)", err)
		return errors.New("Failed to emit signal ActionInvoked (emit signal)")
	}

	for _, subscriber := range server.subscribers {
		log.Printf("ActionInvoked > %+v\n", subscriber)
		if callback, ok := reflect.TypeOf(subscriber).MethodByName("ActionInvoked"); ok {
			log.Printf("ActionInvoked method found > %+v\n", callback)
			subscriber.ActionInvoked(notificationId, actionKey)
		}
	}
	return nil
}

var introspectDataNotification = introspect.Interface{
	Name: "org.freedesktop.Notifications",

	Methods: []introspect.Method{
		{Name: "GetServerInformation", Args: []introspect.Arg{
			{Name: "name", Type: "s", Direction: "out"},
			{Name: "vendor", Type: "s", Direction: "out"},
			{Name: "version", Type: "s", Direction: "out"},
			{Name: "spec_version", Type: "s", Direction: "out"},
		}},

		{Name: "GetCapabilities", Args: []introspect.Arg{
			{Name: "capabilities", Type: "as", Direction: "out"},
		}},

		{Name: "Notify", Args: []introspect.Arg{
			{Name: "app_name", Type: "s", Direction: "in"},
			{Name: "replaces_id", Type: "u", Direction: "in"},
			{Name: "app_icon", Type: "s", Direction: "in"},
			{Name: "summary", Type: "s", Direction: "in"},
			{Name: "body", Type: "s", Direction: "in"},
			{Name: "actions", Type: "as", Direction: "in"},
			{Name: "hints", Type: "a{sv}", Direction: "in"},
			{Name: "expire_timeout", Type: "i", Direction: "in"},
			{Name: "id", Type: "u", Direction: "out"},
		}},

		{Name: "CloseNotification", Args: []introspect.Arg{
			{Name: "id", Direction: "in", Type: "u"},
		}},
	},

	Signals: []introspect.Signal{
		{Name: "NotificationClosed", Args: []introspect.Arg{
			{Name: "id", Type: "u", Direction: ""},
			{Name: "reason", Type: "u", Direction: ""},
		}},

		{Name: "ActionInvoked", Args: []introspect.Arg{
			{Name: "id", Type: "u", Direction: ""},
			{Name: "action_key", Type: "s", Direction: ""},
		}},
	},
}
