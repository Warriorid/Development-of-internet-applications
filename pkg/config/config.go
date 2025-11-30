package config

import (
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceHost string
	ServicePort int
	Redis       RedisConfig
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func NewConfig() (*Config, error) {
	var err error

   configName := "config"
   _ = godotenv.Load()
   if os.Getenv("CONFIG_NAME") != "" {
      configName = os.Getenv("CONFIG_NAME")
   }

   viper.SetConfigName(configName)
   viper.SetConfigType("toml")
   viper.AddConfigPath("./config")
   viper.WatchConfig()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	log.Info("config parsed")

	return cfg, nil
}

func GetRedisHost() string {
	return viper.GetString("redis.host")
}

func GetRedisPort() string {
	return viper.GetString("redis.port")
}

func GetRedisPassword() string {
	return viper.GetString("redis.password")
}

func GetRedisDB() int {
	return viper.GetInt("redis.db")
}