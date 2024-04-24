package app

import (
	"errors"
	"os"
	"path/filepath"

	"webflo-dev/apero/logger"
)

type appConfiguration struct {
	AppId      string
	CssFile    string
	IconFolder string
}

var appConfig appConfiguration

func loadConfig(userConfig *UserConfig) {
	workingDir, err := os.Getwd()
	if err != nil {
		logger.AppLogger.Fatalln("Cannot get working directory", err)
	}

	appConfig = appConfiguration{
		AppId:      "apero.app",
		CssFile:    filepath.Join(workingDir, "apero.css"),
		IconFolder: filepath.Join(workingDir, "icons"),
	}

	if userConfig != nil {
		if userConfig.AppId != "" {
			appConfig.AppId = userConfig.AppId
		}

		if userConfig.IconFolder != "" {
			iconFolder := filepath.Dir(userConfig.IconFolder)
			if iconFolder == "." {
				iconFolder = filepath.Join(workingDir, userConfig.IconFolder)
			} else {
				iconFolder = os.ExpandEnv(userConfig.IconFolder)
			}
			appConfig.IconFolder = iconFolder
		}

		if userConfig.CssFile != "" {
			cssFile := filepath.Dir(userConfig.CssFile)
			if cssFile == "." {
				cssFile = filepath.Join(workingDir, userConfig.CssFile)
			} else {
				cssFile = os.ExpandEnv(userConfig.CssFile)
			}

			if ok := checkUserCssFile(cssFile); ok {
				appConfig.CssFile = cssFile
			} else {
				logger.AppLogger.Println("user CSS file not found:", cssFile)
				logger.AppLogger.Println("default CSS file will be used:", cssFile)
				checkDefaultCssFile()
			}
		} else {
			checkDefaultCssFile()
		}
	}

	logger.AppLogger.Println("AppId:", appConfig.AppId)
	logger.AppLogger.Println("CssFile:", appConfig.CssFile)
	logger.AppLogger.Println("IconFolder:", appConfig.IconFolder)

}

func checkUserCssFile(cssFile string) bool {
	_, err := os.Stat(cssFile)
	return err == nil
}

func checkDefaultCssFile() {
	if _, err := os.Stat(appConfig.CssFile); errors.Is(err, os.ErrNotExist) {
		logger.AppLogger.Println("default CSS file not found. CSS won't be applied", appConfig.CssFile)
		appConfig.CssFile = ""
	}
}

// func _loadConfig() {
// 	viper.SetConfigName("config")
// 	viper.SetConfigType("yaml")

// 	viper.AddConfigPath("$XDG_CONFIG_HOME/apero")
// 	viper.AddConfigPath(".")

// 	err := viper.ReadInConfig()
// 	check(err, "Cannot read configuration file")

// 	err = viper.Unmarshal(&appConfig)
// 	check(err, "Unable to decode configuration file")

// 	appConfig.Dir = filepath.Dir(viper.ConfigFileUsed())
// 	appConfig.CssFile = filepath.Join(appConfig.Dir, appConfig.CssFile)

// 	log.Println(appConfig)
// }
