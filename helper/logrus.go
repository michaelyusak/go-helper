package helper

import "github.com/sirupsen/logrus"

func NewLogrus() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		QuoteEmptyFields: true,
	})
	return log
}
