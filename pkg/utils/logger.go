package utils

import (
	"github.com/sirupsen/logrus"
	"os"
)

type Logger struct {
	*logrus.Logger
}

func NewLogger(svcName string) *Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)

	log.WithFields(logrus.Fields{"module": svcName})

	return &Logger{log}
}

type ErrorMessage struct {
	Error   string `json:"error"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}
