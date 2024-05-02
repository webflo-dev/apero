package notification

import (
	"fmt"
	"log"

	"webflo-dev/apero/services/notifications"
)

type Urgency uint16

const (
	UrgencyLow      Urgency = 0
	UrgencyCritical Urgency = 2
	UrgencyNormal   Urgency = 1
)

type Actions[T any] map[string]func(T)

type notificationEventHandler[T any] struct {
	notifications.NotificationsEventHandler
	id           uint32
	notification Notification[T]
}

type Notification[T any] struct {
	replaceId    uint32
	appName      string
	icon         string
	summary      string
	body         string
	urgency      Urgency
	category     string
	actions      Actions[T]
	actionHandle T
}

func NewNotification[T any](handle T, summary string, body string) Notification[T] {
	return Notification[T]{
		summary: summary,
		body:    body,
		urgency: UrgencyNormal,
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

func (n *Notification[T]) WithActions(actions Actions[T]) *Notification[T] {
	n.actions = actions
	return n
}

func Notify[T any](notification Notification[T]) uint32 {
	actionKeys := make([]string, 0, len(notification.actions))
	for key := range notification.actions {
		actionKeys = append(actionKeys, key)
	}

	id := notifications.Notify(notification.appName, notification.replaceId, notification.icon, notification.summary, notification.body, actionKeys, nil, 0)

	handle := &notificationEventHandler[T]{
		id:           id,
		notification: notification,
	}
	notifications.WatchNotifications(handle)

	return id
}

func (n notificationEventHandler[T]) ActionInvoked(notificationId uint32, actionKey string) {
	log.Println("notificationEventHandler::ActionInvoked", notificationId, actionKey)
	if n.id != notificationId {
		return
	}

	// n.notification.actions[actionKey](n.notification.actionHandle)

	// if action, ok := n.notification.actions[actionKey]; ok {
	// 	action()
	// }
}

///////////////////////////////////////////////////////////:

// type mystruct struct {
// 	name string
// }

// func (my *mystruct) doA(i int) {
// 	fmt.Printf("[doA]: I'm %s, param is: %d\n", my.name, i)
// }

// func (my *mystruct) doB(i int) {
// 	fmt.Printf("[doB]: I'm %s, param is: %d\n", my.name, i)
// }

// func main() {
// 	my1 := &mystruct{"Bob"}
// 	my2 := &mystruct{"Alice"}
// 	lookupMap := map[string]func(int){
// 			"action1": my1.doA,
// 			"action2": my2.doB,
// 	}

// 	lookupMap["action1"](11)
// 	lookupMap["action2"](22)
// }

// func main() {
// 	lookupMap := map[string]func(*mystruct, int){
// 		"action1": (*mystruct).doA,
// 		"action2": (*mystruct).doB,
// 	}

// 	my1 := &mystruct{"Bob"}
// 	my2 := &mystruct{"Alice"}
// 	lookupMap["action1"](my1, 11)
// 	lookupMap["action2"](my2, 22)
// }

type mystruct struct {
	name string
}

func (my *mystruct) doA() {
	fmt.Printf("[doA]: I'm %s\n", my.name)
}

func (my *mystruct) doB() {
	fmt.Printf("[doB]: I'm %s\n", my.name)
}

// Define a type for the function that takes a receiver of any type T and has no return value
type toto[T any] map[string]func(T)

func sdfdsfdsf() {
	// Use the actionFunc type in the lookupMap
	lookupMap := toto[*mystruct]{
		"action1": (*mystruct).doA,
		"action2": (*mystruct).doB,
	}

	my1 := &mystruct{"Bob"}
	my2 := &mystruct{"Alice"}

	lookupMap["action1"](my1)
	lookupMap["action2"](my2)
}
