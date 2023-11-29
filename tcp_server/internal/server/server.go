package server

import (
	"bufio"
	"context"
	"net"
	"strings"

	log "github.com/rs/zerolog/log"
)

func RunServer(_ context.Context) error {
	log.Info().Msg("Launching server...")

	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		return err
	}

	conn, err := listener.Accept()
	if err != nil {
		return err
	}

	for {
		// Будем прослушивать все сообщения разделенные \n
		message, _ := bufio.NewReader(conn).ReadString('\n')
		// Распечатываем полученое сообщение
		log.Info().Msgf("Message Received: %s", string(message))
		// Процесс выборки для полученной строки
		newmessage := strings.ToUpper(message)
		// Отправить новую строку обратно клиенту
		conn.Write([]byte(newmessage + "\n"))
	}

	//return nil
}
