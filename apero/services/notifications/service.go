package notifications

import (
	"webflo-dev/apero/events"
)

type EventType uint16

const (
	EventNewNotification EventType = iota
	EventNotificationClosed
	EventDoNotDisturbChanged
)

type PayloadNotificationClosed struct {
	Id     uint32
	Reason CloseReason
}

type PayloadDoNotDisturb struct {
	DoNotDisturb bool
}

type PayloadEmpty struct {
}

type notificationService struct {
	doNotDisturb bool
	isClearing   bool

	server *server

	evtNewNotification      events.Event[Notification]
	evtNotificationClosed   events.Event[PayloadNotificationClosed]
	evtDoNotDisturbChanged  events.Event[PayloadDoNotDisturb]
	evtNotificationsChanged events.Event[PayloadEmpty]
}

func newService() *notificationService {
	s := &notificationService{
		doNotDisturb: false,
		isClearing:   false,

		evtNewNotification:      events.New[Notification](),
		evtNotificationClosed:   events.New[PayloadNotificationClosed](),
		evtDoNotDisturbChanged:  events.New[PayloadDoNotDisturb](),
		evtNotificationsChanged: events.New[PayloadEmpty](),
	}

	observer := &observer{service: s}
	s.server = newServer(observer)

	return s
}

type observer struct {
	service *notificationService
}

func (o *observer) NotificationsChanged() {
	o.service.evtNotificationsChanged.Publish(PayloadEmpty{})
}
func (o *observer) NewNotification(n Notification) {
	o.service.evtNewNotification.Publish(n)
}

func (o *observer) NotificationClosed(id uint32, reason uint32) {
	o.service.evtNotificationClosed.Publish(PayloadNotificationClosed{id, CloseReason(reason)})
}
