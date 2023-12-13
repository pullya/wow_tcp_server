package config

import (
	"testing"
)

func Test_validatePort(t *testing.T) {
	t.Parallel()
	type args struct {
		in string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "Success #1",
			args:    args{in: "80"},
			want:    80,
			wantErr: false,
		},
		{
			name:    "Failed #1 spaces",
			args:    args{in: " 80,"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Failed #2 negative",
			args:    args{in: "-8081"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Failed #3 too big",
			args:    args{in: "12345567"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Failed #4 char",
			args:    args{in: "eighty"},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validatePort(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePort() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("validatePort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateTimeout(t *testing.T) {
	t.Parallel()
	type args struct {
		in string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "Success #1",
			args:    args{in: "100"},
			want:    100,
			wantErr: false,
		},
		{
			name:    "Success #2 zero",
			args:    args{in: "0"},
			want:    0,
			wantErr: false,
		},
		{
			name:    "Success #3 multi zero",
			args:    args{in: "0000"},
			want:    0,
			wantErr: false,
		},
		{
			name:    "Failed #1 spaces",
			args:    args{in: "10, "},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Failed #2 negative",
			args:    args{in: "-80"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Failed #3 char",
			args:    args{in: "ten"},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateTimeout(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateTimeout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("validateTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateDifficulty(t *testing.T) {
	t.Parallel()
	type args struct {
		in string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "Success #1",
			args:    args{in: "20"},
			want:    20,
			wantErr: false,
		},
		{
			name:    "Success #2 zero",
			args:    args{in: "0"},
			want:    0,
			wantErr: false,
		},
		{
			name:    "Success #3 multi zero",
			args:    args{in: "0000"},
			want:    0,
			wantErr: false,
		},
		{
			name:    "Failed #1 spaces",
			args:    args{in: " 80,"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Failed #2 negative",
			args:    args{in: "-112"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Failed #3 too big",
			args:    args{in: "257"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Failed #4 char",
			args:    args{in: "eighty"},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateDifficulty(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateDifficulty() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("validateDifficulty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateLogLevel(t *testing.T) {
	t.Parallel()
	type args struct {
		in string
	}
	tests := []struct {
		name    string
		args    args
		want    LogLevel
		wantErr bool
	}{
		{
			name:    "Success #1",
			args:    args{in: "Info"},
			want:    LogLevel("Info"),
			wantErr: false,
		},
		{
			name:    "Failed #1 spaces",
			args:    args{in: "Debug, "},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Failed #2 unknown",
			args:    args{in: "Highest"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Failed #3 capital",
			args:    args{in: "ERROR"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateLogLevel(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateLogLevel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("validateLogLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}
