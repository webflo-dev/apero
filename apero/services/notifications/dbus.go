package notifications

import (
	"errors"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
)

type Reason uint32

const (
	path                     = "/org/freedesktop/Notifications"
	iface                    = "org.freedesktop.Notifications"
	signalActionInvoked      = iface + ".ActionInvoked"
	signalNotificationClosed = iface + ".NotificationClosed"

	reasonExpired   Reason = 1
	reasonDismissed Reason = 2
	reasonClosed    Reason = 3
	reasonUnknown   Reason = 4

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

type NotificationList map[uint32]Notification

type NotificationServer interface {
	Start() error
	Stop() error
	IsStarted() bool
	GetCapabilities() []string
	GetServerInformation() (string, string, string, string)
	GetAllNotifications() NotificationList
	GetNotification(id uint32) (Notification, bool)
	HasNotifications() bool
	CloseNotification(id uint32)
	InvokeAction(id uint32, key string) error
}

type ServerObserver interface {
	NewNotification(notification Notification)
	NotificationClosed(id uint32, reason uint32)
	NotificationsChanged(isEmpty bool)
}

type DBusServer interface {
	GetCapabilities() []string
	GetServerInformation() (string, string, string, string)
	Notify(appName string, replacesId uint32, appIcon string, summary string, body string, actions []string, hints hints, expireTimeout int) uint32
	CloseNotification(id uint32)
}

type dbusServer struct {
	DBusServer
	started       bool
	conn          *dbus.Conn
	counter       uint32
	notifications NotificationList
	observer      ServerObserver
	ongoing       map[uint32]chan uint32
}

func StartNewServer(observer ServerObserver) (s NotificationServer, err error) {
	s = &dbusServer{
		started:       false,
		observer:      observer,
		notifications: make(NotificationList),
		ongoing:       make(map[uint32]chan uint32),
	}
	s.Start()
	return
}

func (s *dbusServer) Start() (err error) {

	s.conn, err = dbus.ConnectSessionBus()
	if err != nil {
		logger.Println("Notification server is disabled (failed to connect to session bus)", err)
		return errors.New("Notification server is disabled (failed to connect to session bus)")
	}

	err = s.conn.ExportAll(s, path, iface)
	if err != nil {
		logger.Println("Notification server is disabled (failed to export to dbus)", err)
		return errors.New("Notification server is disabled (failed to export to dbus)")
	}

	reply, err := s.conn.RequestName(iface, dbus.NameFlagDoNotQueue)
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
	err = s.conn.Export(introspect.NewIntrospectable(&node), path, "org.freedesktop.DBus.Introspectable")
	if err != nil {
		logger.Println("Failed to export introspectable interface.", err)
	}

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

func (s *dbusServer) GetAllNotifications() NotificationList {
	return s.notifications
}

func (s *dbusServer) GetNotification(id uint32) (n Notification, ok bool) {
	n, ok = s.notifications[id]
	return
}

func (s *dbusServer) HasNotifications() bool {
	return len(s.notifications) > 0
}

func (s *dbusServer) InvokeAction(id uint32, key string) error {
	if (s.started) == false {
		logger.Println("Notification server is not started. Signal <ActionInvoked> not emitted.")
		return errors.New("Notification server is not started. Signal <ActionInvoked> not emitted.")
	}

	err := s.conn.Emit(path, signalActionInvoked, id, key)
	if err != nil {
		logger.Println("Failed to emit signal ActionInvoked (emit signal)", err)
		return errors.New("Failed to emit signal ActionInvoked (emit signal)")
	}

	return nil
}

func (s *dbusServer) GetCapabilities() []string {
	return serverCapabilities
}

func (s *dbusServer) GetServerInformation() (string, string, string, string) {
	return "apero", "webflo-dev", "0.1", "1.2"
}

func (s *dbusServer) Notify(
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
		id = s.newId()
	} else {
		if replacedNotif, found := s.notifications[replacesId]; found {
			id = replacedNotif.id
		} else {
			id = s.newId()
		}
	}

	n := newNotification(id, appName, appIcon, summary, body, actions, hints)

	s.notifications[n.id] = n

	// logger.Printf("Notification > %+v\n", n)

	go func() {
		s.observer.NewNotification(n)
		s.observer.NotificationsChanged(!s.HasNotifications())
	}()

	if expireTimeout == -1 {
		expireTimeout = defaultTimeout
	}

	if expireTimeout > 0 {
		flag := make(chan uint32, 1)
		s.ongoing[n.id] = flag
		go s.waitForClose(int32(expireTimeout), n.id, uint32(reasonExpired), flag)
	}

	return n.id
}

func (s *dbusServer) CloseNotification(id uint32) {
	reason := uint32(reasonClosed)

	if n, found := s.ongoing[id]; found {
		n <- reason
	} else {
		s.closeNotification(id, reason)
	}
}

func (s *dbusServer) newId() uint32 {
	s.counter++
	return s.counter
}

func (s *dbusServer) waitForClose(timeout int32, id uint32, reason uint32, flag chan uint32) {
	timer := time.NewTimer(time.Duration(timeout) * time.Millisecond)
	defer timer.Stop()

	select {
	case <-timer.C:
		s.closeNotification(id, uint32(reasonExpired))
	case reason = <-flag:
		s.closeNotification(id, uint32(reason))
	}
	delete(s.ongoing, id)
}

func (s *dbusServer) closeNotification(id uint32, reason uint32) (err error) {
	_, exists := s.notifications[id]
	isEmptyPrev := s.HasNotifications()

	delete(s.notifications, id)

	if (s.started) == false {
		logger.Println("Notification server is not started. Signal <NotificationClosed> not emitted.")
		return errors.New("Notification server is not started. Signal <NotificationClosed> not emitted.")
	}

	err = s.conn.Emit(path, signalNotificationClosed, id, reason)
	if err != nil {
		logger.Println("Failed to emit signal <NotificationClosed> (emit signal)", err)
		err = errors.New("Failed to emit signal <NotificationClosed> (emit signal)")
	}

	isEmptyNow := s.HasNotifications()

	if exists {
		go func() { s.observer.NotificationClosed(id, reason) }()
	}

	if isEmptyPrev != isEmptyNow {
		go func() { s.observer.NotificationsChanged(!isEmptyNow) }()
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
