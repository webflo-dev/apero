package notifications

type EventType string

const (
	EventNewNotification      EventType = "NewNotification"
	EventNotificationClosed   EventType = "NotificationClosed"
	EventDoNotDisturbChanged  EventType = "DoNotDisturbchanged"
	EventNotificationsChanged EventType = "NotificationsChanged"
)

type Subscriber interface {
	NewNotification(notification Notification)
	NotificationClosed(id uint32, reason uint32)
	DoNotDisturbChanged(enabled bool)
	NotificationsChanged(isEmpty bool)
}

type service struct {
	doNotDisturb bool
	subscribers  map[EventType][]Subscriber
	isClearing   bool
}

var _server NotificationServer
var _service = newService()

func newService() *service {
	return &service{
		doNotDisturb: false,
		subscribers:  make(map[EventType][]Subscriber),
		isClearing:   false,
	}
}

func (s *service) NewNotification(notification Notification) {
	if s.doNotDisturb == false {
		for _, handler := range s.subscribers[EventNewNotification] {
			handler.NewNotification(notification)
		}
	}
}

func (s *service) NotificationClosed(id uint32, reason uint32) {
	if s.doNotDisturb == false && s.isClearing == false {
		for _, subscriber := range s.subscribers[EventNotificationClosed] {
			subscriber.NotificationClosed(id, reason)
		}
	}
}

func (s *service) NotificationsChanged(isEmpty bool) {
	if s.isClearing == false {
		for _, subscriber := range s.subscribers[EventNotificationsChanged] {
			subscriber.NotificationsChanged(isEmpty)
		}
	}
}

func StartService() {
	_server, _ = StartNewServer(_service)
}

func Register[T Subscriber](handle T, events ...EventType) {
	for _, event := range events {
		_service.subscribers[event] = append(_service.subscribers[event], handle)
	}
}

func Unregister(handle Subscriber, events ...EventType) {
	for _, event := range events {
		handlers := _service.subscribers[event]
		for i, h := range handlers {
			if h == handle {
				_service.subscribers[event] = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}
	}
}

func SetDoNotDisturb(enabled bool) {
	if _service.doNotDisturb == enabled {
		return
	}

	_service.doNotDisturb = enabled

	for _, handler := range _service.subscribers[EventDoNotDisturbChanged] {
		handler.DoNotDisturbChanged(enabled)
	}
}

func DoNotDisturb() bool {
	return _service.doNotDisturb
}

func GetAllNotifications() NotificationList {
	return _server.GetAllNotifications()
}

func GetNotification(id uint32) (n Notification, ok bool) {
	n, ok = _server.GetNotification(id)
	return
}

func HasNotifications() bool {
	return _server.HasNotifications()
}

func ClearAllNotifications(notifyEach bool) {
	if notifyEach {
		_service.isClearing = true
	}
	for id := range _server.GetAllNotifications() {
		_server.CloseNotification(id)
	}
	if notifyEach {
		_service.isClearing = false
	}
}

func InvokeAction(id uint32, key string) (err error) {
	err = _server.InvokeAction(id, key)
	return
}

func CloseNotification(id uint32) {
	_server.CloseNotification(id)
	return
}
