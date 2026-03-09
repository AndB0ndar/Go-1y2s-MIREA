package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func New(service string) *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
	})
	log.SetOutput(os.Stdout)
	level := os.Getenv("LOG_LEVEL")
	if level != "" {
		if lvl, err := logrus.ParseLevel(level); err == nil {
			log.SetLevel(lvl)
		}
	}
	log = log.WithField("service", service).Logger
	return log
}
