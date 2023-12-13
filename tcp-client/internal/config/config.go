package config

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	configFile = "config.yaml"

	logLevelDebug = "Debug"
	logLevelInfo  = "Info"
	logLevelWarn  = "Warn"
	logLevelError = "Error"
	logLevelFatal = "Fatal"
	logLevelPanic = "Panic"
	logLevelTrace = "Trace"
)

var Config Configuration

var logLevelsMap = map[LogLevel]log.Level{
	logLevelDebug: log.DebugLevel,
	logLevelInfo:  log.InfoLevel,
	logLevelWarn:  log.WarnLevel,
	logLevelError: log.ErrorLevel,
	logLevelFatal: log.FatalLevel,
	logLevelPanic: log.PanicLevel,
	logLevelTrace: log.TraceLevel,
}

type LogLevel string

type Configuration struct {
	Address     string `yaml:"address"`
	ServiceName string `yaml:"serviceName"`

	ClientsCount int           `yaml:"clientsCount"`
	ConnInterval time.Duration `yaml:"connInterval"`

	LogLevel LogLevel `yaml:"logLevel"`
}

func ReadConfig() {
	log.SetLevel(log.DebugLevel)
	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("failed to read configuration from file '%s', error: %v", configFile, err)
		return
	}

	err = yaml.Unmarshal(data, &Config)
	if err != nil {
		log.Fatalf("failed to unmarshall config file '%s', error: %v", configFile, err)
		return
	}

	log.Debugf("Configuration read: %v", Config)
}

func (l LogLevel) ToLogrusFormat() log.Level {
	res, ok := logLevelsMap[l]
	if !ok {
		return log.ErrorLevel
	}
	return res
}
