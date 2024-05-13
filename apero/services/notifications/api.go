package notifications

type EventType string

const (
	EventNewNotification      EventType = "NewNotification"
	EventNotificationRemoved  EventType = "NotificationRemoved"
	EventNotificationsCleared EventType = "NotificationsCleared"
	EventDoNotDisturbChanged  EventType = "DoNotDisturbchanged"
	EventActionInvoked        EventType = "ActionInvoked"
)

type Subscriber interface {
	NewNotification(notification Notification)
	NotificationRemoved(id uint32)
	NotificationsCleared()
	DoNotDisturbChanged(enabled bool)
	ActionInvoked(notificationId uint32, actionKey string)
}

func GetServerCapabilities() []string {
	return _service.GetCapabilities()
}

func GetServerInformation() (string, string, string, string) {
	return _service.GetServerInformation()
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

func ClearAllNotifications(notifyEach bool) {
	if notifyEach == false {
		clear(_service.notifications)
		for _, handler := range _service.subscribers[EventNotificationsCleared] {
			handler.NotificationsCleared()
		}
		return
	}

	ids := make([]uint32, 0, len(_service.notifications))
	for id := range _service.notifications {
		ids = append(ids, id)
	}
	for id := range ids {
		_service.CloseNotification(uint32(id))
	}
}

func GetNotifications() []Notification {
	notifications := make([]Notification, 0, len(_service.notifications))
	for _, n := range _service.notifications {
		notifications = append(notifications, n)
	}
	return notifications
}

func GetNotification(id uint32) (Notification, bool) {
	n, ok := _service.notifications[id]
	return n, ok
}

func HasNotifications() bool {
	return len(_service.notifications) > 0
}

func DoNotDisturb() bool {
	return _service.doNotDisturb
}

func Notify(appName string, replacesId uint32, appIcon string, summary string, body string, actions []string, hints hints, expireTimeout int) uint32 {
	return _service.Notify(appName, replacesId, appIcon, summary, body, actions, hints, expireTimeout)
}

func InvokeAction(notificationId uint32, actionKey string) {
	if n, ok := _service.notifications[notificationId]; ok {
		_service.InvokeAction(n.id, actionKey)
	}
}
