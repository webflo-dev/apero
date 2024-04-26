package notifications

import "time"

type notificationServer struct {
	counter       int
	notifications map[int]notification
}

func newNotificationServer() *notificationServer {
	return &notificationServer{
		counter:       0,
		notifications: make(map[int]notification),
	}
}

func (server notificationServer) GetCapabilities() []string {
	return []string{
		"action-icons",
		"actions",
		"body",
		"body-hyperlinks",
		// "body-images",
		"body-markup",
		// "icon-multi",
		"icon-static",
		"persistence",
		// "sound",
	}
}

func (server notificationServer) GetServerInformation() (string, string, string, string) {
	return "apero", "webflo-dev", "0.1", "1.2"
}

func (server notificationServer) Notify(appName string, replacesId int, appIcon string, summary string, body string, actions []string, hints hints, expireTimeout int) int {
	// Logger.Println("Notify!!", appName, replacesId, appIcon, summary, body, actions, hints, expireTimeout)

	id := replacesId
	if id == 0 {
		server.counter++
		id = server.counter
	}

	n := notification{
		id:        id,
		appName:   appName,
		appIcon:   appIcon,
		summary:   summary,
		body:      body,
		actions:   actions,
		hints:     hints,
		timestamp: time.Now().Unix(),
	}

	server.notifications[n.id] = n

	Logger.Printf("Notification > %+v\n", n)

	return n.id
}

func (server notificationServer) CloseNotification(id int) {
	delete(server.notifications, id)
}
