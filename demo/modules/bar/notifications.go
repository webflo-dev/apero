package bar

import (
	"log"
	"webflo-dev/apero/services/notification"
	"webflo-dev/apero/services/notifications"
	"webflo-dev/apero/ui"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

type notificationIndicatorHandler struct {
	notifications.Subscriber
	box  *gtk.Box
	icon *gtk.Image
}

type notificationActionHandler struct {
	id string
}

func (h *notificationActionHandler) doAction1() {
	log.Println("do action1 from", h.id)
}
func (h *notificationActionHandler) doAction2() {
	log.Println("do action2 from", h.id)
}
func doGenericAction() {

}

func newNotificationModule() *gtk.Box {
	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 8)
	box.SetName("notifications")

	button, _ := gtk.ButtonNew()

	icon := ui.NewFontSizeImageFromIconName(Icon_NotificationIndicator)
	button.Add(icon)

	box.Add(button)

	handler := &notificationIndicatorHandler{
		box:  box,
		icon: icon,
	}

	handler.update(notifications.DoNotDisturb(), notifications.HasNotifications())

	notifications.Register(handler, notifications.EventNewNotification, notifications.EventDoNotDisturbChanged, notifications.EventNotificationsCleared)

	button.Connect("button-press-event", func(_ *gtk.Button, ev *gdk.Event) bool {
		btn := gdk.EventButtonNewFromEvent(ev)
		switch btn.Button() {
		case gdk.BUTTON_PRIMARY:
			notifications.SetDoNotDisturb(!notifications.DoNotDisturb())
			return true
		case gdk.BUTTON_SECONDARY:
			notifications.ClearAllNotifications(false)
			return true
		case gdk.BUTTON_MIDDLE:
			notif := notification.NewNotification(&notificationActionHandler{id: "A"}, "Message from Apero (1)", "bla bla bla (1)")
			notif.WithUrgency(notification.UrgencyCritical)
			notif.WithAction("action1", (*notificationActionHandler).doAction1)
			notif.WithAction("action2", (*notificationActionHandler).doAction2)

			id1 := notification.Notify(notif)
			log.Println("Notification ID#1:", id1)

			notifications.InvokeAction(id1, "action2")

			notif2 := notification.NewNotification(&notificationActionHandler{id: "B"}, "Message from Apero (2)", "bla bla bla (2)")
			notif2.WithUrgency(notification.UrgencyCritical)
			notif2.WithAction("action1", (*notificationActionHandler).doAction1)
			notif2.WithAction("action2", (*notificationActionHandler).doAction2)

			id2 := notification.Notify(notif2)
			log.Println("Notification ID#2:", id2)

			notifications.InvokeAction(id2, "action1")

			// notif := &notification.Notification{
			// 	appName: "test",
			// 	actions: notification.Actions{
			// 		"action1": func() {
			// 			fmt.Println("action1")
			// 		},
			// 		"action2": func() {
			// 			fmt.Println("action1")
			// 		},
			// 	},
			// }
			// id := notification.Notify(*notif)

			// log.Println("Notification ID:", id)

			// notifications.InvokeAction(id, "action2")
			return true
		default:
			return false
		}
	})

	// type _ImageData struct {
	// 	/*0*/ Width int
	// 	/*1*/ Height int
	// 	/*2*/ RowStride int
	// 	/*3*/ HasAlpha bool
	// 	/*4*/ BitsPerSample int
	// 	/*5*/ Samples int
	// 	/*6*/ Image []byte
	// }

	// var d = _ImageData{
	// 	r.Max.X, // Width
	// 	r.Max.Y, // Height
	// 	p.Stride,
	// 	true,
	// 	8,
	// 	4,
	// 	p.Pix,
	// }
	// ntf.Hints["image-data"] = dbus.MakeVariant(d)

	return box
}

func (handler *notificationIndicatorHandler) update(doNotDisturb bool, HasNotifications bool) {
	ui.ToggleCSSClassFromBool(&handler.box.Widget, "empty", HasNotifications == false)

	if doNotDisturb {
		ui.SetFontSizeImageFromIconName(handler.icon, Icon_NotificationIndicator_DND)
	} else {
		ui.SetFontSizeImageFromIconName(handler.icon, Icon_NotificationIndicator)
	}
}

func (handler *notificationIndicatorHandler) NewNotification(_ notifications.Notification) {
	handler.update(notifications.DoNotDisturb(), true)
}

func (handler *notificationIndicatorHandler) DoNotDisturbChanged(enabled bool) {
	handler.update(enabled, notifications.HasNotifications())
}

func (handler *notificationIndicatorHandler) NotificationsCleared() {
	handler.update(notifications.DoNotDisturb(), notifications.HasNotifications())
}
