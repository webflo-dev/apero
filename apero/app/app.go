package app

import (
	"os"
	"path/filepath"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func (app *AppConfiguration) SetAppId(appId string) {
	if appId == "" {
		panic("AppId cannot be empty")
	}
	app.AppId = appId
}

func (app *AppConfiguration) SetCssFile(cssFile string) {
	if cssFile == "" {
		panic("css file path cannot be empty")
	}

	path := filepath.Dir(cssFile)
	if path == "." {
		path = filepath.Join(app.workingDir, cssFile)
	} else {
		path = os.ExpandEnv(cssFile)
	}

	fi, err := os.Stat(path)

	if err == nil && !fi.IsDir() {
		app.CssFile = path
		appLogger.Println("user CSS file found:", cssFile)
		return
	}

	if fi.IsDir() {
		appLogger.Println("user CSS file is not a file:", cssFile)
	} else {
		appLogger.Println("user CSS file not found:", cssFile)
	}

	appLogger.Println("default CSS file will be used:", app.CssFile)
	if _, err := os.Stat(app.CssFile); err != nil {
		appLogger.Println("default CSS not found. CSS won't be applied:", app.CssFile)
		app.CssFile = ""
	}
}

func (app *AppConfiguration) SetIconFolder(folderPath string) {
	if folderPath == "" {
		panic("icon folder path cannot be empty")
	}

	path := filepath.Dir(folderPath)
	if path == "." {
		path = filepath.Join(app.workingDir, folderPath)
	} else {
		path = os.ExpandEnv(folderPath)
	}

	fi, err := os.Stat(folderPath)

	if err == nil && fi.IsDir() {
		app.IconFolder = folderPath
		appLogger.Println("user icon folder found:", folderPath)
		return
	}

	if !fi.IsDir() {
		appLogger.Println("user icon folder is not a directory:", folderPath)
	} else {
		appLogger.Println("user icon folder not found:", folderPath)
	}

	appLogger.Println("default icon folder will be used:", app.IconFolder)
	if _, err := os.Stat(app.IconFolder); err != nil {
		appLogger.Println("default icon folder not found. icons won't be applied:", app.IconFolder)
		app.IconFolder = ""
	}
}

func (app *AppConfiguration) SetInspector(use bool) {
	app.UseInspector = use
}

func (app *AppConfiguration) applyCSS() {
	if app.CssFile == "" {
		return
	}

	screen, _ := gdk.ScreenGetDefault()

	if app.cssProvider != nil {
		gtk.RemoveProviderForScreen(screen, app.cssProvider)
		app.cssProvider = nil
	}

	app.cssProvider, _ = gtk.CssProviderNew()
	app.cssProvider.LoadFromPath(app.CssFile)
	gtk.AddProviderForScreen(screen, app.cssProvider, gtk.STYLE_PROVIDER_PRIORITY_USER)
}
