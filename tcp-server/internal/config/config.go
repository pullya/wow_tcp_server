package config

import log "github.com/sirupsen/logrus"

const (
	TcpPort     = ":8081"      // Номер порта для соединения с tcp-сервером
	ServiceName = "tcp-server" // Имя сервиса для отображения в логах

	PowDifficulty = 10 // Условие сложности для Proof of work
	ProofString   = "Find a string that, when hashed, can be proofed"

	LogLevel = log.DebugLevel // Уровень логирования
)
