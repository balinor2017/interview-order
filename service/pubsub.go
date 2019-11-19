package service

// pubsubService is an implementation of pubsub package

import (
	"github.com/balinor2017/interview-order/config"
	"github.com/balinor2017/interview-order/flags"
	"github.com/balinor2017/interview-order/pubsub"
)

// Publisher connection instance
var publisher pubsub.Publisher

// Subscriber connection instance
var subscriber pubsub.Subscriber

// Init starts connection to RabbitMQ, subscriber and publishers
func InitMQPublisher() {
	// Get amqp url
	url := getRabbitURL()
	// Init RabbitMQ publisher connection
	publisher = pubsub.NewRabbitPublisher(url)
}

func InitMQConsumer() {
	// Get amqp url
	url := getRabbitURL()
	// Init RabbitMQ subscriber connection`
	subscriber = pubsub.NewRabbitSubscriber(url)
	go subscriber.Subscribe(flags.SendingNotification, pubsub.NoExchange, handleSendingNotification)
}

// getRabbitURL returns AMQP connection url and crash if no such env found
func getRabbitURL() string {
	// Get connection url
	url := config.MustGetString("rabbitmq.AMQP_URL")
	return url
}
