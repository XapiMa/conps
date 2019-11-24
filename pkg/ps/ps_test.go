// you must do test in go golang:1.13.4-alpine3.10/go/src/github.com/xapima/conps/pkg/ps
package ps_test

import (
	"reflect"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/xapima/conps/pkg/ps"
)

func init() {
	log.SetLevel(log.DebugLevel)
}
func TestCmdline(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"/proc/1", args{"/proc/1"}, []string{"/bin/sh"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ps.Cmdline(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Cmdline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
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
		{"cwd", args{"/proc/1"}, "/go/src/github.com/xapima/conps/pkg/ps", false},
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
		{"golang:1.13.4-alpine3.10", args{"/proc/1"}, map[string]string{
			"PATH":     "/go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
			"HOSTNAME": "36190f47e449", "TERM": "xterm", "GOLANG_VERSION": "1.13.4", "GOPATH": "/go", "HOME": "/root"}, false,
		},
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

func TestExe(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"/proc/1", args{"/proc/1"}, "/bin/busybox", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ps.Exe(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Exe() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFd(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{"proc/1", args{"/proc/1"}, map[string]string{"0": "/dev/pts/0", "1": "/dev/pts/0", "10": "/dev/tty", "2": "/dev/pts/0"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ps.Fd(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatus(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    ps.StatusVals
		wantErr bool
	}{
		{"/proc/1", args{"/proc/1"}, ps.StatusVals{Name: "sh", Umask: 0o0022, State: "S",
			Tgid: 1, Pid: 1, PPid: 0, TracerPid: 0, Uids: []int32{0, 0, 0, 0}, Gids: []int32{0, 0, 0, 0}, Threads: 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ps.Status(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Status() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Status() = %v, want %v", got, tt.want)
			}
		})
	}
}
