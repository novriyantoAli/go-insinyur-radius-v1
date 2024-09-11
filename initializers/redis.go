package initializers

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var REDIS *redis.Client

func ConnectRedis() {
	host := viper.GetString("redis.host")
	port := viper.GetString("redis.port")

	dsn := fmt.Sprintf("%s:%s", host, port)

	REDIS = redis.NewClient(&redis.Options{
		Addr: dsn,
	})
}
