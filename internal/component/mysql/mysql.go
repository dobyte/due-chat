package mysql

import (
	"strings"
	"sync"
	"time"

	_ "github.com/dobyte/due-chat/internal/component/mysql/serializer"
	"github.com/dobyte/due/v2/etc"
	"github.com/dobyte/due/v2/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	loggers "gorm.io/gorm/logger"
)

var (
	once     sync.Once
	instance *gorm.DB
)

type Config struct {
	DSN             string `json:"dsn"`
	LogLevel        string `json:"logLevel"`
	SlowThreshold   int    `json:"slowThreshold"`
	MaxIdleConns    int    `json:"maxIdleConns"`
	MaxOpenConns    int    `json:"maxOpenConns"`
	ConnMaxLifetime int    `json:"connMaxLifetime"`
}

func Instance() *gorm.DB {
	once.Do(func() {
		instance = NewInstance("etc.mysql.default")
	})

	return instance
}

// NewInstance 新建实例
func NewInstance[T string | Config | *Config](config T) *gorm.DB {
	var (
		conf *Config
		v    any = config
	)

	switch c := v.(type) {
	case string:
		conf = &Config{}
		if err := etc.Get(c).Scan(conf); err != nil {
			log.Fatalf("load mysql config failed: %v", err)
		}
	case Config:
		conf = &c
	case *Config:
		conf = c
	}

	var logLevel loggers.LogLevel
	switch strings.ToLower(conf.LogLevel) {
	case "silent":
		logLevel = loggers.Silent
	case "error":
		logLevel = loggers.Error
	case "warn":
		logLevel = loggers.Warn
	case "info":
		logLevel = loggers.Info
	default:
		logLevel = loggers.Warn
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: conf.DSN,
	}), &gorm.Config{Logger: &logger{
		logLevel:                  logLevel,
		ignoreRecordNotFoundError: true,
		slowThreshold:             time.Duration(conf.SlowThreshold) * time.Millisecond,
	}})
	if err != nil {
		log.Fatalf("establish mysql connection failed: %v", err)
	}

	database, err := db.DB()
	if err != nil {
		log.Fatalf("establish mysql connection failed: %v", err)
	}

	if conf.MaxIdleConns > 0 {
		database.SetMaxIdleConns(conf.MaxIdleConns)
	}

	if conf.MaxOpenConns > 0 {
		database.SetMaxOpenConns(conf.MaxOpenConns)
	}

	if conf.ConnMaxLifetime > 0 {
		database.SetConnMaxLifetime(time.Duration(conf.ConnMaxLifetime) * time.Second)
	}

	return db
}
