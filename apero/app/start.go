package app

import (
	"os"
	"path/filepath"
	"webflo-dev/apero/services/hyprland"
	"webflo-dev/apero/services/notifications"
	sysStat "webflo-dev/apero/services/system-stats"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var _app *AppConfiguration

type AppEvent interface {
	LoadWindows() []*gtk.Window
	ConfigLoaded()
}

type AppConfiguration struct {
	workingDir  string
	cssProvider *gtk.CssProvider
	gtkApp      *gtk.Application

	AppId        string
	CssFile      string
	IconFolder   string
	UseInspector bool
}

func NewApp() *AppConfiguration {
	workingDir, err := os.Getwd()
	if err != nil {
		appLogger.Fatalln("Cannot get working directory", err)
	}

	_app := &AppConfiguration{
		workingDir: workingDir,
		AppId:      "apero.app",
		CssFile:    filepath.Join(workingDir, "apero.css"),
		IconFolder: filepath.Join(workingDir, "icons"),
	}
	return _app
}

func (app *AppConfiguration) Start(handle AppEvent) int {
	appLogger.Println("AppId:", app.AppId)
	appLogger.Println("CssFile:", app.CssFile)
	appLogger.Println("IconFolder:", app.IconFolder)

	appLogger.Println("start IPC...")
	startIPC()
	appLogger.Println("IPC started!")

	appLogger.Println("start application...")

	appLogger.Println("AppId:", app.AppId)
	appLogger.Println("CssFile:", app.CssFile)
	appLogger.Println("IconFolder:", app.IconFolder)

	gtk.Init(nil)

	app.gtkApp, _ = gtk.ApplicationNew(app.AppId, glib.APPLICATION_FLAGS_NONE)
	app.gtkApp.Connect("activate", func() {
		app.startServices()
		app.activate()
		handle.LoadWindows()
	})
	appLogger.Println("application started! 🚀")

	return app.gtkApp.Run([]string{})
}

func (app *AppConfiguration) activate() {
	gtk.SetInteractiveDebugging(app.UseInspector)

	app.gtkApp.Hold()

	if app.IconFolder != "" {
		iconTheme, _ := gtk.IconThemeGetDefault()
		iconTheme.AppendSearchPath(app.IconFolder)
	}

	app.applyCSS()
}

func (app *AppConfiguration) startServices() {
	hyprland.StartService()
	sysStat.StartService()
	notifications.StartService()
}
