package main

import (
	"demo/component"
	apero "webflo-dev/apero/app"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func main() {
	apero.Start(&apero.UserConfig{
		Windows: GetWindows,
		// CssFile: "demo.css",
	})
}

func GetWindows() []*gtk.Window {
	return []*gtk.Window{
		component.NewBar(),
	}
}
