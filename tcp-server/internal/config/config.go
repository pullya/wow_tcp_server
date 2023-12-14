package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

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

	envPort        = "WOW_SERVER_PORT"
	envTimeout     = "WOW_SERVER_TIMEOUT"
	envServiceName = "WOW_SERVER_SERVICE_NAME"
	envDifficulty  = "WOW_SERVER_DIFFICULTY"
	envProofString = "WOW_SERVER_PROOF_STRING"
	envLogLevel    = "WOW_SERVER_LOG_LEVEL"

	shardsCount = 8
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
	envTimeout,
	envServiceName,
	envDifficulty,
	envProofString,
	envLogLevel,
}

type LogLevel string

type Configuration struct {
	Port        int    `yaml:"port"`
	Timeout     int    `yaml:"timeout"`
	ServiceName string `yaml:"serviceName"`

	Difficulty  int    `yaml:"difficulty"`
	ProofString string `yaml:"proofString"`

	ShardsCnt int `yaml:"shardsCnt"`

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
	Config.ShardsCnt = shardsCount

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
					log.Debugf("tcpPort set to %d", Config.Port)
				}
			case envTimeout:
				timeout, err := validateTimeout(envVal)
				if err == nil {
					Config.Timeout = timeout
					log.Debugf("timeout set to %d", Config.Timeout)
				}
			case envServiceName:
				Config.ServiceName = envVal
				log.Debugf("serviceName set to '%s'", Config.ServiceName)
			case envDifficulty:
				diff, err := validateDifficulty(envVal)
				if err == nil {
					Config.Difficulty = diff
					log.Debugf("difficulty set to %d", Config.Difficulty)
				}
			case envProofString:
				Config.ProofString = envVal
				log.Debugf("proofString set to '%s'", Config.ProofString)
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

func validateTimeout(in string) (int, error) {
	num, err := strconv.Atoi(in)
	if err != nil {
		return 0, err
	}
	if num < 0 {
		return 0, errors.New("incorrect timeout")
	}
	return num, nil
}

func validateDifficulty(in string) (int, error) {
	num, err := strconv.Atoi(in)
	if err != nil {
		return 0, err
	}
	if num < 0 || num > 256 {
		return 0, errors.New("incorrect difficulty")
	}
	return num, nil
}

func validateLogLevel(in string) (LogLevel, error) {
	if _, ok := logLevelsMap[LogLevel(in)]; !ok {
		return "", errors.New("incorrect log level")
	}
	return LogLevel(in), nil
}

func BuildPort(port int) string {
	return fmt.Sprintf(":%d", port)
}
