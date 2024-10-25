package redis

import (
	"github.com/dobyte/due/v2/etc"
	"github.com/dobyte/due/v2/log"
	"github.com/go-redis/redis/v8"
	"sync"
)

var (
	once     sync.Once
	instance Redis
)

type (
	Redis  = redis.UniversalClient
	Script = redis.Script
)

type Config struct {
	Addrs      []string `json:"addrs"`
	DB         int      `json:"db"`
	Username   string   `json:"username"`
	Password   string   `json:"password"`
	MaxRetries int      `json:"maxRetries"`
}

// Instance 获取单例
func Instance() Redis {
	once.Do(func() {
		instance = NewInstance("etc.redis.default")
	})

	return instance
}

// NewInstance 新建实例
func NewInstance[T string | Config | *Config](config T) Redis {
	var (
		conf *Config
		v    any = config
	)

	switch c := v.(type) {
	case string:
		conf = &Config{}
		if err := etc.Get(c).Scan(conf); err != nil {
			log.Fatalf("load redis config failed: %v", err)
		}
	case Config:
		conf = &c
	case *Config:
		conf = c
	}

	return redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:      conf.Addrs,
		DB:         conf.DB,
		Username:   conf.Username,
		Password:   conf.Password,
		MaxRetries: conf.MaxRetries,
	})
}
