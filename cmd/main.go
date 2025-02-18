package main

import (
	"log"

	"Rabbit-Mail-Guard/config"
	"Rabbit-Mail-Guard/internal/consumer"
	"Rabbit-Mail-Guard/internal/email"
	"Rabbit-Mail-Guard/internal/redis"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// 初始化 Redis 客户端
	redisCli := redis.NewRedisClient(
		cfg.Redis.Host,
		cfg.Redis.Port,
		cfg.Redis.Password,
		cfg.Redis.DB,
	)

	// 初始化邮件服务
	emailSvc := email.NewEmailService(
		cfg.SMTP.Host,
		cfg.SMTP.Port,
		cfg.SMTP.Username,
		cfg.SMTP.Password,
	)

	// 初始化消费者
	consumer, err := consumer.NewEmailConsumer(
		cfg,
		emailSvc,
		redisCli,
	)
	if err != nil {
		log.Fatalf("Error creating consumer: %v", err)
	}

	// 启动消费者
	if err := consumer.Start(); err != nil {
		log.Fatalf("Error starting consumer: %v", err)
	}
}
