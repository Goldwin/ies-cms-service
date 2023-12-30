package infra

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Hosts          string        `yaml:"hosts" default:"127.0.0.1:6379"`
	Username       string        `yaml:"username" default:""`
	Password       string        `yaml:"password" default:""`
	MaxRetries     int           `yaml:"maxRetries" default:"3"`
	ReadTimeout    time.Duration `yaml:"readTimeout"`
	WriteTimeout   time.Duration `yaml:"writeTimeout"`
	RouteByLatency bool          `yaml:"routeByLatency"`
}

func NewRedisClient(r *RedisConfig) redis.UniversalClient {
	var redisClient redis.UniversalClient
	var option redis.UniversalOptions
	if r == nil {
		return nil
	}
	addresses := strings.Split(r.Hosts, ",")
	option = redis.UniversalOptions{
		Addrs:          addresses,
		Password:       r.Password,
		MaxRetries:     r.MaxRetries,
		ReadTimeout:    r.ReadTimeout,
		WriteTimeout:   r.WriteTimeout,
		RouteByLatency: r.RouteByLatency,
	}
	if r != nil {
		redisClient = redis.NewUniversalClient(&option)
	} else {
		log.Fatal("Failed to parse redis config")
	}
	str, err := redisClient.Ping(context.Background()).Result()
	if err != nil || str != "PONG" {
		log.Fatal("Failed to connect redis")
	}
	return redisClient
}
