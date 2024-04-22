package mq

import (
	"tiny-tiktok/api_router/pkg/logger"
	"tiny-tiktok/service/video_service/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

func FailOnError(err error, msg string) {
	if err != nil {
		logger.Log.Errorf("%s: %s", msg, err)
	}
}

func InitMQ() *amqp.Connection {
	url := config.InitRabbitMQUrl()
	conn, err := amqp.Dial(url)
	FailOnError(err, "Failed to connect to RabbitMQ")
	return conn
}
