// Package pubsub provides utility for publishing or subscribing message queue.
// Currently implemented for RabbitMQ.
//
// This package is initiated by:
// - Saggaf Arsyad <saggaf@nusantarabetastudio.com>
// - Muhammad Tri Wibowo <triwibowo@nusantarabetastudio.com>
//
// The reason why publisher and subscriber is separated:
// It's advisable to use separate connections for
// Channel.Publish and Channel.Consume so not to have TCP pushback on publishing
// affect the ability to consume messages, so this parameter is here mostly for
// completeness.
package pubsub

import (
	"github.com/streadway/amqp"
)

type SubscriberCallback func(data []byte, args ...interface{})
type SubscriberCallbackV2 func(m amqp.Delivery, args ...interface{})

type Publisher interface {
	DelayedPublish(queueName string, payload interface{}, delay int64) error
	DelayedPublishJSON(queueName string, payload interface{}, delay int64) error
	Publish(queueName string, payload interface{}) error
	PublishJSON(queueName string, payload interface{}) error
}

type Subscriber interface {
	Subscribe(queueName string, kind string, fn SubscriberCallback, args ...interface{}) error
	SubscribeWithDelivery(numberOfInstance int, queueName string, kind string, fn SubscriberCallbackV2, args ...interface{}) error
	Close() error
}

// NewRabbitPublisher initiate publisher that uses RabbitMQ as datasource
func NewRabbitPublisher(url string) Publisher {
	// Initiate connection to rabbit mq
	conn := NewRabbitConnection(url)
	// Init messenger
	p := RabbitPublisher{conn}
	// Return messenger
	return &p
}

// NewRabbitSubscriber initiate subscriber that uses RabbitMQ as datasource
func NewRabbitSubscriber(url string) Subscriber {
	// Initiate connection to rabbit mq
	conn := NewRabbitConnection(url)
	// Init messenger
	m := RabbitSubscriber{conn}
	// Return messenger
	return &m
}
