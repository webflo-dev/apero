package notifications

// type EventType string

// const (
// 	EventNewNotification      EventType = "NewNotification"
// 	EventNotificationClosed   EventType = "NotificationClosed"
// 	EventDoNotDisturbChanged  EventType = "DoNotDisturbchanged"
// 	EventNotificationsChanged EventType = "NotificationsChanged"
// )

// type Subscriber interface {
// 	NewNotification(notification Notification)
// 	NotificationClosed(id uint32, reason uint32)
// 	DoNotDisturbChanged(enabled bool)
// 	NotificationsChanged(isEmpty bool)
// }

type notificationService struct {
	doNotDisturb bool
	// subscribers  map[EventType][]Subscriber
	isClearing bool

	server *server
}

func newService() *notificationService {
	return &notificationService{
		doNotDisturb: false,
		// subscribers:  make(map[EventType][]Subscriber),
		isClearing: false,
		server:     newServer(),
	}
}

func (s *notificationService) newNotification(notification Notification) {
	// if s.doNotDisturb == false {
	// 	for _, handler := range s.subscribers[EventNewNotification] {
	// 		handler.NewNotification(notification)
	// 	}
	// }
}

func (s *notificationService) notificationClosed(id uint32, reason uint32) {
	// if s.doNotDisturb == false && s.isClearing == false {
	// 	for _, subscriber := range s.subscribers[EventNotificationClosed] {
	// 		subscriber.NotificationClosed(id, reason)
	// 	}
	// }
}

func (s *notificationService) notificationsChanged(isEmpty bool) {
	// if s.isClearing == false {
	// 	for _, subscriber := range s.subscribers[EventNotificationsChanged] {
	// 		subscriber.NotificationsChanged(isEmpty)
	// 	}
	// }
}
