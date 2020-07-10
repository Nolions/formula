package db

import (
	"errors"
	"net/url"
	"strings"
	"time"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

// ConfigNode database node config
type ConfigNode struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Address  string `mapstructure:"address"`
}

// Config database config
type Config struct {
	Driver          string     `mapstructure:"driver"`
	Database        string     `mapstructure:"database"`
	Master          ConfigNode `mapstructure:"master"`
	Slave           ConfigNode `mapstructure:"slave"`
	DialTimeout     string     `mapstructure:"dial_timeout"`
	ReadTimeout     string     `mapstructure:"read_timeout"`
	WriteTimeout    string     `mapstructure:"write_timeout"`
	DBTimezone      string     `mapstructure:"db_timezone"`
	AppTimezone     string     `mapstructure:"app_timezone"`
	ConnMaxLifeTime string     `mapstructure:"conn_max_life_time"`
	MaxIdleConns    int        `mapstructure:"max_idle_conns"`
	MaxOpenConns    int        `mapstructure:"max_open_conns"`
}

func (c Config) dbTimezone() (*time.Location, error) {
	return time.LoadLocation(c.DBTimezone)
}

func (c Config) appTimezone() (*time.Location, error) {
	return time.LoadLocation(c.AppTimezone)
}

func (c Config) connMaxLifeTime() (time.Duration, error) {
	return time.ParseDuration(c.ConnMaxLifeTime)
}

// New create db instance
func New(conf Config) (*xorm.EngineGroup, error) {
	if conf.Master.Address == "" {
		return nil, errors.New("master cannot be nil, master db is required")
	}

	master, err := newNode(conf, conf.Master)
	if err != nil {
		return nil, err
	}

	slaves := []*xorm.Engine{}
	if conf.Slave.Address != "" {
		slave, err := newNode(conf, conf.Slave)
		if err != nil {
			return nil, err
		}
		slaves = append(slaves, slave)
	}

	return xorm.NewEngineGroup(master, slaves)
}

func newNode(conf Config, nodeConf ConfigNode) (*xorm.Engine, error) {
	dbTimezone, err := conf.dbTimezone()
	if err != nil {
		return nil, err
	}
	appTimezone, err := conf.appTimezone()
	if err != nil {
		return nil, err
	}
	connMaxLifeTime, err := conf.connMaxLifeTime()
	if err != nil {
		return nil, err
	}

	db, err := xorm.NewEngine(conf.Driver, buildDSN(nodeConf.Username, nodeConf.Password, nodeConf.Address, conf.Database, conf.DialTimeout, conf.ReadTimeout, conf.WriteTimeout, conf.DBTimezone))
	if err != nil {
		return nil, err
	}

	db.SetTZDatabase(dbTimezone)
	db.SetTZLocation(appTimezone)
	db.SetConnMaxLifetime(connMaxLifeTime)
	db.SetMaxIdleConns(conf.MaxIdleConns)
	db.SetMaxOpenConns(conf.MaxOpenConns)

	return db, nil
}

func buildDSN(username, password, address, database, dialTimeout, readTimeout, writeTimeout, dbTimezone string) string {
	var s strings.Builder

	s.WriteString(username)
	s.WriteString(":")
	s.WriteString(password)
	s.WriteString("@tcp(")
	s.WriteString(address)
	s.WriteString(")/")
	s.WriteString(database)
	s.WriteString("?timeout=")
	s.WriteString(dialTimeout)
	s.WriteString("&readTimeout=")
	s.WriteString(readTimeout)
	s.WriteString("&writeTimeout=")
	s.WriteString(writeTimeout)
	s.WriteString("&parseTime=true&loc=")
	s.WriteString(url.QueryEscape(dbTimezone))

	return s.String()
}
