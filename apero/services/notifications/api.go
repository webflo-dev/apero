package notifications

import "webflo-dev/apero/events"

var _service = newService()

func StartService() {
	_service.server.start()
}

func OnNewNotification(f func(payload Notification)) (events.ID, error) {
	return _service.evtNewNotification.RegisterHandler(events.HandlerFunc[Notification](f))
}
func OnNotificationClosed(id string, f func(payload PayloadNotificationClosed)) (events.ID, error) {
	return _service.evtNotificationClosed.RegisterHandler(events.HandlerFunc[PayloadNotificationClosed](f))
}
func OnNotificationsChanged(id string, f func(payload PayloadEmpty)) (events.ID, error) {
	return _service.evtNotificationsChanged.RegisterHandler(events.HandlerFunc[PayloadEmpty](f))
}

func SetDoNotDisturb(enabled bool) {
	if _service.doNotDisturb == enabled {
		return
	}

	_service.doNotDisturb = enabled

	_service.evtDoNotDisturbChanged.Publish(PayloadDoNotDisturb{DoNotDisturb: enabled})
}

func DoNotDisturb() bool {
	return _service.doNotDisturb
}

func GetAllNotifications() NotificationList {
	return _service.server.notifications
}

func GetNotification(id uint32) (n Notification, ok bool) {
	n, ok = _service.server.notifications[id]
	return
}

func HasNotifications() bool {
	return len(_service.server.notifications) > 0
}

func ClearAllNotifications(notifyEach bool) {
	if notifyEach {
		_service.isClearing = true
	}
	for id := range _service.server.notifications {
		_service.server.closeNotification(id)
	}
	if notifyEach {
		_service.isClearing = false
	}
}

func InvokeAction(id uint32, key string) bool {
	return _service.server.invokeAction(id, key)
}

func CloseNotification(id uint32) {
	_service.server.closeNotification(id)
	_service.evtNotificationClosed.Publish(PayloadNotificationClosed{Id: id, Reason: CloseReasonClosed})
}
