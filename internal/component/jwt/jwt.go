package jwt

import (
	"github.com/dobyte/due-chat/internal/component/redis"
	"github.com/dobyte/due/v2/etc"
	"github.com/dobyte/due/v2/log"
	"github.com/dobyte/jwt"
	"sync"
)

var (
	once     sync.Once
	instance *jwt.JWT
)

type (
	JWT     = jwt.JWT
	Token   = jwt.Token
	Payload = jwt.Payload
)

type Config struct {
	Issuer        string        `json:"issuer"`
	ValidDuration int           `json:"validDuration"`
	SecretKey     string        `json:"secretKey"`
	IdentityKey   string        `json:"identityKey"`
	Locations     string        `json:"locations"`
	Store         *redis.Config `json:"store"`
}

// Instance JWT实例
func Instance() *JWT {
	once.Do(func() {
		instance = NewInstance("etc.jwt.default")
	})

	return instance
}

// NewInstance 使用配置创建
func NewInstance[T string | Config | *Config](config T) *JWT {
	var (
		conf *Config
		v    any = config
	)

	switch c := v.(type) {
	case string:
		conf = &Config{}
		if err := etc.Get(c).Scan(conf); err != nil {
			log.Fatalf("load jwt config failed: %v", err)
		}
	case Config:
		conf = &c
	case *Config:
		conf = c
	}

	opts := make([]jwt.Option, 0, 6)
	opts = append(opts, jwt.WithIssuer(conf.Issuer))
	opts = append(opts, jwt.WithIdentityKey(conf.IdentityKey))
	opts = append(opts, jwt.WithSecretKey(conf.SecretKey))
	opts = append(opts, jwt.WithValidDuration(conf.ValidDuration))
	opts = append(opts, jwt.WithLookupLocations(conf.Locations))

	if conf.Store != nil {
		opts = append(opts, jwt.WithStore(&store{redis: redis.NewInstance(conf.Store)}))
	}

	jt, err := jwt.NewJWT(opts...)
	if err != nil {
		log.Fatalf("new a jwt instance failed: %v", err)
	}

	return jt
}
