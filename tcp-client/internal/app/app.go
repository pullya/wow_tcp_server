package app

import (
	"context"
	"sync"
	"time"

	"github.com/pullya/wow_tcp_server/tcp-client/internal/client"
	"github.com/pullya/wow_tcp_server/tcp-client/internal/config"
	"github.com/pullya/wow_tcp_server/tcp-client/internal/model"
)

type App struct {
	client    client.ClientProvider
	challenge Challenger
	wg        *sync.WaitGroup
}

func New(client client.ClientProvider, challenge Challenger) App {
	return App{
		client:    client,
		challenge: challenge,
		wg:        &sync.WaitGroup{},
	}
}

func (a *App) Run(ctx context.Context) {

	for i := 1; i <= config.Config.ClientsCount; i++ {
		select {
		case <-ctx.Done():
			config.Logger.Warn("Client shutting down...")
			return
		default:
			a.wg.Add(1)
			go a.startWork(ctx, i)
			time.Sleep(config.Config.ConnInterval * time.Millisecond)
		}
	}
	a.wg.Wait()
}

func (a *App) startWork(ctx context.Context, id int) {
	defer a.wg.Done()

	conn, err := a.client.Run(ctx)
	if err != nil {
		config.Logger.WithField("connection", id).Errorf("Error while establishing connection to tcp-server: %v", err)
		return
	}
	defer a.client.CloseConn(conn)

	requestMessage := model.PrepareMessage("", model.MessageTypeRequest, "", 0)

	if err = a.client.SendMessage(ctx, conn, requestMessage.AsJsonString()); err != nil {
		config.Logger.WithField("connection", id).Errorf("Error while sending request message: %v", err)
		return
	}

	taskMessage, err := a.client.ReceiveMessage(ctx, conn)
	if err != nil {
		config.Logger.WithField("connection", id).Errorf("Error reading PoW challenge: %v", err)
		return
	}

	config.Logger.WithField("connection", id).Debugf("Message from server received: %s", taskMessage)

	sm, err := model.ParseServerMessage(taskMessage)
	if err != nil {
		config.Logger.WithField("connection", id).Errorf("Unable to unmarshal server message: %v\n", err)
		return
	}
	a.client.CloseConn(conn)

	a.challenge.SetDifficulty(sm.Difficulty)
	nonce := a.challenge.GenerateSolution(ctx, sm.MessageString)
	config.Logger.WithField("connection", id).Infof("Found solution: %s", nonce)

	responseMessage := model.PrepareMessage(sm.RequestID, model.MessageTypeSolution, nonce, sm.Difficulty)

	conn, err = a.client.Run(ctx)
	if err != nil {
		config.Logger.WithField("connection", id).Errorf("Error while establishing connection to tcp-server: %v", err)
		return
	}

	if err = a.client.SendMessage(ctx, conn, responseMessage.AsJsonString()); err != nil {
		config.Logger.WithField("connection", id).Errorf("Error while sending message: %v", err)
		return
	}

	message, err := a.client.ReceiveMessage(ctx, conn)
	if err != nil {
		return
	}
	config.Logger.WithField("connection", id).Infof("Message from server received: %s", message)
	sm, err = model.ParseServerMessage(message)
	if err != nil {
		config.Logger.WithField("connection", id).Errorf("Unable to unmarshal server message: %v\n", err)
		return
	}
	config.Logger.WithField("connection", id).Infof("Words of Wisdom: %s", sm.MessageString)
}
