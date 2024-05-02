package bar

import (
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

type notificationHandler struct {
	notifications.NotificationsEventHandler
	box  *gtk.Box
	icon *gtk.Image
}

func newNotificationModule() *gtk.Box {
	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 8)
	box.SetName("notifications")

	button, _ := gtk.ButtonNew()

	icon := ui.NewFontSizeImageFromIconName(Icon_NotificationIndicator)
	button.Add(icon)

	box.Add(button)

	handler := &notificationHandler{
		box:  box,
		icon: icon,
	}

	handler.update(notifications.DoNotDisturb(), notifications.HasNotifications())

	notifications.WatchNotifications(handler)

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
			notifications.EmitActionInvoked()
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

func (handler *notificationHandler) update(doNotDisturb bool, HasNotifications bool) {
	ui.ToggleCSSClassFromBool(&handler.box.Widget, "empty", HasNotifications == false)

	if doNotDisturb {
		ui.SetFontSizeImageFromIconName(handler.icon, Icon_NotificationIndicator_DND)
	} else {
		ui.SetFontSizeImageFromIconName(handler.icon, Icon_NotificationIndicator)
	}
}

func (handler *notificationHandler) NewNotification(_ notifications.Notification) {
	handler.update(notifications.DoNotDisturb(), true)
}

func (handler *notificationHandler) DoNotDisturbChanged(enabled bool) {
	handler.update(enabled, notifications.HasNotifications())
}

func (handler *notificationHandler) NotificationsCleared() {
	handler.update(notifications.DoNotDisturb(), notifications.HasNotifications())
}
