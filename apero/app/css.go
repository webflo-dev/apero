package app

import (
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

var cssProvider *gtk.CSSProvider

func ApplyCSS(filepath string) {
	if filepath == "" {
		filepath = appConfig.CssFile
	}

	if cssProvider != nil {
		gtk.StyleContextRemoveProviderForDisplay(gdk.DisplayGetDefault(), cssProvider)
		cssProvider = nil
	}

	cssProvider = gtk.NewCSSProvider()
	cssProvider.LoadFromPath(filepath)
	gtk.StyleContextAddProviderForDisplay(gdk.DisplayGetDefault(), cssProvider, gtk.STYLE_PROVIDER_PRIORITY_USER)
}
