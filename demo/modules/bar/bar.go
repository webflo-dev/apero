package bar

import (
	"webflo-dev/apero/ui"

	"github.com/gotk3/gotk3/gtk"
)

func NewBar() *gtk.Window {
	window := ui.NewWindow(ui.WindowProps{
		Name:          "bar",
		Layer:         ui.LayerTop,
		ExclusiveZone: ui.ExclusiveZoneAuto,
		Anchor:        ui.PositionTop | ui.PositionLeft | ui.PositionRight,
	})

	ui.SetMargin(window, ui.PositionTop, 5)
	ui.SetMargin(window, ui.PositionRight, 20)
	ui.SetMargin(window, ui.PositionBottom, 0)
	ui.SetMargin(window, ui.PositionLeft, 20)

	ui.AddCSSClass(&window.Widget, "bar")

	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	box.PackStart(newStartBarModule(), true, true, 0)
	box.Add(newCenterBarModule())
	box.PackEnd(newEndBarModule(), true, true, 0)

	window.Add(box)
	window.SetVisible(true)

	window.ShowAll()
	return window
}

func newStartBarModule() *gtk.Box {
	container, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)

	container.Add(newWorkspaces())

	return container
}

func newCenterBarModule() *gtk.Box {
	container, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	container.SetHExpand(false)

	container.Add(newDateTime())

	return container
}

func newEndBarModule() *gtk.Box {
	container, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)

	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 8)
	box.SetHExpand(true)
	box.SetHAlign(gtk.ALIGN_END)

	box.Add(newSystemInfo())

	container.Add(box)
	return container
}
