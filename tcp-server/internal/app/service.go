package app

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/pullya/wow_tcp_server/tcp-server/internal/config"
	"github.com/pullya/wow_tcp_server/tcp-server/internal/model"
	"github.com/pullya/wow_tcp_server/tcp-server/internal/server"
	"github.com/pullya/wow_tcp_server/tcp-server/internal/storage"
	log "github.com/sirupsen/logrus"
)

var connCnt = 0

type IWowService interface {
	doProofOfWork(ctx context.Context, id int) error
}

type WowService struct {
	Server    server.IServer
	Storage   storage.IStorage
	Challenge IChallenger
}

func NewWowService(tcpServer server.IServer, storage storage.IStorage, challenge IChallenger) WowService {
	return WowService{
		Server:    tcpServer,
		Storage:   storage,
		Challenge: challenge,
	}
}

func (ws *WowService) Run(ctx context.Context) error {
	listener, err := ws.Server.RunServer(ctx)
	if err != nil {
		log.WithField("service", config.ServiceName).Errorf("Error while starting TCP-server: %v", err)
		return err
	}
	defer listener.Close()

	log.WithField("service", config.ServiceName).Debug("Waiting for connections...")

	for {
		select {
		case <-ctx.Done():
			log.WithField("service", config.ServiceName).Warn("Server shutting down...")
			return nil
		default:
			conn, err := listener.Accept()
			if err != nil {
				return err
			}
			connCnt++
			ws.Server.SetConn(conn)

			log.WithFields(log.Fields{"service": config.ServiceName, "connection": connCnt}).Debug("New connection established!")

			go ws.HandleConnection(ctx, connCnt)
		}
	}
}

func (ws *WowService) HandleConnection(ctx context.Context, id int) {
	defer ws.Server.CloseConn()

	if err := ws.doProofOfWork(ctx, id); err != nil {
		log.WithFields(log.Fields{"service": config.ServiceName, "connection": id}).Errorf("Proof of work failed: %v", err)
		return
	}

	wow := ws.Storage.GetRandomWoW(ctx)
	wowMessage := model.PrepareMessage(model.MessageTypeWow, wow, 0)

	log.WithFields(log.Fields{"service": config.ServiceName, "connection": id}).Debugf("Prepared response: %s", string(wow))

	if err := ws.Server.SendMessage(ctx, wowMessage.AsJsonString()); err != nil {
		log.WithFields(log.Fields{"service": config.ServiceName, "connection": id}).Errorf("Error while sending response: %v", err)
		return
	}
}

func (ws *WowService) doProofOfWork(ctx context.Context, id int) error {
	challengeMessage := model.PrepareMessage(model.MessageTypeChallenge, generatePoWChallenge(id), ws.Challenge.GetPowDifficulty())

	if err := ws.Server.SendMessage(ctx, challengeMessage.AsJsonString()); err != nil {
		log.WithFields(log.Fields{"service": config.ServiceName, "connection": id}).Errorf("Error while sending challenge: %v", err)
		return err
	}

	response, err := ws.Server.ReceiveMessage(ctx)
	if err != nil {
		log.WithFields(log.Fields{"service": config.ServiceName, "connection": id}).Errorf("Error reading response: %v", err)
		return err
	}

	log.WithFields(log.Fields{"service": config.ServiceName, "connection": id}).Infof("Client response received: %s", string(response))

	clientResponse, err := model.ParseServerMessage(response)
	if err != nil {
		log.WithFields(log.Fields{"service": config.ServiceName, "connection": id}).
			Errorf("Unable to unmarshal client message: %v\n", err)
		return err
	}

	solution, err := clientResponse.GetUint64()
	if err != nil {
		log.WithFields(log.Fields{"service": config.ServiceName, "connection": id}).Error("Unable to parse solution. Closing connection")
		return errors.New("unable to parse solution")
	}

	if !ws.Challenge.IsValidPoW(generatePoWChallenge(id), solution) {
		log.WithFields(log.Fields{"service": config.ServiceName, "connection": id}).Error("PoW verification failed. Closing connection")
		return errors.New("pow verification failed")
	}

	log.WithFields(log.Fields{"service": config.ServiceName, "connection": id}).Debug("PoW verification successful. Allowing connection")

	return nil
}

func generatePoWChallenge(cnt int) string {
	return fmt.Sprintf("%s %d", config.ProofString, cnt)
}
