package app

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/pkg/errors"
	"github.com/pullya/wow_tcp_server/tcp-server/internal/config"
	"github.com/pullya/wow_tcp_server/tcp-server/internal/model"
	"github.com/pullya/wow_tcp_server/tcp-server/internal/server"
	"github.com/pullya/wow_tcp_server/tcp-server/internal/storage"
)

var connCnt = 0

type Apper interface {
	doProofOfWork(ctx context.Context, id int) error
}

type App struct {
	server    server.ServerProvider
	storage   storage.Storageer
	challenge Challenger
}

func New(tcpServer server.ServerProvider, storage storage.Storageer, challenge Challenger) App {
	return App{
		server:    tcpServer,
		storage:   storage,
		challenge: challenge,
	}
}

func (a *App) Run(ctx context.Context) error {
	listener, err := a.server.Run(ctx)
	if err != nil {
		config.Logger.Errorf("Error while starting TCP-server: %v", err)
		return err
	}
	defer listener.Close()

	config.Logger.Debug("Waiting for connections...")

	for {
		select {
		case <-ctx.Done():
			config.Logger.Warn("Server shutting down...")
			return nil
		default:
			conn, err := listener.Accept()
			if err != nil {
				return err
			}
			connCnt++

			config.Logger.WithField("connection", connCnt).Debug("New connection established!")
			go a.handleConnection(ctx, conn, connCnt)
		}
	}
}

func (a *App) handleConnection(ctx context.Context, conn net.Conn, id int) {
	defer conn.Close()

	if err := conn.SetReadDeadline(time.Now().Add(a.server.GetTimeout())); err != nil {
		config.Logger.WithField("connection", id).Errorf("Error while setting timeout: %v", err)
		return
	}
	if err := a.doProofOfWork(ctx, conn, id); err != nil {
		return
	}

	wow := a.storage.GetRandomWOW(ctx)
	wowMessage := model.PrepareMessage(model.MessageTypeWow, wow, 0)

	config.Logger.WithField("connection", id).Debugf("Prepared response: %s", string(wow))

	if err := a.server.SendMessage(ctx, conn, wowMessage.AsJsonString()); err != nil {
		config.Logger.WithField("connection", id).Errorf("Error while sending response: %v", err)
		return
	}
}

func (a *App) doProofOfWork(ctx context.Context, conn net.Conn, id int) error {
	challengeMessage := model.PrepareMessage(model.MessageTypeChallenge, generatePOWChallenge(id), a.challenge.Difficulty())

	if err := a.server.SendMessage(ctx, conn, challengeMessage.AsJsonString()); err != nil {
		config.Logger.WithField("connection", id).Errorf("Error while sending challenge: %v", err)
		return err
	}

	response, err := a.server.ReceiveMessage(ctx, conn)
	if err != nil {
		config.Logger.WithField("connection", id).Errorf("Error reading response: %v", err)
		return err
	}

	config.Logger.WithField("connection", id).Infof("Client response received: %s", string(response))

	clientResponse, err := model.ParseServerMessage(response)
	if err != nil {
		config.Logger.WithField("connection", id).Errorf("Unable to unmarshal client message: %v\n", err)
		return err
	}

	solution, err := clientResponse.GetUint64()
	if err != nil {
		config.Logger.WithField("connection", id).Error("Unable to parse solution. Closing connection")
		return errors.New("unable to parse solution")
	}

	if !a.challenge.IsValid(generatePOWChallenge(id), solution) {
		config.Logger.WithField("connection", id).Error("PoW verification failed. Closing connection")
		return errors.New("pow verification failed")
	}

	config.Logger.WithField("connection", id).Debug("PoW verification successful. Allowing connection")

	return nil
}

func generatePOWChallenge(cnt int) string {
	return fmt.Sprintf("%s %d", config.Config.ProofString, cnt)
}
