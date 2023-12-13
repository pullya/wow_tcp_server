package app

import (
	"context"
	"net"
	"testing"

	"github.com/pkg/errors"
	"github.com/pullya/wow_tcp_server/tcp-server/internal/app/mocks"
	"github.com/pullya/wow_tcp_server/tcp-server/internal/config"
	"github.com/pullya/wow_tcp_server/tcp-server/internal/server"
	serverMocks "github.com/pullya/wow_tcp_server/tcp-server/internal/server/mocks"
	"github.com/pullya/wow_tcp_server/tcp-server/internal/storage"
	storageMocks "github.com/pullya/wow_tcp_server/tcp-server/internal/storage/mocks"
	"github.com/stretchr/testify/mock"
)

func Test_generatePOWChallenge(t *testing.T) {
	t.Parallel()
	type args struct {
		cnt int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Success #1",
			args: args{cnt: 1},
			want: config.Config.ProofString + " 1",
		},
		{
			name: "Success #2",
			args: args{cnt: -1234},
			want: config.Config.ProofString + " -1234",
		},
		{
			name: "Success #3",
			args: args{cnt: 0},
			want: config.Config.ProofString + " 0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generatePOWChallenge(tt.args.cnt); got != tt.want {
				t.Errorf("generatePOWChallenge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWowService_doProofOfWork(t *testing.T) {
	tConn := *new(net.Conn)
	config.InitLogger()

	type fields struct {
		Server    server.ServerProvider
		Storage   storage.Storageer
		Challenge Challenger
	}
	type args struct {
		ctx  context.Context
		conn net.Conn
		id   int
	}
	tests := []struct {
		name    string
		fields  func() fields
		args    args
		wantErr bool
	}{
		{
			name: "Send challenge message Error",
			fields: func() fields {
				serverMock := &serverMocks.ServerProvider{}
				storageMock := &storageMocks.Storageer{}
				challengeMock := &mocks.Challenger{}

				challengeMock.On("Difficulty").Return(10)
				serverMock.On("SendMessage", mock.Anything, tConn, mock.Anything).Return(errors.New("error"))

				return fields{
					Server:    serverMock,
					Storage:   storageMock,
					Challenge: challengeMock,
				}
			},
			args: args{
				context.Background(),
				tConn,
				21,
			},
			wantErr: true,
		},
		{
			name: "Receive client message Error",
			fields: func() fields {
				serverMock := &serverMocks.ServerProvider{}
				storageMock := &storageMocks.Storageer{}
				challengeMock := &mocks.Challenger{}

				challengeMock.On("Difficulty").Return(10)
				serverMock.On("SendMessage", mock.Anything, tConn, mock.Anything).Return(nil)
				serverMock.On("ReceiveMessage", mock.Anything, tConn).Return("message", errors.New("error"))

				return fields{
					Server:    serverMock,
					Storage:   storageMock,
					Challenge: challengeMock,
				}
			},
			args: args{
				context.Background(),
				tConn,
				21,
			},
			wantErr: true,
		},
		{
			name: "Parse message Error",
			fields: func() fields {
				serverMock := &serverMocks.ServerProvider{}
				storageMock := &storageMocks.Storageer{}
				challengeMock := &mocks.Challenger{}

				challengeMock.On("Difficulty").Return(10)
				serverMock.On("SendMessage", mock.Anything, tConn, mock.Anything).Return(nil)
				serverMock.On("ReceiveMessage", mock.Anything, tConn).Return("message", nil)

				return fields{
					Server:    serverMock,
					Storage:   storageMock,
					Challenge: challengeMock,
				}
			},
			args: args{
				context.Background(),
				tConn,
				21,
			},
			wantErr: true,
		},
		{
			name: "Get uint64 Error",
			fields: func() fields {
				serverMock := &serverMocks.ServerProvider{}
				storageMock := &storageMocks.Storageer{}
				challengeMock := &mocks.Challenger{}

				challengeMock.On("Difficulty").Return(10)
				serverMock.On("SendMessage", mock.Anything, tConn, mock.Anything).Return(nil)
				serverMock.On("ReceiveMessage", mock.Anything, tConn).Return("{\"message_type\":\"solution\",\"message_string\":\"-2450\",\"difficulty\":10}", nil)

				return fields{
					Server:    serverMock,
					Storage:   storageMock,
					Challenge: challengeMock,
				}
			},
			args: args{
				context.Background(),
				tConn,
				21,
			},
			wantErr: true,
		},
		{
			name: "Error",
			fields: func() fields {
				serverMock := &serverMocks.ServerProvider{}
				storageMock := &storageMocks.Storageer{}
				challengeMock := &mocks.Challenger{}

				challengeMock.On("Difficulty").Return(10)
				serverMock.On("SendMessage", mock.Anything, tConn, mock.Anything).Return(nil)
				serverMock.On("ReceiveMessage", mock.Anything, tConn).Return("{\"message_type\":\"solution\",\"message_string\":\"2450\",\"difficulty\":10}", nil)
				challengeMock.On("IsValid", mock.Anything, uint64(2450)).Return(false)

				return fields{
					Server:    serverMock,
					Storage:   storageMock,
					Challenge: challengeMock,
				}
			},
			args: args{
				context.Background(),
				tConn,
				21,
			},
			wantErr: true,
		},
		{
			name: "Error send challenge message",
			fields: func() fields {
				serverMock := &serverMocks.ServerProvider{}
				storageMock := &storageMocks.Storageer{}
				challengeMock := &mocks.Challenger{}

				challengeMock.On("Difficulty").Return(10)
				serverMock.On("SendMessage", mock.Anything, tConn, mock.Anything).Return(nil)
				serverMock.On("ReceiveMessage", mock.Anything, tConn).Return("{\"message_type\":\"solution\",\"message_string\":\"2450\",\"difficulty\":10}", nil)
				challengeMock.On("IsValid", mock.Anything, uint64(2450)).Return(true)

				return fields{
					Server:    serverMock,
					Storage:   storageMock,
					Challenge: challengeMock,
				}
			},
			args: args{
				context.Background(),
				tConn,
				21,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ws := &App{
				server:    tt.fields().Server,
				storage:   tt.fields().Storage,
				challenge: tt.fields().Challenge,
			}
			if err := ws.doProofOfWork(tt.args.ctx, tt.args.conn, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("App.doProofOfWork() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
