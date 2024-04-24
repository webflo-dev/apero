package bar

import (
	"time"

	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func newDateTime() *gtk.Box {

	dateLabel := gtk.NewLabel(glib.NewDateTimeNowLocal().Format("%A %d %B"))
	timeLabel := gtk.NewLabel(time.Now().Format("15:04"))

	box := gtk.NewBox(gtk.OrientationHorizontal, 16)
	box.SetHExpand(true)
	box.SetHAlign(gtk.AlignCenter)
	box.SetName("date-time")
	box.SetCSSClasses([]string{"date-time"})

	dateContainer := gtk.NewBox(gtk.OrientationHorizontal, 8)
	dateContainer.SetCSSClasses([]string{"date"})
	dateContainer.Append(gtk.NewImageFromIconName("_calendar-day-symbolic"))
	dateContainer.Append(dateLabel)

	timeContainer := gtk.NewBox(gtk.OrientationHorizontal, 8)
	timeContainer.SetCSSClasses([]string{"time"})
	timeContainer.Append(gtk.NewImageFromIconName("_clock-symbolic"))
	timeContainer.Append(timeLabel)

	box.Append(dateContainer)
	box.Append(timeContainer)

	go func() {
		for t := range time.Tick(time.Second) {
			glib.IdleAdd(func() {
				timeLabel.SetText(t.Format("15:04"))
			})
		}
	}()

	return box
}
