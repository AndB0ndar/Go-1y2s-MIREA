package consumer

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
	log     *logrus.Logger
}

func NewConsumer(amqpURL, queue string, log *logrus.Logger) (*Consumer, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare a durable queue
	_, err = ch.QueueDeclare(
		queue,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	// Set prefetch count to 1 to handle one message at a time
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to set Qos: %w", err)
	}

	return &Consumer{
		conn:    conn,
		channel: ch,
		queue:   queue,
		log:     log,
	}, nil
}

func (c *Consumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Consumer) Consume() {
	msgs, err := c.channel.Consume(
		c.queue,
		"",    // consumer
		false, // auto-ack (we'll ack manually)
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // args
	)
	if err != nil {
		c.log.WithError(err).Fatal("failed to register consumer")
	}

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			c.log.WithFields(logrus.Fields{
				"message_id":  msg.MessageId,
				"routing_key": msg.RoutingKey,
			}).Info("received message")

			// Process message (just log its content)
			var event map[string]interface{}
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				c.log.WithError(err).Error("failed to unmarshal message")
				// Reject and don't requeue (poison message)
				msg.Nack(false, false)
				continue
			}

			c.log.WithField("event", event).Info("event data")

			// Acknowledge successful processing
			if err := msg.Ack(false); err != nil {
				c.log.WithError(err).Error("failed to ack message")
			} else {
				c.log.Info("message acknowledged")
			}
		}
	}()

	c.log.Info("waiting for messages. To exit press CTRL+C")
	<-forever
}
