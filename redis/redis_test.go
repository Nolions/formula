package redis_test

import (
	"testing"

	"github.com/redpkg/formula/redis"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)

	_, err := redis.New(newConfig())

	assert.NoError(err)
}

func newConfig() redis.Config {
	return redis.Config{
		Address:  "localhost:6379",
		Password: "",
		DB:       0,
	}
}
