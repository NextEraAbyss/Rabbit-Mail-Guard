package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"Rabbit-Mail-Guard/config"
	"Rabbit-Mail-Guard/internal/email"
	"Rabbit-Mail-Guard/internal/redis"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EmailConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
	dlQueue string // 死信队列
	email   *email.Service
	redis   *redis.Client
}

type EmailMessage struct {
	To string `json:"to"`
}

func NewEmailConsumer(cfg *config.Config, emailSvc *email.Service, redisCli *redis.Client) (*EmailConsumer, error) {
	// 使用最基础的连接方式
	url := fmt.Sprintf("amqp://guest:guest@%s:%s/",
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
	)
	log.Printf("Attempting to connect to RabbitMQ at: %s", url)
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Failed to open channel: %v", err)
		conn.Close()
		return nil, err
	}

	// 声明队列
	_, err = ch.QueueDeclare(
		cfg.RabbitMQ.Queue, // name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		log.Printf("Failed to declare queue: %v", err)
		ch.Close()
		conn.Close()
		return nil, err
	}

	// 声明死信队列
	dlQueue := cfg.RabbitMQ.Queue + "_failed"
	_, err = ch.QueueDeclare(
		dlQueue, // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Printf("Failed to declare dead letter queue: %v", err)
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &EmailConsumer{
		conn:    conn,
		channel: ch,
		queue:   cfg.RabbitMQ.Queue,
		dlQueue: dlQueue,
		email:   emailSvc,
		redis:   redisCli,
	}, nil
}

func (c *EmailConsumer) Start() error {
	msgs, err := c.channel.Consume(
		c.queue,
		"",
		false, // manual ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var msg EmailMessage
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				c.moveToDeadLetter(d.Body)
				d.Ack(false)
				continue
			}

			code, err := c.email.SendVerificationCode(msg.To)
			if err != nil {
				log.Printf("Error sending verification code: %v", err)
				c.moveToDeadLetter(d.Body)
				d.Ack(false)
				continue
			}

			if err := c.redis.SetVerificationCode(context.Background(), msg.To, code); err != nil {
				log.Printf("Error storing verification code: %v", err)
				c.moveToDeadLetter(d.Body)
				d.Ack(false)
				continue
			}

			log.Printf("Sent verification code to %s", msg.To)
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

	return nil
}

// 添加移动到死信队列的方法
func (c *EmailConsumer) moveToDeadLetter(body []byte) {
	err := c.channel.Publish(
		"",        // exchange
		c.dlQueue, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Printf("Failed to move message to dead letter queue: %v", err)
	} else {
		log.Printf("Message moved to dead letter queue")
	}
}
