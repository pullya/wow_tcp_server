package app

import (
	"testing"

	"github.com/pullya/wow_tcp_server/tcp-server/internal/config"
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
