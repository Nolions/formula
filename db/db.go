package db

import (
	"fmt"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

type ConfigNode struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Config struct {
	Master          ConfigNode    `mapstructure:"master"`
	Slave           ConfigNode    `mapstructure:"slave"`
	Database        string        `mapstructure:"database"`
	Timezone        string        `mapstructure:"timezone"`
	DialTimeout     string        `mapstructure:"dial_timeout"`
	ReadTimeout     string        `mapstructure:"read_timeout"`
	WriteTimeout    string        `mapstructure:"write_timeout"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
}

func (c Config) timezone() (*time.Location, error) {
	return time.LoadLocation(c.Timezone)
}

// New create db instance
func New(conf Config) (*xorm.EngineGroup, error) {
	master, err := newNode(conf.Master, conf)
	if err != nil {
		return nil, err
	}

	slaves := []*xorm.Engine{}
	if conf.Slave.Host != "" {
		slave, err := newNode(conf.Slave, conf)
		if err != nil {
			return nil, err
		}
		slaves = append(slaves, slave)
	}

	return xorm.NewEngineGroup(master, slaves)
}

func newNode(nodeConf ConfigNode, conf Config) (*xorm.Engine, error) {
	timezone, err := conf.timezone()
	if err != nil {
		return nil, err
	}

	db, err := xorm.NewEngine("mysql", buildDSN(nodeConf.Host, nodeConf.Port, nodeConf.Username, nodeConf.Password, conf.Database, conf.DialTimeout, conf.ReadTimeout, conf.WriteTimeout, conf.Timezone))
	if err != nil {
		return nil, err
	}

	db.SetTZDatabase(timezone)
	db.SetTZLocation(timezone)
	db.SetConnMaxLifetime(conf.ConnMaxLifetime)
	db.SetMaxIdleConns(conf.MaxIdleConns)
	db.SetMaxOpenConns(conf.MaxOpenConns)

	return db, nil
}

func buildDSN(host string, port int, username, password, database, dialTimeout, readTimeout, writeTimeout, timezone string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%s&readTimeout=%s&writeTimeout=%s&parseTime=true&loc=%s", username, password, host, port, database, dialTimeout, readTimeout, writeTimeout, url.QueryEscape(timezone))
}
