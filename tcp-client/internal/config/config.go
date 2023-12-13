package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
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

	envPort         = "WOW_CLIENT_PORT"
	envServiceName  = "WOW_CLIENT_SERVICE_NAME"
	envClientsCount = "WOW_CLIENT_CLIENTS_COUNT"
	envConnInterval = "WOW_CLIENT_CONN_INTERVAL"
	envLogLevel     = "WOW_CLIENT_LOG_LEVEL"
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

var envArray = []string{
	envPort,
	envServiceName,
	envClientsCount,
	envConnInterval,
	envLogLevel,
}

type LogLevel string

type Configuration struct {
	Port        int    `yaml:"port"`
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

	log.Debugf("Default configuration read: %v", Config)

	checkEnv()
}

func (l LogLevel) ToLogrusFormat() log.Level {
	res, ok := logLevelsMap[l]
	if !ok {
		return log.ErrorLevel
	}
	return res
}

func checkEnv() {

	for _, name := range envArray {
		if envVal := os.Getenv(name); envVal != "" {
			switch name {
			case envPort:
				port, err := validatePort(envVal)
				if err == nil {
					Config.Port = port
					log.Debugf("address set to %d", Config.Port)
				}
			case envServiceName:
				Config.ServiceName = envVal
				log.Debugf("serviceName set to '%s'", Config.ServiceName)
			case envClientsCount:
				cc, err := validateClientsCount(envVal)
				if err == nil {
					Config.ClientsCount = cc
					log.Debugf("clientsCount set to %d", Config.ClientsCount)
				}
			case envConnInterval:
				ci, err := validateConnInterval(envVal)
				if err == nil {
					Config.ConnInterval = time.Duration(ci)
					log.Debugf("connInterval set to '%s'", Config.ConnInterval)
				}
			case envLogLevel:
				ll, err := validateLogLevel(envVal)
				if err == nil {
					Config.LogLevel = ll
					log.Debugf("logLevel set to '%v'", Config.LogLevel)
				}
			}
		}
	}
}

func validatePort(in string) (int, error) {
	num, err := strconv.Atoi(in)
	if err != nil {
		return 0, err
	}
	if num < 0 || num > 65535 {
		return 0, errors.New("incorrect port number")
	}
	return num, nil
}

func validateClientsCount(in string) (int, error) {
	num, err := strconv.Atoi(in)
	if err != nil {
		return 0, err
	}
	if num < 0 {
		return 0, errors.New("incorrect clients count")
	}
	return num, nil
}

func validateConnInterval(in string) (int, error) {
	num, err := strconv.Atoi(in)
	if err != nil {
		return 0, err
	}
	if num < 0 {
		return 0, errors.New("incorrect conn interval")
	}
	return num, nil
}

func validateLogLevel(in string) (LogLevel, error) {
	if _, ok := logLevelsMap[LogLevel(in)]; !ok {
		return "", errors.New("incorrect log level")
	}
	return LogLevel(in), nil
}

func BuildAddress(port int) string {
	return fmt.Sprintf("tcp_server:%d", port)
}
