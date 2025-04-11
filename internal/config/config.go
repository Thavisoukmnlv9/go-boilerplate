package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBConn     string `mapstructure:"DB_CONN"`
	ServerPort string `mapstructure:"SERVER_PORT"`
	JWTSecret  string `mapstructure:"JWT_SECRET"`
	RedisAddr  string `mapstructure:"REDIS_ADDR"`
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		panic("Failed to read .env file: " + err.Error())
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		panic("Failed to unmarshal config: " + err.Error())
	}
	return &cfg
}
