package mq

import (
	"log"
	"tiny-tiktok/service/video_service/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}

func InitMQ() *amqp.Connection {
	url := config.InitRabbitMQUrl()
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")
	return conn
}
