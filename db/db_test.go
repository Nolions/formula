package db_test

import (
	"testing"
	"time"

	"github.com/redpkg/formula/db"
	"github.com/stretchr/testify/assert"
)

func TestNewNoMaster(t *testing.T) {
	assert := assert.New(t)

	_, err := db.New(newConfig())

	assert.EqualError(err, "master cannot be nil, master db is required")
}

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
}

func newConfig() db.Config {
	return db.Config{
		Driver:          "mysql",
		Database:        "test",
		DialTimeout:     "10s",
		ReadTimeout:     "30s",
		WriteTimeout:    "60s",
		DBTimezone:      "UTC",
		AppTimezone:     "Asia/Taipei",
		ConnMaxLifeTime: time.Second * 0,
		MaxIdleConns:    2,
		MaxOpenConns:    0,
	}
}

func newConfigNode(role string) db.ConfigNode {
	return db.ConfigNode{
		Username: role,
		Password: role,
		Address:  "localhost:3306",
	}
}
