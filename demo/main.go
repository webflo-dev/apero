package main

import (
	"demo/modules/bar"
	"os"
	apero "webflo-dev/apero/app"

	"github.com/gotk3/gotk3/gtk"
)

func main() {

	useInspector := false

	if argsWithoutProg := os.Args[1:]; len(argsWithoutProg) != 0 {
		useInspector = argsWithoutProg[0] == "-i" || argsWithoutProg[0] == "--inspector"
	}

	apero.Start(&apero.UserConfig{
		Windows:      GetWindows,
		UseInspector: useInspector,
	})
}

func GetWindows() []*gtk.Window {
	return []*gtk.Window{
		bar.NewBar(),
	}
}
