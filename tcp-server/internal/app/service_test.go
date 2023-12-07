package app

import (
	"context"
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

func Test_generatePoWChallenge(t *testing.T) {
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
			want: config.ProofString + " 1",
		},
		{
			name: "Success #2",
			args: args{cnt: -1234},
			want: config.ProofString + " -1234",
		},
		{
			name: "Success #3",
			args: args{cnt: 0},
			want: config.ProofString + " 0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generatePoWChallenge(tt.args.cnt); got != tt.want {
				t.Errorf("generatePoWChallenge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWowService_doProofOfWork(t *testing.T) {
	type fields struct {
		Server    server.IServer
		Storage   storage.IStorage
		Challenge IChallenger
	}
	type args struct {
		ctx context.Context
		id  int
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
				serverMock := &serverMocks.IServer{}
				storageMock := &storageMocks.IStorage{}
				challengeMock := &mocks.IChallenger{}

				challengeMock.On("GetPowDifficulty").Return(10)
				serverMock.On("SendMessage", mock.Anything, mock.Anything).Return(errors.New("error"))

				return fields{
					Server:    serverMock,
					Storage:   storageMock,
					Challenge: challengeMock,
				}
			},
			args: args{
				context.Background(),
				21,
			},
			wantErr: true,
		},
		{
			name: "Receive client message Error",
			fields: func() fields {
				serverMock := &serverMocks.IServer{}
				storageMock := &storageMocks.IStorage{}
				challengeMock := &mocks.IChallenger{}

				challengeMock.On("GetPowDifficulty").Return(10)
				serverMock.On("SendMessage", mock.Anything, mock.Anything).Return(nil)
				serverMock.On("ReceiveMessage", mock.Anything).Return("message", errors.New("error"))

				return fields{
					Server:    serverMock,
					Storage:   storageMock,
					Challenge: challengeMock,
				}
			},
			args: args{
				context.Background(),
				21,
			},
			wantErr: true,
		},
		{
			name: "Parse message Error",
			fields: func() fields {
				serverMock := &serverMocks.IServer{}
				storageMock := &storageMocks.IStorage{}
				challengeMock := &mocks.IChallenger{}

				challengeMock.On("GetPowDifficulty").Return(10)
				serverMock.On("SendMessage", mock.Anything, mock.Anything).Return(nil)
				serverMock.On("ReceiveMessage", mock.Anything).Return("message", nil)

				return fields{
					Server:    serverMock,
					Storage:   storageMock,
					Challenge: challengeMock,
				}
			},
			args: args{
				context.Background(),
				21,
			},
			wantErr: true,
		},
		{
			name: "Get uint64 Error",
			fields: func() fields {
				serverMock := &serverMocks.IServer{}
				storageMock := &storageMocks.IStorage{}
				challengeMock := &mocks.IChallenger{}

				challengeMock.On("GetPowDifficulty").Return(10)
				serverMock.On("SendMessage", mock.Anything, mock.Anything).Return(nil)
				serverMock.On("ReceiveMessage", mock.Anything).Return("{\"message_type\":\"solution\",\"message_string\":\"-2450\",\"difficulty\":10}", nil)

				return fields{
					Server:    serverMock,
					Storage:   storageMock,
					Challenge: challengeMock,
				}
			},
			args: args{
				context.Background(),
				21,
			},
			wantErr: true,
		},
		{
			name: "Error",
			fields: func() fields {
				serverMock := &serverMocks.IServer{}
				storageMock := &storageMocks.IStorage{}
				challengeMock := &mocks.IChallenger{}

				challengeMock.On("GetPowDifficulty").Return(10)
				serverMock.On("SendMessage", mock.Anything, mock.Anything).Return(nil)
				serverMock.On("ReceiveMessage", mock.Anything).Return("{\"message_type\":\"solution\",\"message_string\":\"2450\",\"difficulty\":10}", nil)
				challengeMock.On("IsValidPoW", mock.Anything, uint64(2450)).Return(false)

				return fields{
					Server:    serverMock,
					Storage:   storageMock,
					Challenge: challengeMock,
				}
			},
			args: args{
				context.Background(),
				21,
			},
			wantErr: true,
		},
		{
			name: "Error send challenge message",
			fields: func() fields {
				serverMock := &serverMocks.IServer{}
				storageMock := &storageMocks.IStorage{}
				challengeMock := &mocks.IChallenger{}

				challengeMock.On("GetPowDifficulty").Return(10)
				serverMock.On("SendMessage", mock.Anything, mock.Anything).Return(nil)
				serverMock.On("ReceiveMessage", mock.Anything).Return("{\"message_type\":\"solution\",\"message_string\":\"2450\",\"difficulty\":10}", nil)
				challengeMock.On("IsValidPoW", mock.Anything, uint64(2450)).Return(true)

				return fields{
					Server:    serverMock,
					Storage:   storageMock,
					Challenge: challengeMock,
				}
			},
			args: args{
				context.Background(),
				21,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ws := &WowService{
				Server:    tt.fields().Server,
				Storage:   tt.fields().Storage,
				Challenge: tt.fields().Challenge,
			}
			if err := ws.doProofOfWork(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("WowService.doProofOfWork() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
