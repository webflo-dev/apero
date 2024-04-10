package app

import (
	"log"
	"os"
	"strings"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/spf13/viper"
)



func check(err error) {
	if err != nil {
		log.Fatalf("fatal error: %w", err)
	}
}

func loadConfig(){
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$XDG_CONFIG_HOME/apero")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	check(err)
}

func loadCSS() {
	cssBytes, err := os.ReadFile(viper.GetString("css-file"))
	if err != nil {
		log.Println("Failed to read CSS file:", err)
		return
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

	gtk.StyleContextAddProviderForDisplay(
		gdk.DisplayGetDefault(), provider,
		gtk.STYLE_PROVIDER_PRIORITY_APPLICATION,
	)
}

