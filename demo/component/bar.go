package component

import (
	"webflo-dev/apero/ui"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func NewBar() *gtk.Window {
	window := ui.NewWindow(ui.WindowProps{
		Name:   "bar",
		Layer:  ui.LayerTop,
		Anchor: ui.PositionTop | ui.PositionLeft | ui.PositionRight,
	})

	window.SetChild(gtk.NewLabel("Hello, World!"))
	window.SetVisible(true)

	return window
}
