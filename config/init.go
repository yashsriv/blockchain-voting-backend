package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

func init() {
	// read from config file
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("[warn] unable to locate config file")
	}

	// Auto read config values from env
	viper.SetEnvPrefix("voting")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}
