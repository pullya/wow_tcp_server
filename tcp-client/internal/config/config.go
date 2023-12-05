package config

import log "github.com/sirupsen/logrus"

const (
	Address     = "tcp_server:8081" // Адрес для соединения с tcp-сервером. Имя "tcp_server" должно совпадать с именем сервера в файле docker-compose.yml
	ServiceName = "tcp-client"      // Имя сервиса для отображения в логах

	ClientsCount = 5    // Количество клиентов, которое будет запущено
	ConnInterval = 3000 // Интервал в миллисекундах между запуском горутин с клиентами

	LogLevel = log.DebugLevel // Уровень логирования
)
