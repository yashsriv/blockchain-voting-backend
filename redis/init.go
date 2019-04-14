package redis

import (
	"fmt"

	radix "github.com/mediocregopher/radix/v3"
	"github.com/spf13/viper"
)

var Client radix.Client

func init() {
	network := viper.GetString("redis.network")
	host := viper.GetString("redis.host")
	port := viper.GetString("redis.port")
	address := fmt.Sprintf("%s:%s", host, port)
	poolSize := viper.GetInt("redis.pool")

	pool, err := radix.NewPool(network, address, poolSize)
	if err != nil {
		// handle error
		panic(err)
	}

	Client = pool
}
