package app

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"plugin"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type AppPlugin interface {
	Windows() []*gtk.Window
}

func loadPlugins(configDir string, plugins []string) {
	log.Println(plugins)

	for _, item := range plugins {

		plugin, err := plugin.Open(filepath.Join(configDir, item))
		check(err, "Cannot load plugins")

		appPlugins, err := lookUpSymbol[AppPlugin](plugin, "Windows")
		check(err, "Cannot find Windows symbol")

		(*appPlugins).Windows()
	}
}

func lookUpSymbol[M any](plugin *plugin.Plugin, symbolName string) (*M, error) {
	symbol, err := plugin.Lookup(symbolName)
	if err != nil {
		return nil, err
	}
	switch symbol.(type) {
	case *M:
		return symbol.(*M), nil
	case M:
		result := symbol.(M)
		return &result, nil
	default:
		return nil, errors.New(fmt.Sprintf("unexpected type from module symbol: %T", symbol))
	}
}
