package redisClient

import (
	"context"
	"fmt"
	"os"
	"time"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)



func NewClient() (*redis.Client, error) {
	
	password := os.Getenv("REDIS_PASSWORD")
	host := viper.GetString("redis.host")
	port := viper.GetString("redis.port")
	
	addr := fmt.Sprintf("%s:%s", host, port)
	
	opts := &redis.Options{
		Addr:     addr,
		Password: password,
		DB:       viper.GetInt("redis.db"),
	}
	
	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		logrus.Fatalf("Redis ping error: %s", err)
		return nil, err
	}
	
	logrus.Infof("Successfully connected to Redis at %s", addr)
	return client, nil
}