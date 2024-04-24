package bar

import (
	"webflo-dev/apero/ui"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
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

	window.SetCSSClasses([]string{"bar"})

	box := gtk.NewCenterBox()
	box.SetStartWidget(newStartBarModule())
	box.SetCenterWidget(newCenterBarModule())
	box.SetEndWidget(newEndBarModule())
	window.SetChild(box)

	window.SetVisible(true)

	return window
}

func newStartBarModule() *gtk.Box {
	container := gtk.NewBox(gtk.OrientationHorizontal, 0)
	container.SetCSSClasses([]string{"start"})

	container.Append(newWorkspaces())

	return container
}

func newCenterBarModule() *gtk.Box {
	container := gtk.NewBox(gtk.OrientationHorizontal, 0)
	container.SetCSSClasses([]string{"center"})
	container.SetHExpand(false)

	// container.Append(newDateTime())

	return container
}

func newEndBarModule() *gtk.Box {
	container := gtk.NewBox(gtk.OrientationHorizontal, 0)
	container.SetCSSClasses([]string{"end"})

	box := gtk.NewBox(gtk.OrientationHorizontal, 8)
	box.SetHExpand(true)
	box.SetHAlign(gtk.AlignEnd)

	// box.Append(newSystemInfo())

	container.Append(box)
	return container
}
