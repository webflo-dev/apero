package app

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

var cssProvider *gtk.CssProvider

func ApplyCSS(filepath string) {
	if filepath == "" {
		filepath = appConfig.CssFile
	}

	screen, _ := gdk.ScreenGetDefault()

	if cssProvider != nil {
		gtk.RemoveProviderForScreen(screen, cssProvider)
		cssProvider = nil
	}

	cssProvider, _ = gtk.CssProviderNew()
	cssProvider.LoadFromPath(filepath)
	gtk.AddProviderForScreen(screen, cssProvider, gtk.STYLE_PROVIDER_PRIORITY_USER)
}
