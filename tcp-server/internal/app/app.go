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
	sendChallenge(ctx context.Context, conn net.Conn, id int) error
	validatePOW(ctx context.Context, clientResponse model.Message, id int) error
	sendWOW(ctx context.Context, conn net.Conn, uid string, id int) error
}

type App struct {
	server       server.ServerProvider
	storage      storage.Storageer
	requeststore storage.Requester
	challenge    Challenger
}

func New(tcpServer server.ServerProvider, storage storage.Storageer, requeststore storage.Requester, challenge Challenger) App {
	return App{
		server:       tcpServer,
		storage:      storage,
		requeststore: requeststore,
		challenge:    challenge,
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

	request, err := a.server.ReceiveMessage(ctx, conn)
	if err != nil {
		config.Logger.WithField("connection", id).Errorf("Error reading request: %v", err)
		return
	}

	config.Logger.WithField("connection", id).Debugf("Request from client received: %s", request)

	clientRequest, err := model.ParseServerMessage(request)
	if err != nil {
		config.Logger.WithField("connection", id).Errorf("Unable to unmarshal client message: %v\n", err)
		return
	}

	switch clientRequest.MessageType {
	case model.MessageTypeRequest:
		if err := a.sendChallenge(ctx, conn, id); err != nil {
			config.Logger.WithField("connection", id).Errorf("Error while sending challenge: %v", err)
			return
		}
	case model.MessageTypeSolution:
		if err := a.validatePOW(ctx, clientRequest, id); err != nil {
			config.Logger.WithField("connection", id).Errorf("Failed to validate POW: %v", err)
			return
		}
		if err = a.sendWOW(ctx, conn, clientRequest.RequestID, id); err != nil {
			return
		}
	default:
		config.Logger.WithField("connection", id).Errorf("Unknown message type: %s", clientRequest.MessageType)
		return
	}
}

func (a *App) sendChallenge(ctx context.Context, conn net.Conn, id int) error {
	uid := storage.GenUID()
	challengeMessage := model.PrepareMessage(uid, model.MessageTypeChallenge, generatePOWChallenge(uid), a.challenge.Difficulty())

	if err := a.server.SendMessage(ctx, conn, challengeMessage.AsJsonString()); err != nil {
		return err
	}

	a.requeststore.Add(ctx, uid)

	return nil
}

func (a *App) validatePOW(ctx context.Context, clientResponse model.Message, id int) error {
	ok, err := a.requeststore.Get(ctx, clientResponse.RequestID)
	if err != nil {
		config.Logger.WithField("connection", id).Errorf("Failed to find request '%s' in store: %v", clientResponse.RequestID, err)
		return err
	}
	if ok {
		config.Logger.WithField("connection", id).Errorf("This POW was already handled '%s'", clientResponse.RequestID)
		return errors.New("Double work")
	}

	solution, err := clientResponse.GetUint64()
	if err != nil {
		config.Logger.WithField("connection", id).Error("Unable to parse solution. Closing connection")
		return errors.New("unable to parse solution")
	}

	if !a.challenge.IsValid(generatePOWChallenge(clientResponse.RequestID), solution) {
		config.Logger.WithField("connection", id).Error("PoW verification failed. Closing connection")
		return errors.New("pow verification failed")
	}

	config.Logger.WithField("connection", id).Debug("PoW verification successful. Allowing connection")

	return nil
}

func (a *App) sendWOW(ctx context.Context, conn net.Conn, uid string, id int) error {
	wow := a.storage.GetRandomWOW(ctx)
	wowMessage := model.PrepareMessage(uid, model.MessageTypeWow, wow, 0)

	config.Logger.WithField("connection", id).Debugf("Prepared response: %s", string(wow))

	if err := a.server.SendMessage(ctx, conn, wowMessage.AsJsonString()); err != nil {
		config.Logger.WithField("connection", id).Errorf("Error while sending response: %v", err)
		return err
	}

	if err := a.requeststore.Set(ctx, uid); err != nil {
		config.Logger.WithField("connection", id).Errorf("Failed to set status for request '%s': %v", uid, err)
		return err
	}

	return nil
}

func generatePOWChallenge(cnt string) string {
	return fmt.Sprintf("%s %s", config.Config.ProofString, cnt)
}
