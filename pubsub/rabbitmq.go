package pubsub

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const (
	ExchangeTypeDelayedMessage = "x-delayed-message"
	NoExchange                 = ""
)

type RabbitDatasource struct {
	*amqp.Connection
}

// @todo separate implementation in sub-package rabbit
type RabbitSubscriber struct {
	conn RabbitDatasource
}

// SubscribeWithDelivery subcribe with automatic ack
func (rs *RabbitSubscriber) Subscribe(queueName string, kind string, fn SubscriberCallback, args ...interface{}) error {
	// Open channel
	conn := rs.conn
	ch, err := conn.Channel()
	if err != nil {
		log.Printf(" Error %s", err)
	}
	// Close channel on returns
	defer ch.Close()
	// Declare queue if not exist
	queue, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Error(err)
		return err
	}
	// Init consumed queue
	var consumedQueue string
	// If type delayed message, create exchanger and temp queue
	if kind == ExchangeTypeDelayedMessage {
		// Generate exchange name
		exchangeName := getDelayedExchangeName(queueName)
		// Declare delayed exchange
		err := declareDelayedExchange(ch, exchangeName)
		// Declare random queue if not exist, this queue will not survive rabbit restart
		randomQueue, err := ch.QueueDeclare("", false, false, false, false, nil)
		if err != nil {
			log.Error(err)
			return err
		}
		// Set QoS
		err = ch.Qos(1, 0, false)
		// Bind queue
		err = ch.QueueBind(randomQueue.Name, queueName, exchangeName, false, nil)
		// Set consumed queue
		consumedQueue = randomQueue.Name
	} else {
		consumedQueue = queue.Name
	}
	// Start consuming queue
	messages, err := ch.Consume(consumedQueue, "", false, false, false, false, nil)
	go func() {
		for d := range messages {
			log.Debugf("Received message for queue %s", queueName)
			if fn != nil {
				// Execute callback
				fn(d.Body, args...)
			}
			// Send ack
			err = d.Ack(false)
			if err != nil {
				log.Errorf("[rabbitmq] error when acknowledge message for queue: %s -> %s", queueName, err.Error())
			}
		}
	}()
	// Reading from an empty channel to block function from returning
	// Source: https://blog.sgmansfield.com/2016/06/how-to-block-forever-in-go/
	log.Infof("%s receiver initiated", queueName)
	forever := make(chan bool)
	<-forever
	// Return success
	return nil
}

// SubscribeWithDelivery subcribe with manual ack or nack
func (rs *RabbitSubscriber) SubscribeWithDelivery(numberOfInstance int, queueName string, kind string, fn SubscriberCallbackV2, args ...interface{}) error {
	// Open channel
	conn := rs.conn
	ch, err := conn.Channel()
	if err != nil {
		log.Printf(" Error %s", err)
	}
	// Close channel on returns
	defer ch.Close()
	// Declare queue if not exist
	queue, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Error(err)
		return err
	}
	// Init consumed queue
	var consumedQueue string
	// If type delayed message, create exchanger and temp queue
	if kind == ExchangeTypeDelayedMessage {
		// Generate exchange name
		exchangeName := getDelayedExchangeName(queueName)
		// Declare delayed exchange
		err := declareDelayedExchange(ch, exchangeName)
		// Declare random queue if not exist, this queue will not survive rabbit restart
		randomQueue, err := ch.QueueDeclare("", false, false, false, false, nil)
		if err != nil {
			log.Error(err)
			return err
		}
		// Set QoS
		err = ch.Qos(1, 0, false)
		// Bind queue
		err = ch.QueueBind(randomQueue.Name, queueName, exchangeName, false, nil)
		// Set consumed queue
		consumedQueue = randomQueue.Name
	} else {
		consumedQueue = queue.Name
	}
	// Start consuming queue

	// default if one instance
	if numberOfInstance == 0 {
		numberOfInstance = 1
	}
	for i := 0; i < numberOfInstance; i++ {
		consumerName := fmt.Sprintf("instance-%d", i+1)
		messages, _ := ch.Consume(consumedQueue, consumerName, false, false, false, false, nil)
		go func() {
			for d := range messages {
				log.Debugf("Received message for queue %s", queueName)
				if fn != nil {
					// Execute callback
					fn(d, args...)
				}
			}
		}()
	}
	// Reading from an empty channel to block function from returning
	// Source: https://blog.sgmansfield.com/2016/06/how-to-block-forever-in-go/
	log.Infof("%s receiver initiated", queueName)
	forever := make(chan bool)
	<-forever
	// Return success
	return nil
}

func (rs *RabbitSubscriber) Close() error {
	conn := rs.conn
	err := conn.Close()
	if err != nil {
		log.Errorf("error when shutdown subscriber, err: %s", err.Error())
	}

	return err
}

// @todo separate implementation in sub-package rabbit
type RabbitPublisher struct {
	conn RabbitDatasource
}

func (rp *RabbitPublisher) DelayedPublish(queueName string, payload interface{}, delay int64) error {
	// Get connection
	conn := rp.conn
	// Open channel
	ch, err := conn.Channel()
	if err != nil {
		log.Error(err)
		return err
	}
	// Close channel on return
	defer ch.Close()
	// Create payload headers
	headers := make(amqp.Table)
	headers["x-delay"] = delay
	// Encode payload
	bodyBytes, err := EncodeTextPlain(payload)
	if err != nil {
		return err
	}
	p := amqp.Publishing{
		Headers:      headers,
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         bodyBytes,
	}
	// Exchange name
	exchangeName := getDelayedExchangeName(queueName)
	err = ch.Publish(exchangeName, queueName, false, false, p)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Debugf("Message sent to queue %s", queueName)
	return nil
}

func (rp *RabbitPublisher) DelayedPublishJSON(queueName string, payload interface{}, delay int64) error {
	// Get connection
	conn := rp.conn
	// Open channel
	ch, err := conn.Channel()
	if err != nil {
		log.Error(err)
		return err
	}
	// Close channel on return
	defer ch.Close()
	// Create payload headers
	headers := make(amqp.Table)
	headers["x-delay"] = delay
	// Encode payload
	bodyBytes, err := EncodeJSON(payload)
	// log.Debugf("[DelayedPublishJSON] body bytes %s", string(bodyBytes))
	if err != nil {
		return err
	}
	p := amqp.Publishing{
		Headers:      headers,
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         bodyBytes,
	}
	// Exchange name
	exchangeName := getDelayedExchangeName(queueName)
	err = ch.Publish(exchangeName, queueName, false, false, p)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Debugf("Message sent to queue %s", queueName)
	return nil
}

func (rp *RabbitPublisher) Publish(queueName string, payload interface{}) error {
	// Get connection
	conn := rp.conn
	// Open channel
	ch, err := conn.Channel()
	if err != nil {
		log.Error(err)
		return err
	}
	log.Debugf("Message open a channel to queue %s", queueName)
	// Close channel on return
	defer ch.Close()
	// Encode payload
	bodyBytes, err := EncodeTextPlain(payload)
	if err != nil {
		return err
	}
	log.Debugf("Message Encode text for queue %s", queueName)
	// Publish
	err = ch.Publish(NoExchange, queueName, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         bodyBytes,
	})
	if err != nil {
		log.Error(err)
		return err
	}
	log.Debugf("Message sent to queue %s", queueName)
	return nil
}

func (rp *RabbitPublisher) PublishJSON(queueName string, payload interface{}) error {
	// Get connection
	conn := rp.conn
	// Open channel
	ch, err := conn.Channel()
	if err != nil {
		log.Error(err)
		return err
	}
	log.Debugf("Message open a channel to queue %s", queueName)
	// Close channel on return
	defer ch.Close()
	// Encode payload
	bodyBytes, err := EncodeJSON(payload)
	if err != nil {
		return err
	}
	log.Debugf("Message Encode json for queue %s", queueName)
	// Publish
	err = ch.Publish(NoExchange, queueName, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         bodyBytes,
	})
	if err != nil {
		log.Error(err)
		return err
	}
	log.Debugf("Message sent to queue %s", queueName)
	return nil
}

// NewRabbitConnection initiate rabbit connection and crash if fails
func NewRabbitConnection(url string) RabbitDatasource {
	// Initiate connection
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Errorf("failed to connect to RabbitMQ Datasource. URL: %s, Error: %s", url, err.Error())
		os.Exit(3)
	}
	log.Info("RabbitMQ Connection Started")
	// Return DataSource interface
	return RabbitDatasource{conn}
}

// getDelayedExchangeName returns delayed exchange name pattern based on queue name
func getDelayedExchangeName(queueName string) string {
	return queueName + ".delayed"
}

// declareDelayedExchange makes sure that delayed exchange is declared in RabbitMQ
func declareDelayedExchange(ch *amqp.Channel, name string) error {
	// Add arguments
	exchangeArgs := make(amqp.Table)
	exchangeArgs["x-delayed-type"] = "direct"
	// Declare delayed exchange
	err := ch.ExchangeDeclare(name, ExchangeTypeDelayedMessage, true, false, false, false, exchangeArgs)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// EncodeTextPlain encodes struct into text/plain body
func EncodeTextPlain(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(data)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return buf.Bytes(), nil
}

// EncodeJSON encodes struct into application/json body
func EncodeJSON(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(data)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return buf.Bytes(), nil
}

// DecodeTextPlain decodes bytes of text/plain into struct
func DecodeTextPlain(bts []byte, data interface{}) error {
	buf := bytes.NewBuffer(bts)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(data)
	if err != nil {
		return err
	}
	return nil
}

// DecodeJSON decodes bytes of json into struct
func DecodeJSON(bts []byte, data interface{}) error {
	buf := bytes.NewBuffer(bts)
	dec := json.NewDecoder(buf)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	return nil
}
