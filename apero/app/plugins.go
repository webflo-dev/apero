package app

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"plugin"
	"webflo-dev/apero/logger"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type AppPlugin interface {
	Windows() []*gtk.Window
}

func loadPlugins(configDir string, plugins []string) {
	log.Println(plugins)

	for _, item := range plugins {

		plugin, err := plugin.Open(filepath.Join(configDir, item))
		if err != nil {
			logger.AppLogger.Fatalf("Cannot load plugins. %w", err)
		}

		appPlugins, err := lookUpSymbol[AppPlugin](plugin, "Windows")
		if err != nil {
			logger.AppLogger.Fatalf("Cannot find Windows symbol. %w", err)
		}

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
