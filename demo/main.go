package main

import (
	"demo/modules/bar"
	"log"
	"os"
	apero "webflo-dev/apero/app"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func main() {

	useInspector := false

	if argsWithoutProg := os.Args[1:]; len(argsWithoutProg) != 0 {
		useInspector = argsWithoutProg[0] == "-i" || argsWithoutProg[0] == "--inspector"
	}

	log.Println("use inspector: ", useInspector)

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
