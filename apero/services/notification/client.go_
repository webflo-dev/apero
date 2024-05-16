package notification

import (
	"webflo-dev/apero/services/notifications"
)

type Urgency uint16

const (
	UrgencyLow      Urgency = 0
	UrgencyCritical Urgency = 2
	UrgencyNormal   Urgency = 1
)

type ActionRunner[T any] func(T)
type Actions[T any] map[string]ActionRunner[T]

func (a Actions[T]) getKeys() []string {
	keys := make([]string, 0, len(a))
	for key := range a {
		keys = append(keys, key)
	}
	return keys
}

type notificationEventHandler[T any] struct {
	notifications.Subscriber
	id           uint32
	notification Notification[T]
}

type Notification[T any] struct {
	replaceId uint32
	appName   string
	icon      string
	summary   string
	body      string
	urgency   Urgency
	category  string
	actions   Actions[T]
	handle    T
}

func NewNotification[T any](handle T, summary string, body string) Notification[T] {
	return Notification[T]{
		summary: summary,
		body:    body,
		urgency: UrgencyNormal,
		actions: make(Actions[T]),
		handle:  handle,
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

func (n *Notification[T]) WithAction(actionKey string, action ActionRunner[T]) *Notification[T] {
	n.actions[actionKey] = action
	return n
}

func Notify[T any](n Notification[T]) uint32 {
	id := notifications.Notify(n.appName, n.replaceId, n.icon, n.summary, n.body, n.actions.getKeys(), nil, 0)

	handle := &notificationEventHandler[T]{
		id:           id,
		notification: n,
	}

	notifications.Register(handle, notifications.EventActionInvoked)

	return id
}

func (n *notificationEventHandler[T]) ActionInvoked(notificationId uint32, actionKey string) {
	if n.id != notificationId {
		return
	}

	if action, ok := n.notification.actions[actionKey]; ok {
		action(n.notification.handle)
		notifications.Unregister(n, notifications.EventActionInvoked)
	}
}
