package notifications

import (
	"errors"
)

type NotificationsEventHandler interface {
	NewNotification(notification Notification)
	NotificationRemoved(id uint32)
	NotificationsCleared()
	DoNotDisturbChanged(enabled bool)
}

func GetServerCapabilities() []string {
	return server.GetCapabilities()
}
func GetServerInformation() (string, string, string, string) {
	return server.GetServerInformation()
}

func WatchNotifications(handler NotificationsEventHandler) error {
	if server.started == false {
		Logger.Println("notifications server is not started")
		return errors.New("notifications server is not started")
	}

	server.subscribers = append(server.subscribers, handler)

	return nil
}

func SetDoNotDisturb(enabled bool) {
	if server.doNotDisturb == enabled {
		return
	}

	server.doNotDisturb = enabled

	for _, handler := range server.subscribers {
		handler.DoNotDisturbChanged(enabled)
	}
}

func ClearAllNotifications(notifyEach bool) {
	if notifyEach == false {
		clear(server.notifications)
		for _, handler := range server.subscribers {
			handler.NotificationsCleared()
		}
	} else {
		ids := make([]uint32, 0, len(server.notifications))
		for id := range server.notifications {
			ids = append(ids, id)
		}
		for id := range ids {
			server.CloseNotification(uint32(id))
		}
	}
}

func GetNotifications() []Notification {
	notifications := make([]Notification, 0, len(server.notifications))
	for _, n := range server.notifications {
		notifications = append(notifications, n)
	}
	return notifications
}

func GetNotification(id uint32) (Notification, bool) {
	n, ok := server.notifications[id]
	return n, ok
}

func HasNotifications() bool {
	return len(server.notifications) > 0
}

func DoNotDisturb() bool {
	return server.doNotDisturb
}

func InvokeAction(notificationId uint32, actionKey string) {

}
