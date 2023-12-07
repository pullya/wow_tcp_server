package app

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/pullya/wow_tcp_server/tcp-client/internal/app/mocks"
	"github.com/pullya/wow_tcp_server/tcp-client/internal/client"
	clientMocks "github.com/pullya/wow_tcp_server/tcp-client/internal/client/mocks"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWowService_startWork(t *testing.T) {
	type fields struct {
		Client    client.IClient
		Challenge IChallenger
		wg        *sync.WaitGroup
	}
	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name   string
		fields func() fields
		args   args
		want   string
	}{
		{
			name: "Error RunClient",
			fields: func() fields {
				clientMock := &clientMocks.IClient{}
				challengeMock := &mocks.IChallenger{}
				wg := sync.WaitGroup{}
				wg.Add(1)

				clientMock.On("RunClient", mock.Anything).Return(errors.New("error"))

				return fields{
					Client:    clientMock,
					Challenge: challengeMock,
					wg:        &wg,
				}
			},
			args: args{
				ctx: context.Background(),
				id:  12,
			},
			want: fmt.Sprintf("time=\"%s\" level=error msg=\"Error while establishing connection to tcp-server: error\" connection=12 service=tcp-client\n",
				time.Now().Format("2006-01-02T15:04:05-07:00")),
		},
		{
			name: "Challenge Message wasn't received",
			fields: func() fields {
				clientMock := &clientMocks.IClient{}
				challengeMock := &mocks.IChallenger{}
				wg := sync.WaitGroup{}
				wg.Add(1)

				clientMock.On("RunClient", mock.Anything).Return(nil)
				clientMock.On("CloseConn").Return()
				clientMock.On("ReceiveMessage", mock.Anything).Return("message", errors.New("error"))

				return fields{
					Client:    clientMock,
					Challenge: challengeMock,
					wg:        &wg,
				}
			},
			args: args{
				ctx: context.Background(),
				id:  12,
			},
			want: fmt.Sprintf("time=\"%s\" level=error msg=\"Error reading PoW challenge: error\" connection=12 service=tcp-client\n",
				time.Now().Format("2006-01-02T15:04:05-07:00")),
		},
		{
			name: "Error Parse server message",
			fields: func() fields {
				clientMock := &clientMocks.IClient{}
				challengeMock := &mocks.IChallenger{}
				wg := sync.WaitGroup{}
				wg.Add(1)

				clientMock.On("RunClient", mock.Anything).Return(nil)
				clientMock.On("CloseConn").Return()
				clientMock.On("ReceiveMessage", mock.Anything).Return("message", nil)

				return fields{
					Client:    clientMock,
					Challenge: challengeMock,
					wg:        &wg,
				}
			},
			args: args{
				ctx: context.Background(),
				id:  12,
			},
			want: fmt.Sprintf("time=\"%s\" level=debug msg=\"Message from server received: message\" connection=12 service=tcp-client\ntime=\"%s\" level=error msg=\"Unable to unmarshal server message: invalid character 'm' looking for beginning of value\\n\" connection=12 service=tcp-client\n",
				time.Now().Format("2006-01-02T15:04:05-07:00"), time.Now().Format("2006-01-02T15:04:05-07:00")),
		},
		{
			name: "Error Send answer to server",
			fields: func() fields {
				clientMock := &clientMocks.IClient{}
				challengeMock := &mocks.IChallenger{}
				wg := sync.WaitGroup{}
				wg.Add(1)

				clientMock.On("RunClient", mock.Anything).Return(nil)
				clientMock.On("CloseConn").Return()
				clientMock.On("ReceiveMessage", mock.Anything).Return("{\"message_type\":\"challenge\",\"message_string\":\"Find a string that, when hashed, can be proofed 1\",\"difficulty\":10}", nil)

				challengeMock.On("SetPowDifficulty", 10).Return()
				challengeMock.On("GenerateSolution", mock.Anything, "Find a string that, when hashed, can be proofed 1").Return("123")

				clientMock.On("SendMessage", mock.Anything, mock.Anything).Return(errors.New("error"))

				return fields{
					Client:    clientMock,
					Challenge: challengeMock,
					wg:        &wg,
				}
			},
			args: args{
				ctx: context.Background(),
				id:  12,
			},
			want: fmt.Sprintf("time=\"%s\" level=debug msg=\"Message from server received: {\\\"message_type\\\":\\\"challenge\\\",\\\"message_string\\\":\\\"Find a string that, when hashed, can be proofed 1\\\",\\\"difficulty\\\":10}\" connection=12 service=tcp-client\ntime=\"%s\" level=info msg=\"Found solution: 123\" connection=12 service=tcp-client\ntime=\"%s\" level=error msg=\"Error while sending message: error\" connection=12 service=tcp-client\n",
				time.Now().Format("2006-01-02T15:04:05-07:00"), time.Now().Format("2006-01-02T15:04:05-07:00"), time.Now().Format("2006-01-02T15:04:05-07:00")),
		},
	}

	var logBuffer bytes.Buffer

	log.StandardLogger().SetLevel(log.DebugLevel)
	log.StandardLogger().SetOutput(&logBuffer)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			fields := tt.fields()
			ws := &WowService{
				Client:    fields.Client,
				Challenge: fields.Challenge,
				wg:        fields.wg,
			}

			ws.startWork(tt.args.ctx, tt.args.id)

			actual := logBuffer.String()
			assert.Equal(t, tt.want, actual)
			logBuffer.Reset()
		})
	}
}
