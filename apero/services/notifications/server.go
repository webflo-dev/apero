package notifications

import (
	"errors"
	"time"

	gdb "webflo-dev/apero/dbus"
	"webflo-dev/apero/services"

	"github.com/godbus/dbus/v5/introspect"
)

type closeReason uint32

const (
	path                     = "/org/freedesktop/Notifications"
	iface                    = "org.freedesktop.Notifications"
	signalActionInvoked      = iface + ".ActionInvoked"
	signalNotificationClosed = iface + ".NotificationClosed"

	closeReasonExpired   closeReason = 1
	closeReasonDismissed closeReason = 2
	closeReasonClosed    closeReason = 3
	closeReasonUnknown   closeReason = 4

	defaultTimeout = 5000
)

var serverCapabilities = []string{
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

const (
	ServerInfoName        = "apero"
	ServerInfoVendor      = "webflo-dev"
	ServerInfoVersion     = "0.1"
	ServerInfoSpecVersion = "1.2"
)

type NotificationList map[uint32]Notification

type server struct {
	base services.Service
	bus  *gdb.Bus

	counter       uint32
	notifications NotificationList
	ongoing       map[uint32]chan uint32
}

type spec struct {
	server *server
}

func newServer() *server {
	return &server{
		base:          services.NewService(),
		notifications: make(NotificationList),
		ongoing:       make(map[uint32]chan uint32),
	}
}

func (s *server) start() (ok bool) {
	s.bus, _ = gdb.ConnectToSessionBus(logger)

	spec := &spec{server: s}

	if ok = s.bus.ExportAll(spec, path, iface); !ok {
		return false
	}

	if ok = s.bus.RequestName(iface); !ok {
		return false
	}

	if ok = s.bus.ExportIntrospectable(path, introspectDataNotification); !ok {
		return false
	}

	return s.base.Start(s.bus.Close, nil)
}

func (s *server) stop() {
	s.base.Stop()
}

func (s *spec) GetCapabilities() []string {
	return serverCapabilities
}

func (s *spec) GetServerInformation() (string, string, string, string) {
	return ServerInfoName, ServerInfoVendor, ServerInfoVersion, ServerInfoSpecVersion
}

func (s *spec) Notify(
	appName string,
	replacesId uint32,
	appIcon string,
	summary string,
	body string,
	actions []string,
	hints hints,
	expireTimeout int,
) uint32 {
	var id uint32

	if replacesId == 0 {
		id = s.server.newId()
	} else {
		if replacedNotif, found := s.server.notifications[replacesId]; found {
			id = replacedNotif.id
		} else {
			id = s.server.newId()
		}
	}

	n := newNotification(id, appName, appIcon, summary, body, actions, hints)

	s.server.notifications[n.id] = n

	// logger.Printf("Notification > %+v\n", n)

	// go func() {
	// 	s.observer.NewNotification(n)
	// 	s.observer.NotificationsChanged(!s.HasNotifications())
	// }()

	if expireTimeout == -1 {
		expireTimeout = defaultTimeout
	}

	if expireTimeout > 0 {
		flag := make(chan uint32, 1)
		s.server.ongoing[n.id] = flag
		go s.server.waitForClose(int32(expireTimeout), n.id, uint32(closeReasonExpired), flag)
	}

	return n.id
}

func (s *spec) CloseNotification(id uint32) {
	s.server.closeNotification(id)
}

func (s *server) invokeAction(id uint32, key string) bool {
	if s.base.IsStarted() == false {
		logger.Println("Notification server is not started. Signal <ActionInvoked> not emitted.")
		return false
	}

	return s.bus.EmitSignal(path, signalActionInvoked, id, key)
}

func (s *server) newId() uint32 {
	s.counter++
	return s.counter
}

func (s *server) waitForClose(timeout int32, id uint32, reason uint32, flag chan uint32) {
	timer := time.NewTimer(time.Duration(timeout) * time.Millisecond)
	defer timer.Stop()

	select {
	case <-timer.C:
		s.removeNotification(id, uint32(closeReasonExpired))
	case reason = <-flag:
		s.removeNotification(id, uint32(reason))
	}
	delete(s.ongoing, id)
}

func (s *server) closeNotification(id uint32) {
	reason := uint32(closeReasonClosed)

	if n, found := s.ongoing[id]; found {
		n <- reason
	} else {
		s.removeNotification(id, reason)
	}
}

func (s *server) removeNotification(id uint32, reason uint32) (err error) {
	_, exists := s.notifications[id]
	isEmptyPrev := len(s.notifications) == 0

	delete(s.notifications, id)

	if (s.base.IsStarted()) == false {
		logger.Println("Notification server is not started. Signal <NotificationClosed> not emitted.")
		return errors.New("Notification server is not started. Signal <NotificationClosed> not emitted.")
	}

	s.bus.EmitSignal(path, signalNotificationClosed, id, reason)

	isEmptyNow := len(s.notifications) == 0

	if exists {
		// go func() { s.observer.NotificationClosed(id, reason) }()
	}

	if isEmptyPrev != isEmptyNow {
		// go func() { s.observer.NotificationsChanged(!isEmptyNow) }()
	}

	return err
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
			{Name: "id", Type: "u", Direction: "out"},
			{Name: "reason", Type: "u", Direction: "out"},
		}},

		{Name: "ActionInvoked", Args: []introspect.Arg{
			{Name: "id", Type: "u", Direction: "out"},
			{Name: "action_key", Type: "s", Direction: "out"},
		}},
	},
}
