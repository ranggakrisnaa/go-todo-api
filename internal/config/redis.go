package config

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type RedisConfig struct {
	Host        string
	Port        string
	MaxActive   int
	MaxIdle     int
	IdleTimeout int
}

func NewRedisPool(cfg *RedisConfig) *redis.Pool {
	return &redis.Pool{
		MaxActive:   cfg.MaxActive,
		MaxIdle:     cfg.MaxIdle,
		Wait:        true,
		IdleTimeout: time.Duration(cfg.IdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", cfg.Host+":"+cfg.Port)
		},
	}
}
