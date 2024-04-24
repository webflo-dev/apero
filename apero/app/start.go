package app

import (
	"webflo-dev/apero/logger"

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

	logger.AppLogger.Println("load configuration...")
	loadConfig(userConfig)
	logger.AppLogger.Println("configuration loaded!")

	logger.AppLogger.Println("start IPC...")
	startIPC()
	logger.AppLogger.Println("IPC started!")

	logger.AppLogger.Println("start application...")

	gtk.Init(nil)

	application, _ := gtk.ApplicationNew(appConfig.AppId, glib.APPLICATION_FLAGS_NONE)
	application.Connect("activate", func() { activate(application, userConfig.UseInspector, userConfig.Windows) })
	logger.AppLogger.Println("application started! 🚀")

	return application.Run([]string{})
}

func activate(application *gtk.Application, useInspector bool, windowsLoader func() []*gtk.Window) {
	gtk.SetInteractiveDebugging(useInspector)

	application.Hold()

	if appConfig.IconFolder != "" {
		iconTheme, _ := gtk.IconThemeGetDefault()
		iconTheme.AppendSearchPath(appConfig.IconFolder)
	}

	logger.AppLogger.Println("loading CSS from " + appConfig.CssFile)
	ApplyCSS(appConfig.CssFile)
	logger.AppLogger.Println("CSS loaded!")

	if windowsLoader != nil {
		windowsLoader()
	}

}
