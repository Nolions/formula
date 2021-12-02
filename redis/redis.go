package redis

import (
	"fmt"

	goredis "github.com/go-redis/redis/v8"
)

type Config struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// New create redis instance
func New(conf Config) (*goredis.Client, error) {
	return goredis.NewClient(&goredis.Options{
		Addr:     buildAddress(conf.Host, conf.Port),
		Password: conf.Password,
		DB:       conf.DB,
	}), nil
}

func buildAddress(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}
