package service

// pubsubService is an implementation of pubsub package

import (
	"fmt"
	"os"

	"github.com/interview-order/flags"
	"github.com/interview-order/pubsub"
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
	url, ok := os.LookupEnv("AMQP_URL")
	if !ok {
		fmt.Println("AMQP_URL is not set in system environment")
		os.Exit(2)
	}
	return url
}
