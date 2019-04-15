package main

import (
	"fmt"

	"github.com/spf13/viper"

	_ "blockchain-voting/config"
	"blockchain-voting/http"
	"ethlib"
)

func main() {

	host := viper.GetString("http.host")
	port := viper.GetString("http.port")
	address := fmt.Sprintf("%s:%s", host, port)
	var err error
	http.VC, err = ethlib.NewVotingContractWrapper()
	if err != nil {
		fmt.Printf("Unrecoverable error occured %v", err)
		return
	}
	router := http.Router()
	router.Run(address)
}
