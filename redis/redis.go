package redis

import (
	goredis "github.com/go-redis/redis/v8"
)

// Config redis config
type Config struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// New create redis instance
func New(conf Config) (*goredis.Client, error) {
	return goredis.NewClient(&goredis.Options{
		Addr:     conf.Address,
		Password: conf.Password,
		DB:       conf.DB,
	}), nil
}
