package app

import (
	"webflo-dev/apero/logger"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
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

	logger.AppLogger.Println("AppId:", appConfig.AppId)
	logger.AppLogger.Println("CssFile:", appConfig.CssFile)
	logger.AppLogger.Println("IconFolder:", appConfig.IconFolder)

	logger.AppLogger.Println("start IPC...")
	startIPC()
	logger.AppLogger.Println("IPC started!")

	logger.AppLogger.Println("start application...")
	application = gtk.NewApplication(appConfig.AppId, gio.ApplicationFlagsNone)
	application.ConnectActivate(func() { activate(application, userConfig.UseInspector, userConfig.Windows) })
	logger.AppLogger.Println("application started! 🚀")

	return application.Run([]string{})
}

func activate(application *gtk.Application, useInspector bool, windowsLoader func() []*gtk.Window) {
	gtk.WindowSetInteractiveDebugging(useInspector)

	application.Hold()

	iconTheme := gtk.IconThemeGetForDisplay(gdk.DisplayGetDefault())
	if appConfig.IconFolder != "" {
		iconTheme.AddSearchPath(appConfig.IconFolder)
	}

	logger.AppLogger.Println("loading CSS from " + appConfig.CssFile)
	ApplyCSS(appConfig.CssFile)
	logger.AppLogger.Println("CSS loaded!")

	windowsLoader()

}
