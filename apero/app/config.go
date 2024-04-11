package app

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

type appConfiguration struct {
	AppId   string `mapstructure:"app-id"`
	CssFile string `mapstructure:"css-file"`
}

var appConfig appConfiguration

func loadConfig(userConfig *UserConfig) {
	workingDir, err := os.Getwd()
	check(err, "Cannot get working directory")

	appConfig = appConfiguration{
		AppId:   "apero.app",
		CssFile: filepath.Join(workingDir, "apero.css"),
	}

	if userConfig != nil {
		if userConfig.AppId != "" {
			appConfig.AppId = userConfig.AppId
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
				log.Println("user CSS file not found:", cssFile)
				log.Println("default CSS file will be used:", cssFile)
				checkDefaultCssFile()
			}
		} else {
			checkDefaultCssFile()
		}
	}
}

func checkUserCssFile(cssFile string) bool {
	_, err := os.Stat(cssFile)
	return err == nil
}

func checkDefaultCssFile() {
	if _, err := os.Stat(appConfig.CssFile); errors.Is(err, os.ErrNotExist) {
		log.Println("default CSS file not found. CSS won't be applied", appConfig.CssFile)
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
