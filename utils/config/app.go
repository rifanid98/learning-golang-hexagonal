package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	//driver
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	// set setting based on env
	LoadEnvVars()
	MongoConnect()
	//OpenDbPool()
}

func LoadEnvVars() {
	// Bind OS environment variable
	viper.SetEnvPrefix("learning_golang_hexagonal")

	err := viper.BindEnv("env")
	if err != nil {
		panic(fmt.Errorf("Fatal error setting file: %s", err))
	}

	err = viper.BindEnv("app_path")
	if err != nil {
		panic(fmt.Errorf("Fatal error setting file: %s", err))
	}

	cwd, _ := os.Getwd()
	dirString := strings.Split(cwd, "learning-golang-hexagonal")
	dir := strings.Join([]string{dirString[0], "learning-golang-hexagonal"}, "")
	AppPath := dir

	cfg := "config"
	var env string
	if viper.Get("env") != nil {
		env = viper.Get("env").(string)
	}
	if strings.HasPrefix(env, "dev") {
		cfg += "_dev"
	} else if strings.HasPrefix(env, "test") {
		cfg += "_test"
		if viper.Get("app_path") != nil {
			AppPath = viper.Get("app_path").(string)
		}
	}

	viper.SetConfigName(cfg)
	viper.SetConfigType("json")
	viper.AddConfigPath(AppPath)

	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error setting file: %s", err))
	}
}
