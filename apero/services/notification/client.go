package notification

import (
	gdb "webflo-dev/apero/dbus"

	"github.com/godbus/dbus/v5"
	// "github.com/godbus/dbus/v5"
)

const (
	path                = "/org/freedesktop/Notifications"
	iface               = "org.freedesktop.Notifications"
	methodNotify        = iface + ".Notify"
	methodClose         = iface + ".CloseNotification"
	signalActionInvoked = iface + ".ActionInvoked"
)

type Urgency uint16

const (
	UrgencyLow      Urgency = 0
	UrgencyCritical Urgency = 2
	UrgencyNormal   Urgency = 1
)

type Timeout int32

const (
	ExpiresDefault Timeout = -1
	ExpiresNever   Timeout = 0
)

type ActionRunner[T any] func(handle T)
type Actions[T any] map[string]ActionRunner[T]

func (a Actions[T]) getKeys() []string {
	keys := make([]string, 0, len(a))
	for key := range a {
		keys = append(keys, key)
	}
	return keys
}

type Action struct {
	key  string
	text string
}

type Notification[T any] struct {
	id        uint32
	replaceId uint32
	appName   string
	icon      string
	summary   string
	body      string
	urgency   Urgency
	category  string
	timeout   Timeout
	actions   Actions[T]
	resident  map[string]bool
	handle    T
}

func NewNotification[T any](handle T, summary string, body string) Notification[T] {
	return Notification[T]{
		summary:  summary,
		body:     body,
		urgency:  UrgencyNormal,
		actions:  make(Actions[T]),
		handle:   handle,
		timeout:  ExpiresDefault,
		resident: make(map[string]bool),
	}
}

func (n *Notification[T]) WithSummary(summary string) *Notification[T] {
	n.summary = summary
	return n
}

func (n *Notification[T]) WithBody(body string) *Notification[T] {
	n.body = body
	return n
}

func (n *Notification[T]) WithAppName(appName string) *Notification[T] {
	n.appName = appName
	return n
}

func (n *Notification[T]) WithIcon(icon string) *Notification[T] {
	n.icon = icon
	return n
}

func (n *Notification[T]) WithUrgency(urgency Urgency) *Notification[T] {
	n.urgency = urgency
	return n
}

func (n *Notification[T]) WithCategory(category string) *Notification[T] {
	n.category = category
	return n
}

func (n *Notification[T]) Replace(notificationId uint32) *Notification[T] {
	n.replaceId = notificationId
	return n
}

func (n *Notification[T]) WithTimeout(timeout Timeout) *Notification[T] {
	n.timeout = timeout
	return n
}

func (n *Notification[T]) WithAction(actionKey string, resident bool, action ActionRunner[T]) *Notification[T] {
	n.actions[actionKey] = action
	n.resident[actionKey] = resident
	return n
}

func (n *Notification[T]) Notify() (uint32, bool) {
	bus, err := gdb.ConnectToSessionBus(logger)
	if err != nil {
		return 0, false
	}

	var id uint32
	if ok := bus.Call(iface, path, methodNotify, &id, n.appName,
		n.replaceId,
		n.icon,
		n.summary,
		n.body,
		n.actions.getKeys(),
		make(map[string]any),
		n.timeout); !ok {
		return 0, false
	}

	n.id = id

	if len(n.actions) != 0 {
		return id, bus.WatchSignal(n.waitForAction, dbus.WithMatchObjectPath(path), dbus.WithMatchInterface(iface))
	}

	return id, true
}

func (n *Notification[T]) Close() {
	if bus, err := gdb.ConnectToSessionBus(logger); err == nil {
		bus.Call(iface, path, methodClose, 0, n.id)
	}
}

func (n *Notification[T]) waitForAction(signal *dbus.Signal) gdb.WatchBehavior {
	_id := signal.Body[0].(uint32)
	_key := signal.Body[1].(string)

	if signal.Path == path && signal.Name == signalActionInvoked && _id == n.id {
		if handle, ok := n.actions[_key]; ok {
			handle(n.handle)
			if n.resident[_key] == false {
				n.Close()
			}
			return gdb.WatchBehaviorStop
		}
	}

	return gdb.WatchBehaviorContinue
}
