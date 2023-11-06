package config

import (
	"fmt"
	"os"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

func ConnectRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%v:%v", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		DB:   1,
	})

	_, err := client.Ping().Result()
	if err != nil {
		logrus.Error(err.Error())
		return nil
	}

	logrus.Info("Connection established")

	return client
}
