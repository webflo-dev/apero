package ui

import "github.com/gotk3/gotk3/gtk"

func AddCSSClass(widget *gtk.Widget, className string) {
	ctx, _ := widget.GetStyleContext()
	ctx.AddClass(className)
}

func RemoveCSSClass(widget *gtk.Widget, className string) {
	ctx, _ := widget.GetStyleContext()
	ctx.RemoveClass(className)
}

func ToggleCSSClass(widget *gtk.Widget, className string, shouldAdd bool) {
	ctx, _ := widget.GetStyleContext()
	if shouldAdd {
		ctx.AddClass(className)
	} else {
		ctx.RemoveClass(className)
	}
}

func ToggleCSSClassFromFunc(widget *gtk.Widget, className string, shouldAdd func() bool) {
	ctx, _ := widget.GetStyleContext()
	if shouldAdd() {
		ctx.AddClass(className)
	} else {
		ctx.RemoveClass(className)
	}
}
