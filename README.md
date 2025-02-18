# Rabbit-Mail-Guard

> 🚀 A robust email verification service powered by RabbitMQ and Redis

## 项目概述
基于 Go 语言开发的邮件验证码服务，提供可靠的验证码发送和管理功能。

### 核心特性
- 🚀 基于 RabbitMQ 的消息队列处理
- 📧 邮件验证码发送
- 📦 Redis 验证码管理
- 🔄 死信队列容错机制

## 技术栈
- 💻 Go 1.21+
- 🐰 RabbitMQ
- 🔥 Redis
- ✉️ SMTP

## 项目结构
```
Rabbit-Mail-Guard/
├── cmd/                    # 应用入口
│   └── main.go            # 主程序
├── config/                 # 配置管理
│   └── config.go          # 配置加载器
├── internal/              # 内部包
│   ├── consumer/          # 消息消费者
│   │   └── email_consumer.go
│   ├── email/            # 邮件服务
│   │   └── service.go
│   └── redis/            # Redis 客户端
│       └── client.go
├── .env                   # 环境配置
└── go.mod                # 依赖管理
```

## 快速开始

### 环境要求
- Go 1.21+
- RabbitMQ
- Redis
- SMTP 服务

### 安装
```bash
# 克隆项目
git clone https://github.com/yourusername/Rabbit-Mail-Guard.git

# 安装依赖
go mod tidy
```

### 配置
创建 `.env` 文件：
```env
# RabbitMQ
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USERNAME=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_QUEUE=email_queue
RABBITMQ_EXCHANGE=email_exchange

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=your_password
REDIS_DB=0

# SMTP
SMTP_HOST=smtp.example.com
SMTP_PORT=465
SMTP_USERNAME=your_email
SMTP_PASSWORD=your_password
```

### 运行
```bash
go run cmd/main.go
```

## 使用说明

### 发送验证码
向 RabbitMQ 队列发送消息：
```json
{
    "to": "user@example.com"
}
```

### 错误处理
服务包含完整的错误处理机制：
- ✅ 消息解析失败 -> 死信队列
- ✅ 邮件发送失败 -> 死信队列
- ✅ Redis 存储失败 -> 死信队列

### 监控指标
- RabbitMQ 连接状态
- 消息处理状态
- 邮件发送结果
- 错误统计

## 最佳实践
1. 定期检查死信队列
2. 监控服务状态
3. 及时处理错误日志
4. 定期更新依赖

## 开发计划
- [ ] 消息重试机制
- [ ] 验证码验证 API
- [ ] 性能监控
- [ ] 并发控制优化

## 贡献指南
欢迎提交 Issue 和 Pull Request

## 许可证
MIT License

## 联系方式
- 作者：NextEraAbyss
- 邮箱：1578347363@qq.com
