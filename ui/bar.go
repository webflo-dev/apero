package ui

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)


func NewBar() GtkWindow {
	window := NewWindow(WindowProps{ 
		name: "bar",
		layer: LayerTop,
		anchor: PositionTop | PositionLeft | PositionRight,
	})

	window.SetChild(gtk.NewLabel("Hello, World!"))
	window.SetVisible(true)

	return window
}