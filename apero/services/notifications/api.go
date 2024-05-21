package notifications

var _service = newService()

func StartService() {
	_service.server.start()
}

// func Register[T Subscriber](handle T, events ...EventType) {
// 	for _, event := range events {
// 		_service.subscribers[event] = append(_service.subscribers[event], handle)
// 	}
// }

// func Unregister(handle Subscriber, events ...EventType) {
// 	for _, event := range events {
// 		handlers := _service.subscribers[event]
// 		for i, h := range handlers {
// 			if h == handle {
// 				_service.subscribers[event] = append(handlers[:i], handlers[i+1:]...)
// 				break
// 			}
// 		}
// 	}
// }

func SetDoNotDisturb(enabled bool) {
	if _service.doNotDisturb == enabled {
		return
	}

	_service.doNotDisturb = enabled

	// for _, handler := range _service.subscribers[EventDoNotDisturbChanged] {
	// 	handler.DoNotDisturbChanged(enabled)
	// }
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
}
