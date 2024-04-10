package app

import (
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/spf13/viper"
)

var application *gtk.Application

func Start() int {
	loadConfig()

	application = gtk.NewApplication(viper.GetString("app-id"), gio.ApplicationFlagsNone)
	application.ConnectActivate(func() { activate(application) })

	return application.Run([]string{})
}


func activate(application *gtk.Application) {
	application.Hold()
	
	loadCSS()
}


