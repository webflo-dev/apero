package app

import (
	"log"
	"os"
	"strings"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func loadCSS(filePath string) (*gtk.CSSProvider, error) {
	if filePath == "" {
		return nil, nil
	}

	cssBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	css := string(cssBytes)

	provider := gtk.NewCSSProvider()
	provider.ConnectParsingError(func(sec *gtk.CSSSection, err error) {
		// Optional line parsing routine.
		loc := sec.StartLocation()
		lines := strings.Split(css, "\n")
		log.Printf("CSS error (%v) at line: %q", err, lines[loc.Lines()])
	})

	provider.LoadFromData(css)

	return provider, nil
}

func applyCSS(cssProvider *gtk.CSSProvider) {
	gtk.StyleContextAddProviderForDisplay(gdk.DisplayGetDefault(), cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
}
