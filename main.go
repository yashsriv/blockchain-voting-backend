package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"

	"blockchain-voting/http"
)

func main() {
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

	host := viper.GetString("http.host")
	port := viper.GetString("http.port")
	address := fmt.Sprintf("%s:%s", host, port)

	router := http.Router()
	router.Run(address)
}
