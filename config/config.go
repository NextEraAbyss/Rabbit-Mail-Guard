package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	RabbitMQ struct {
		Host     string
		Port     string
		Username string
		Password string
		Queue    string
		Exchange string
	}
	Redis struct {
		Host     string
		Port     string
		Password string
		DB       int
	}
	SMTP struct {
		Host     string
		Port     int
		Username string
		Password string
	}
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{}
	
	// RabbitMQ 配置
	config.RabbitMQ.Host = viper.GetString("RABBITMQ_HOST")
	config.RabbitMQ.Port = viper.GetString("RABBITMQ_PORT")
	config.RabbitMQ.Username = viper.GetString("RABBITMQ_USERNAME")
	config.RabbitMQ.Password = viper.GetString("RABBITMQ_PASSWORD")
	config.RabbitMQ.Queue = viper.GetString("RABBITMQ_QUEUE")
	config.RabbitMQ.Exchange = viper.GetString("RABBITMQ_EXCHANGE")

	// Redis 配置
	config.Redis.Host = viper.GetString("REDIS_HOST")
	config.Redis.Port = viper.GetString("REDIS_PORT")
	config.Redis.Password = viper.GetString("REDIS_PASSWORD")
	config.Redis.DB = viper.GetInt("REDIS_DB")

	// SMTP 配置
	config.SMTP.Host = viper.GetString("SMTP_HOST")
	config.SMTP.Port = viper.GetInt("SMTP_PORT")
	config.SMTP.Username = viper.GetString("SMTP_USERNAME")
	config.SMTP.Password = viper.GetString("SMTP_PASSWORD")

	return config, nil
} 