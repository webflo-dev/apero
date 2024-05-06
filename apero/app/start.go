package app

import (
	"webflo-dev/apero/services/hyprland"
	sysStat "webflo-dev/apero/services/system-stats"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var application *gtk.Application

type WindowsLoader func() []*gtk.Window
type UserConfig struct {
	AppId        string
	CssFile      string
	IconFolder   string
	Windows      WindowsLoader
	UseInspector bool
}

func Start(userConfig *UserConfig) int {

	appLogger.Println("load configuration...")
	loadConfig(userConfig)
	appLogger.Println("configuration loaded!")

	appLogger.Println("start IPC...")
	startIPC()
	appLogger.Println("IPC started!")

	appLogger.Println("start application...")

	gtk.Init(nil)

	application, _ := gtk.ApplicationNew(appConfig.AppId, glib.APPLICATION_FLAGS_NONE)
	application.Connect("activate", func() { activate(application, userConfig.UseInspector, userConfig.Windows) })
	appLogger.Println("application started! 🚀")

	return application.Run([]string{})
}

func activate(application *gtk.Application, useInspector bool, windowsLoader func() []*gtk.Window) {
	gtk.SetInteractiveDebugging(useInspector)

	application.Hold()

	if appConfig.IconFolder != "" {
		iconTheme, _ := gtk.IconThemeGetDefault()
		iconTheme.AppendSearchPath(appConfig.IconFolder)
	}

	appLogger.Println("loading CSS from " + appConfig.CssFile)
	ApplyCSS(appConfig.CssFile)
	appLogger.Println("CSS loaded!")

	sysStat.StartService()
	hyprland.StartService()

	if windowsLoader != nil {
		windowsLoader()
	}

}
