package main

import (
	"fmt"

	"github.com/spf13/viper"

	_ "blockchain-voting/config"
	"blockchain-voting/http"
)

func main() {

	host := viper.GetString("http.host")
	port := viper.GetString("http.port")
	address := fmt.Sprintf("%s:%s", host, port)

	router := http.Router()
	router.Run(address)
}
