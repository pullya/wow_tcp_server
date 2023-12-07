package app

import (
	"context"
	"sync"
	"time"

	"github.com/pullya/wow_tcp_server/tcp-client/internal/client"
	"github.com/pullya/wow_tcp_server/tcp-client/internal/config"
	"github.com/pullya/wow_tcp_server/tcp-client/internal/model"
	log "github.com/sirupsen/logrus"
)

type WowService struct {
	Client    client.IClient
	Challenge IChallenger
	wg        *sync.WaitGroup
}

func NewWowService(client client.IClient, challenge IChallenger) WowService {
	return WowService{
		Client:    client,
		Challenge: challenge,
		wg:        &sync.WaitGroup{},
	}
}

func (ws *WowService) Run(ctx context.Context) {

	for i := 1; i <= config.ClientsCount; i++ {
		select {
		case <-ctx.Done():
			log.WithField("service", config.ServiceName).Warn("Client shutting down...")
			return
		default:
			ws.wg.Add(1)
			go ws.startWork(ctx, i)
			time.Sleep(config.ConnInterval * time.Millisecond)
		}
	}
	ws.wg.Wait()
}

func (ws *WowService) startWork(ctx context.Context, id int) {
	defer ws.wg.Done()

	err := ws.Client.RunClient(ctx)
	if err != nil {
		log.WithFields(log.Fields{"service": config.ServiceName, "connection": id}).
			Errorf("Error while establishing connection to tcp-server: %v", err)
		return
	}
	defer ws.Client.CloseConn()

	taskMessage, err := ws.Client.ReceiveMessage(ctx)
	if err != nil {
		log.WithFields(log.Fields{"service": config.ServiceName, "connection": id}).
			Errorf("Error reading PoW challenge: %v", err)
		return
	}

	log.WithFields(log.Fields{"service": config.ServiceName, "connection": id}).
		Debugf("Message from server received: %s", taskMessage)

	sm, err := model.ParseServerMessage(taskMessage)
	if err != nil {
		log.WithFields(log.Fields{"service": config.ServiceName, "connection": id}).
			Errorf("Unable to unmarshal server message: %v\n", err)
		return
	}

	ws.Challenge.SetPowDifficulty(sm.Difficulty)
	nonce := ws.Challenge.GenerateSolution(ctx, sm.MessageString)
	log.WithFields(log.Fields{"service": config.ServiceName, "connection": id}).
		Infof("Found solution: %s", nonce)

	responseMessage := model.PrepareMessage(model.MessageTypeSolution, nonce, sm.Difficulty)

	if err = ws.Client.SendMessage(ctx, responseMessage.AsJsonString()); err != nil {
		log.WithFields(log.Fields{"service": config.ServiceName, "connection": id}).
			Errorf("Error while sending message: %v", err)
		return
	}

	message, _ := ws.Client.ReceiveMessage(ctx)
	log.WithFields(log.Fields{"service": config.ServiceName, "connection": id}).
		Infof("Message from server received: %s", message)
	sm, err = model.ParseServerMessage(message)
	if err != nil {
		log.WithFields(log.Fields{"service": config.ServiceName, "connection": id}).
			Errorf("Unable to unmarshal server message: %v\n", err)
		return
	}
	log.WithFields(log.Fields{"service": config.ServiceName, "connection": id}).
		Infof("Words of Wisdom: %s", sm.MessageString)
}
