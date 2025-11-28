package config

import (
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewRedisClient(viper *viper.Viper, log *logrus.Logger) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr: viper.GetString("database.redis.hots") + ":" + viper.GetString("database.redis.port"),
		DB:   viper.GetInt("redis.db"),
	})

	return redisClient
}
