package ps_test

import (
	"reflect"
	"testing"

	"github.com/xapima/conps/pkg/ps"
)

func TestCmdline(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"1", args{"./test/proc/1"}, "/bin/sh", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ps.Cmdline(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Cmdline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Cmdline() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCwd(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"..", args{"./test/proc/1"}, "..", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ps.Cwd(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Cwd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Cwd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{"ENV=envPATH=/bin", args{"./test/proc/1"}, map[string]string{"ENV": "env", "PATH": "/bin"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ps.Env(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Env() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Env() = %v, want %v", got, tt.want)
			}
		})
	}
}
