package main

import (
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"

	"app/services/worker/internal/consumer"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	rabbitURL := os.Getenv("RABBIT_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}
	queue := os.Getenv("JOB_QUEUE")
	if queue == "" {
		queue = "task_jobs"
	}
	dlqQueue := os.Getenv("DLQ_QUEUE")
	if dlqQueue == "" {
		dlqQueue = "task_jobs_dlq"
	}

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.WithError(err).Fatal("failed to connect to RabbitMQ")
	}
	defer conn.Close()

	jobConsumer := consumer.NewJobConsumer(conn, queue, dlqQueue, log)

	go func() {
		if err := jobConsumer.Start(); err != nil {
			log.WithError(err).Fatal("job consumer failed")
		}
	}()

	log.Info("worker started")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down")
}
