package db_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/redpkg/formula/v2/db"
	"github.com/stretchr/testify/assert"
)

func TestNewOnlyMaster(t *testing.T) {
	assert := assert.New(t)

	conf := newConfig()
	conf.Master = newConfigNode("master")

	db, err := db.New(conf)
	if !assert.NoError(err) {
		return
	}

	assert.Len(db.Slaves(), 0)
	assert.Equal(db.Master(), db.Slave())

	fmt.Printf("%+v\n", db.Master())
}

func TestNewMasterSlave(t *testing.T) {
	assert := assert.New(t)

	conf := newConfig()
	conf.Master = newConfigNode("master")
	conf.Slave = newConfigNode("slave")

	db, err := db.New(conf)
	if !assert.NoError(err) {
		return
	}

	assert.Len(db.Slaves(), 1)
	assert.NotEqual(db.Master(), db.Slave())

	fmt.Printf("%+v\n", db.Master())
	fmt.Printf("%+v\n", db.Slave())
}

func newConfig() db.Config {
	return db.Config{
		Database:        "default",
		Timezone:        "UTC",
		DialTimeout:     "10s",
		ReadTimeout:     "30s",
		WriteTimeout:    "60s",
		ConnMaxLifetime: time.Hour,
		MaxIdleConns:    5,
		MaxOpenConns:    10,
	}
}

func newConfigNode(role string) db.ConfigNode {
	return db.ConfigNode{
		Host:     "localhost",
		Port:     3306,
		Username: role,
		Password: role,
	}
}
