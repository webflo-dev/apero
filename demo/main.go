package main

import (
	"demo/modules/bar"
	"os"
	apero "webflo-dev/apero/app"

	"github.com/gotk3/gotk3/gtk"
)

type appEvent struct {
	apero.AppEvent
	windows []*gtk.Window
}

func main() {

	useInspector := false

	if argsWithoutProg := os.Args[1:]; len(argsWithoutProg) != 0 {
		useInspector = argsWithoutProg[0] == "-i" || argsWithoutProg[0] == "--inspector"
	}

	app := apero.NewApp()
	app.SetAppId("demo.app")
	app.UseInspector = useInspector

	handle := &appEvent{}

	app.Start(handle)
}

func (a appEvent) LoadWindows() []*gtk.Window {
	a.windows = []*gtk.Window{
		bar.NewBar(),
	}
	return a.windows
}
