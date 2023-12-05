package model

import (
	"reflect"
	"testing"
)

func TestMessage_AsJsonString(t *testing.T) {
	t.Parallel()
	type fields struct {
		MessageType   string
		MessageString string
		Difficulty    int
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "Success #1",
			fields: fields{
				MessageType:   "solution",
				MessageString: "Challenge string 1000",
				Difficulty:    3,
			},
			want: []byte("{\"message_type\":\"solution\",\"message_string\":\"Challenge string 1000\",\"difficulty\":3}\n"),
		},
		{
			name: "Success #2",
			fields: fields{
				MessageType:   "",
				MessageString: "",
				Difficulty:    0,
			},
			want: []byte("{\"message_type\":\"\",\"message_string\":\"\",\"difficulty\":0}\n"),
		},
		{
			name:   "Success #3",
			fields: fields{},
			want:   []byte("{\"message_type\":\"\",\"message_string\":\"\",\"difficulty\":0}\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Message{
				MessageType:   tt.fields.MessageType,
				MessageString: tt.fields.MessageString,
				Difficulty:    tt.fields.Difficulty,
			}
			if got := m.AsJsonString(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Message.AsJsonString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseServerMessage(t *testing.T) {
	t.Parallel()
	type args struct {
		message string
	}
	tests := []struct {
		name    string
		args    args
		want    Message
		wantErr bool
	}{
		{
			name:    "Successfully parsed #1",
			args:    args{message: "{\"message_type\":\"challenge\",\"message_string\":\"Find a string that 2\",\"difficulty\":1}"},
			want:    Message{MessageType: "challenge", MessageString: "Find a string that 2", Difficulty: 1},
			wantErr: false,
		},
		{
			name:    "Successfully parsed #2",
			args:    args{message: "{\"message_type\":\"solution\",\"message_string\":\"2\",\"difficulty\":5}"},
			want:    Message{MessageType: "solution", MessageString: "2", Difficulty: 5},
			wantErr: false,
		},
		{
			name:    "Failed to parse #1",
			args:    args{message: "Some message"},
			want:    Message{},
			wantErr: true,
		},
		{
			name:    "Failed to parse #2",
			args:    args{message: "{\"id\":\"73635\",\"result_string\":2,\"difficulty\":5}"},
			want:    Message{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseServerMessage(tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseServerMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseServerMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateMessageType(t *testing.T) {
	t.Parallel()
	type args struct {
		messageType string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Success #1",
			args: args{messageType: MessageTypeChallenge},
			want: true,
		},
		{
			name: "Success #2",
			args: args{messageType: MessageTypeWow},
			want: true,
		},
		{
			name: "Failed (empty) #1",
			args: args{messageType: ""},
			want: false,
		},
		{
			name: "Failed (wrong) #2",
			args: args{messageType: "something"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateMessageType(tt.args.messageType); got != tt.want {
				t.Errorf("validateMessageType() = %v, want %v", got, tt.want)
			}
		})
	}
}
