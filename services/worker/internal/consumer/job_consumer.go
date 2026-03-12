package consumer

import (
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"

	"app/services/worker/internal/storage"
)

const maxAttempts = 3

type JobConsumer struct {
	conn      *amqp.Connection
	queue     string
	dlqQueue  string
	log       *logrus.Logger
	processed *storage.ProcessedStore
}

func NewJobConsumer(conn *amqp.Connection, queue, dlqQueue string, log *logrus.Logger) *JobConsumer {
	return &JobConsumer{
		conn:      conn,
		queue:     queue,
		dlqQueue:  dlqQueue,
		log:       log,
		processed: storage.NewProcessedStore(),
	}
}

func (c *JobConsumer) Start() error {
	ch, err := c.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()

	// Объявляем основную очередь и DLQ (durable)
	_, err = ch.QueueDeclare(c.queue, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}
	_, err = ch.QueueDeclare(c.dlqQueue, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare DLQ: %w", err)
	}

	// Prefetch 1
	err = ch.Qos(1, 0, false)
	if err != nil {
		return fmt.Errorf("failed to set prefetch: %w", err)
	}

	msgs, err := ch.Consume(c.queue, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to consume: %w", err)
	}

	c.log.WithField("queue", c.queue).Info("started consuming jobs")

	for msg := range msgs {
		c.processMessage(msg)
	}
	return nil
}

type jobMessage struct {
	Job       string `json:"job"`
	TaskID    string `json:"task_id"`
	Attempt   int    `json:"attempt"`
	MessageID string `json:"message_id"`
	Timestamp string `json:"timestamp"`
}

func (c *JobConsumer) processMessage(msg amqp.Delivery) {
	var job jobMessage
	if err := json.Unmarshal(msg.Body, &job); err != nil {
		c.log.WithError(err).Error("invalid job message, sending to DLQ")
		// Отправляем в DLQ
		c.sendToDLQ(msg.Body, "invalid json")
		msg.Ack(false) // убираем из основной очереди
		return
	}

	logger := c.log.WithFields(logrus.Fields{
		"message_id": job.MessageID,
		"task_id":    job.TaskID,
		"attempt":    job.Attempt,
	})

	// Проверка идемпотентности
	if c.processed.Exists(job.MessageID) {
		logger.Info("duplicate message, skipping")
		msg.Ack(false)
		return
	}

	// Имитация обработки (может упасть с ошибкой)
	err := c.doWork(job)
	if err != nil {
		logger.WithError(err).Error("work failed")

		if job.Attempt < maxAttempts {
			// Retry: публикуем новое сообщение с увеличенным attempt
			job.Attempt++
			newBody, _ := json.Marshal(job)
			c.publishRetry(newBody)
			logger.WithField("new_attempt", job.Attempt).Info("retrying job")
			msg.Ack(false) // исходное удаляем
		} else {
			// Превышены попытки -> DLQ
			logger.Error("max attempts exceeded, sending to DLQ")
			c.sendToDLQ(msg.Body, "max attempts exceeded")
			msg.Ack(false)
		}
		return
	}

	// Успех
	logger.Info("job processed successfully")
	c.processed.Add(job.MessageID)
	msg.Ack(false)
}

// doWork имитирует выполнение задачи
func (c *JobConsumer) doWork(job jobMessage) error {
	// Имитация длительности
	time.Sleep(2 * time.Second)

	// Для демонстрации ошибок: если task_id содержит "fail", генерируем ошибку
	if job.TaskID == "fail" || job.TaskID == "t_fail" {
		return fmt.Errorf("simulated processing error")
	}
	return nil
}

func (c *JobConsumer) publishRetry(body []byte) {
	ch, err := c.conn.Channel()
	if err != nil {
		c.log.WithError(err).Error("failed to open channel for retry")
		return
	}
	defer ch.Close()
	err = ch.Publish("", c.queue, true, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         body,
		DeliveryMode: amqp.Persistent,
	})
	if err != nil {
		c.log.WithError(err).Error("failed to publish retry message")
	}
}

func (c *JobConsumer) sendToDLQ(body []byte, reason string) {
	ch, err := c.conn.Channel()
	if err != nil {
		c.log.WithError(err).Error("failed to open channel for DLQ")
		return
	}
	defer ch.Close()
	err = ch.Publish("", c.dlqQueue, true, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         body,
		DeliveryMode: amqp.Persistent,
		Headers: amqp.Table{
			"dlq-reason": reason,
		},
	})
	if err != nil {
		c.log.WithError(err).Error("failed to publish to DLQ")
	}
}
