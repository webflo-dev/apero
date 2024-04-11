package app

import (
	"log"

	"webflo-dev/apero/ipc"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

var application *gtk.Application
var userConfig *UserConfig

type WindowsLoader func() []*gtk.Window
type UserConfig struct {
	AppId   string
	CssFile string
	Windows WindowsLoader
}

func Start(userConfig *UserConfig) int {

	log.Println("[apero] load configuration...")
	loadConfig(userConfig)
	log.Println("[apero] configuration loaded!")

	log.Println("[apero] loading CSS from " + appConfig.CssFile)
	cssProvider, err := loadCSS(appConfig.CssFile)
	check(err, "Cannot load CSS")
	log.Println("[apero] CSS loaded!")

	log.Println("[apero] start IPC...")
	ipc.StartIPC()
	log.Println("[apero] IPC started!")

	log.Println("[apero] start application...")
	application = gtk.NewApplication(appConfig.AppId, gio.ApplicationFlagsNone)
	application.ConnectActivate(func() { activate(application, cssProvider, userConfig.Windows) })
	log.Println("[apero] application started! 🚀")

	return application.Run([]string{})
}

func activate(application *gtk.Application, cssProvider *gtk.CSSProvider, windowsLoader func() []*gtk.Window) {
	application.Hold()

	if cssProvider != nil {
		applyCSS(cssProvider)
	}

	windowsLoader()

}
