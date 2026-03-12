package main

import (
	"os"
	"os/signal"
	"syscall"

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
	queueName := os.Getenv("QUEUE_NAME")
	if queueName == "" {
		queueName = "task_events"
	}

	cons, err := consumer.NewConsumer(rabbitURL, queueName, log)
	if err != nil {
		log.WithError(err).Fatal("failed to create consumer")
	}
	defer cons.Close()

	// Graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		log.Info("shutting down consumer...")
		cons.Close()
		os.Exit(0)
	}()

	log.Info("worker started, consuming messages...")
	cons.Consume()
}
