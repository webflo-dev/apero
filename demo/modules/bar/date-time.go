package bar

import (
	"time"
	"webflo-dev/apero/ui"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func newDateTimeModule() *gtk.Box {

	// dateLabel, _ := gtk.LabelNew(glib.NewDateTimeNowLocal().Format("%A %d %B"))
	dateLabel, _ := gtk.LabelNew(time.Now().Format("Monday 02 January"))
	timeLabel, _ := gtk.LabelNew(time.Now().Format("15:04"))

	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 16)
	box.SetHExpand(true)
	box.SetHAlign(gtk.ALIGN_CENTER)
	box.SetName("date-time")
	ui.AddCSSClass(&box.Widget, "date-time")

	dateContainer, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 8)
	ui.AddCSSClass(&dateContainer.Widget, "date")
	dateIcon := ui.NewFontSizeImageFromIconName(Icon_DateTime_Date)
	dateContainer.Add(dateIcon)
	dateContainer.Add(dateLabel)

	timeContainer, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 8)
	ui.AddCSSClass(&timeContainer.Widget, "time")
	timeIcon := ui.NewFontSizeImageFromIconName(Icon_DateTime_Time)
	timeContainer.Add(timeIcon)
	timeContainer.Add(timeLabel)

	box.Add(dateContainer)
	box.Add(timeContainer)

	go func() {
		for t := range time.Tick(time.Second) {
			glib.IdleAdd(func() {
				timeLabel.SetText(t.Format("15:04"))
			})
		}
	}()

	return box
}
