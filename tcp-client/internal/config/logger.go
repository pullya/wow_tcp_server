package config

import (
	log "github.com/sirupsen/logrus"
)

var Logger *log.Entry

func InitLogger() {
	Logger = log.WithFields(
		log.Fields{
			"service": Config.ServiceName,
		})
}
