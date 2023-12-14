package app

import (
	"context"
	"net"
	"testing"

	"github.com/pkg/errors"
	"github.com/pullya/wow_tcp_server/tcp-server/internal/app/mocks"
	"github.com/pullya/wow_tcp_server/tcp-server/internal/config"
	"github.com/pullya/wow_tcp_server/tcp-server/internal/model"
	"github.com/pullya/wow_tcp_server/tcp-server/internal/server"
	serverMocks "github.com/pullya/wow_tcp_server/tcp-server/internal/server/mocks"
	"github.com/pullya/wow_tcp_server/tcp-server/internal/storage"
	storageMocks "github.com/pullya/wow_tcp_server/tcp-server/internal/storage/mocks"
	"github.com/stretchr/testify/mock"
)

func Test_generatePOWChallenge(t *testing.T) {
	t.Parallel()
	type args struct {
		cnt string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Success #1",
			args: args{cnt: "1q2w3e"},
			want: config.Config.ProofString + " 1q2w3e",
		},
		{
			name: "Success #2",
			args: args{cnt: "-1234"},
			want: config.Config.ProofString + " -1234",
		},
		{
			name: "Success #3",
			args: args{cnt: "0"},
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

func TestApp_sendChallenge(t *testing.T) {
	t.Parallel()

	tConn := *new(net.Conn)
	config.InitLogger()
	ctx := context.Background()

	type fields struct {
		server       server.ServerProvider
		storage      storage.Storageer
		requeststore storage.Requester
		challenge    Challenger
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
			name: "Error send challenge message",
			fields: func() fields {
				serverMock := &serverMocks.ServerProvider{}
				storageMock := &storageMocks.Storageer{}
				requeststoreMock := &storageMocks.Requester{}
				challengeMock := &mocks.Challenger{}

				challengeMock.On("Difficulty").Return(10)
				serverMock.On("SendMessage", mock.Anything, tConn, mock.Anything).Return(errors.New("error"))

				return fields{
					server:       serverMock,
					storage:      storageMock,
					requeststore: requeststoreMock,
					challenge:    challengeMock,
				}
			},
			args: args{
				ctx,
				tConn,
				21,
			},
			wantErr: true,
		},
		{
			name: "Challenge message successfully sent",
			fields: func() fields {
				serverMock := &serverMocks.ServerProvider{}
				storageMock := &storageMocks.Storageer{}
				requeststoreMock := &storageMocks.Requester{}
				challengeMock := &mocks.Challenger{}

				challengeMock.On("Difficulty").Return(10)
				serverMock.On("SendMessage", mock.Anything, tConn, mock.Anything).Return(nil)
				requeststoreMock.On("Add", mock.Anything, mock.Anything).Return()

				return fields{
					server:       serverMock,
					storage:      storageMock,
					requeststore: requeststoreMock,
					challenge:    challengeMock,
				}
			},
			args: args{
				ctx,
				tConn,
				21,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &App{
				server:       tt.fields().server,
				storage:      tt.fields().storage,
				requeststore: tt.fields().requeststore,
				challenge:    tt.fields().challenge,
			}
			if err := a.sendChallenge(tt.args.ctx, tt.args.conn, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("App.sendChallenge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApp_validatePOW(t *testing.T) {
	t.Parallel()

	config.InitLogger()
	ctx := context.Background()
	uid := storage.GenUID()

	type fields struct {
		server       server.ServerProvider
		storage      storage.Storageer
		requeststore storage.Requester
		challenge    Challenger
	}
	type args struct {
		ctx            context.Context
		clientResponse model.Message
		id             int
	}
	tests := []struct {
		name    string
		fields  func() fields
		args    args
		wantErr bool
	}{
		{
			name: "Error find request in store",
			fields: func() fields {
				serverMock := &serverMocks.ServerProvider{}
				storageMock := &storageMocks.Storageer{}
				requeststoreMock := &storageMocks.Requester{}
				challengeMock := &mocks.Challenger{}

				requeststoreMock.On("Get", mock.Anything, uid).Return(false, errors.New("error"))

				return fields{
					server:       serverMock,
					storage:      storageMock,
					requeststore: requeststoreMock,
					challenge:    challengeMock,
				}
			},
			args: args{
				ctx,
				model.Message{RequestID: uid, MessageType: model.MessageTypeRequest, MessageString: "answer", Difficulty: 21},
				21,
			},
			wantErr: true,
		},
		{
			name: "Error request already handled",
			fields: func() fields {
				serverMock := &serverMocks.ServerProvider{}
				storageMock := &storageMocks.Storageer{}
				requeststoreMock := &storageMocks.Requester{}
				challengeMock := &mocks.Challenger{}

				requeststoreMock.On("Get", mock.Anything, uid).Return(true, nil)

				return fields{
					server:       serverMock,
					storage:      storageMock,
					requeststore: requeststoreMock,
					challenge:    challengeMock,
				}
			},
			args: args{
				ctx,
				model.Message{RequestID: uid, MessageType: model.MessageTypeRequest, MessageString: "answer", Difficulty: 21},
				21,
			},
			wantErr: true,
		},
		{
			name: "Error parse client response",
			fields: func() fields {
				serverMock := &serverMocks.ServerProvider{}
				storageMock := &storageMocks.Storageer{}
				requeststoreMock := &storageMocks.Requester{}
				challengeMock := &mocks.Challenger{}

				requeststoreMock.On("Get", mock.Anything, uid).Return(false, nil)

				return fields{
					server:       serverMock,
					storage:      storageMock,
					requeststore: requeststoreMock,
					challenge:    challengeMock,
				}
			},
			args: args{
				ctx,
				model.Message{RequestID: uid, MessageType: model.MessageTypeRequest, MessageString: "answer", Difficulty: 21},
				21,
			},
			wantErr: true,
		},
		{
			name: "Error validate POW",
			fields: func() fields {
				serverMock := &serverMocks.ServerProvider{}
				storageMock := &storageMocks.Storageer{}
				requeststoreMock := &storageMocks.Requester{}
				challengeMock := &mocks.Challenger{}

				requeststoreMock.On("Get", mock.Anything, uid).Return(false, nil)
				challengeMock.On("IsValid", mock.Anything, uint64(2450)).Return(false)

				return fields{
					server:       serverMock,
					storage:      storageMock,
					requeststore: requeststoreMock,
					challenge:    challengeMock,
				}
			},
			args: args{
				ctx,
				model.Message{RequestID: uid, MessageType: model.MessageTypeSolution, MessageString: "2450", Difficulty: 21},
				21,
			},
			wantErr: true,
		},
		{
			name: "Success",
			fields: func() fields {
				serverMock := &serverMocks.ServerProvider{}
				storageMock := &storageMocks.Storageer{}
				requeststoreMock := &storageMocks.Requester{}
				challengeMock := &mocks.Challenger{}

				requeststoreMock.On("Get", mock.Anything, uid).Return(false, nil)
				challengeMock.On("IsValid", mock.Anything, uint64(2450)).Return(true)

				return fields{
					server:       serverMock,
					storage:      storageMock,
					requeststore: requeststoreMock,
					challenge:    challengeMock,
				}
			},
			args: args{
				ctx,
				model.Message{RequestID: uid, MessageType: model.MessageTypeSolution, MessageString: "2450", Difficulty: 21},
				21,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &App{
				server:       tt.fields().server,
				storage:      tt.fields().storage,
				requeststore: tt.fields().requeststore,
				challenge:    tt.fields().challenge,
			}
			if err := a.validatePOW(tt.args.ctx, tt.args.clientResponse, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("App.validatePOW() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApp_sendWOW(t *testing.T) {
	t.Parallel()

	tConn := *new(net.Conn)
	config.InitLogger()
	ctx := context.Background()
	uid := storage.GenUID()

	type fields struct {
		server       server.ServerProvider
		storage      storage.Storageer
		requeststore storage.Requester
		challenge    Challenger
	}
	type args struct {
		ctx  context.Context
		conn net.Conn
		uid  string
		id   int
	}
	tests := []struct {
		name    string
		fields  func() fields
		args    args
		wantErr bool
	}{
		{
			name: "Error sending message",
			fields: func() fields {
				serverMock := &serverMocks.ServerProvider{}
				storageMock := &storageMocks.Storageer{}
				requeststoreMock := &storageMocks.Requester{}
				challengeMock := &mocks.Challenger{}

				storageMock.On("GetRandomWOW", mock.Anything).Return("Random Word of Wisdom")
				serverMock.On("SendMessage", mock.Anything, tConn, mock.Anything).Return(errors.New("error"))

				return fields{
					server:       serverMock,
					storage:      storageMock,
					requeststore: requeststoreMock,
					challenge:    challengeMock,
				}
			},
			args: args{
				ctx,
				tConn,
				uid,
				21,
			},
			wantErr: true,
		},
		{
			name: "Error sending message",
			fields: func() fields {
				serverMock := &serverMocks.ServerProvider{}
				storageMock := &storageMocks.Storageer{}
				requeststoreMock := &storageMocks.Requester{}
				challengeMock := &mocks.Challenger{}

				storageMock.On("GetRandomWOW", mock.Anything).Return("Random Word of Wisdom")
				serverMock.On("SendMessage", mock.Anything, tConn, mock.Anything).Return(nil)
				requeststoreMock.On("Set", mock.Anything, uid).Return(errors.New("error"))

				return fields{
					server:       serverMock,
					storage:      storageMock,
					requeststore: requeststoreMock,
					challenge:    challengeMock,
				}
			},
			args: args{
				ctx,
				tConn,
				uid,
				21,
			},
			wantErr: true,
		},
		{
			name: "Success",
			fields: func() fields {
				serverMock := &serverMocks.ServerProvider{}
				storageMock := &storageMocks.Storageer{}
				requeststoreMock := &storageMocks.Requester{}
				challengeMock := &mocks.Challenger{}

				storageMock.On("GetRandomWOW", mock.Anything).Return("Random Word of Wisdom")
				serverMock.On("SendMessage", mock.Anything, tConn, mock.Anything).Return(nil)
				requeststoreMock.On("Set", mock.Anything, uid).Return(nil)

				return fields{
					server:       serverMock,
					storage:      storageMock,
					requeststore: requeststoreMock,
					challenge:    challengeMock,
				}
			},
			args: args{
				ctx,
				tConn,
				uid,
				21,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &App{
				server:       tt.fields().server,
				storage:      tt.fields().storage,
				requeststore: tt.fields().requeststore,
				challenge:    tt.fields().challenge,
			}
			if err := a.sendWOW(tt.args.ctx, tt.args.conn, tt.args.uid, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("App.sendWOW() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
