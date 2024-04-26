package notifications

import (
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
)

const (
	path  = "/org/freedesktop/Notifications"
	iface = "org.freedesktop.Notifications"
)

func StartNotificationServer() {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		Logger.Println("Failed to connect to session bus. Notification server is now disabled", err)
		return
	}

	server := newNotificationServer()

	err = conn.ExportAll(server, path, iface)
	if err != nil {
		Logger.Println("Failed to export notification server. Notification server is now disabled", err)
		return
	}

	reply, err := conn.RequestName("org.freedesktop.Notifications", dbus.NameFlagDoNotQueue)
	if err != nil {
		Logger.Println("Failed to request name on session bus. Notification server is now disabled", err)
		return
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		Logger.Println("Name already taken. Notification server is now disabled", err)
		return
	}

	node := introspect.Node{
		Name: path,
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			prop.IntrospectData,
			IntrospectDataNotification,
		},
	}
	err = conn.Export(introspect.NewIntrospectable(&node), path, "org.freedesktop.DBus.Introspectable")
	if err != nil {
		Logger.Println("Failed to export introspectable interface. Notification server is now disabled", err)
		return
	}

	go func() {
		defer conn.Close()

		Logger.Println("Listening on org.freedesktop.Notifications / /org/freedesktop/Notifications ...")

		select {}
	}()
}

var IntrospectDataNotification = introspect.Interface{
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
