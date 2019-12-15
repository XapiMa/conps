package docker

import (
	"reflect"
	"testing"
)

func TestNewDockerApi(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"non error", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDockerApi()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDockerApi() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDockerApi_IsContainerProcess(t *testing.T) {
	type args struct {
		pid int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"alp", args{17335}, true},
		{"init", args{1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDockerApi()
			if err != nil {
				t.Errorf("cant create DockerAPI: %v", err)
			}
			if got := d.IsContainerProcess(tt.args.pid); got != tt.want {
				t.Errorf("DockerApi.IsContainerProcess() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDockerApi_GetPidWithCid(t *testing.T) {
	type args struct {
		cid string
	}
	tests := []struct {
		name         string
		args         args
		containerPid int
		want         int
		wantErr      bool
	}{
		{"alp", args{"c32f07b91e9cc82517e2b2eaa5808fca6c64ba359c52926532b565b5fac2f92f"}, 17335, 17317, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDockerApi()
			if err != nil {
				t.Errorf("cant create DockerAPI: %v", err)
			}
			if ok := d.IsContainerProcess(tt.containerPid); !ok {
				t.Errorf("pid %v is not container process", tt.containerPid)
			}
			got, err := d.GetPidWithCid(tt.args.cid)
			if (err != nil) != tt.wantErr {
				t.Errorf("DockerApi.PidWithCid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DockerApi.PidWithCid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDockerApi_GetCidWithPid(t *testing.T) {
	type args struct {
		pid int
	}
	tests := []struct {
		name         string
		args         args
		containerPid int
		want         string
		wantErr      bool
	}{
		{"alp", args{17335}, 17335, "c32f07b91e9cc82517e2b2eaa5808fca6c64ba359c52926532b565b5fac2f92f", false},
		{"notContainerPid", args{17317}, 17335, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDockerApi()
			if err != nil {
				t.Errorf("cant create DockerAPI: %v", err)
			}
			if ok := d.IsContainerProcess(tt.containerPid); !ok {
				t.Errorf("pid %v is not container process", tt.containerPid)
			}
			got, err := d.GetCidWithPid(tt.args.pid)
			if (err != nil) != tt.wantErr {
				t.Errorf("DockerApi.GetCidWithPid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DockerApi.GetCidWithPid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDockerApi_GetNameWithCid(t *testing.T) {

	type args struct {
		cid string
	}
	tests := []struct {
		name         string
		args         args
		containerPid int
		want         string
		wantErr      bool
	}{
		{"alp", args{"c32f07b91e9cc82517e2b2eaa5808fca6c64ba359c52926532b565b5fac2f92f"}, 17335, "alp", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDockerApi()
			if err != nil {
				t.Errorf("cant create DockerAPI: %v", err)
			}
			if ok := d.IsContainerProcess(tt.containerPid); !ok {
				t.Errorf("pid %v is not container process", tt.containerPid)
			}
			got, err := d.GetNameWithCid(tt.args.cid)
			if (err != nil) != tt.wantErr {
				t.Errorf("DockerApi.NameWithCid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DockerApi.NameWithCid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDockerApi_GetCidWithName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name         string
		args         args
		containerPid int
		want         string
		wantErr      bool
	}{
		{"alp", args{"alp"}, 17335, "c32f07b91e9cc82517e2b2eaa5808fca6c64ba359c52926532b565b5fac2f92f", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDockerApi()
			if err != nil {
				t.Errorf("cant create DockerAPI: %v", err)
			}
			if ok := d.IsContainerProcess(tt.containerPid); !ok {
				t.Errorf("pid %v is not container process", tt.containerPid)
			}
			got, err := d.GetCidWithName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("DockerApi.CidWithName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DockerApi.CidWithName() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestDockerApi_GetUserNameWithCidUid(t *testing.T) {
	type args struct {
		cid string
		uid int
	}
	tests := []struct {
		name         string
		args         args
		containerPid int
		want         string
		wantErr      bool
	}{
		{"alp_root", args{"c32f07b91e9cc82517e2b2eaa5808fca6c64ba359c52926532b565b5fac2f92f", 0}, 17335, "root", false},
		{"alp_nobody", args{"c32f07b91e9cc82517e2b2eaa5808fca6c64ba359c52926532b565b5fac2f92f", 65534}, 17335, "nobody", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDockerApi()
			if err != nil {
				t.Errorf("cant create DockerAPI: %v", err)
			}
			if ok := d.IsContainerProcess(tt.containerPid); !ok {
				t.Errorf("pid %v is not container process", tt.containerPid)
			}
			got, err := d.GetUserNameWithCidUid(tt.args.cid, tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("DockerApi.GetUserNameWithCidUid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DockerApi.GetUserNameWithCidUid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDockerApi_GetGroupNameWithCidGid(t *testing.T) {
	type args struct {
		cid string
		gid int
	}
	tests := []struct {
		name         string
		args         args
		containerPid int
		want         string
		wantErr      bool
	}{
		{"alp_root", args{"c32f07b91e9cc82517e2b2eaa5808fca6c64ba359c52926532b565b5fac2f92f", 0}, 17335, "root", false},
		{"alp_ping", args{"c32f07b91e9cc82517e2b2eaa5808fca6c64ba359c52926532b565b5fac2f92f", 999}, 17335, "ping", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDockerApi()
			if err != nil {
				t.Errorf("cant create DockerAPI: %v", err)
			}
			if ok := d.IsContainerProcess(tt.containerPid); !ok {
				t.Errorf("pid %v is not container process", tt.containerPid)
			}
			got, err := d.GetGroupNameWithCidGid(tt.args.cid, tt.args.gid)
			if (err != nil) != tt.wantErr {
				t.Errorf("DockerApi.GetGroupNameWithCidUid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DockerApi.GetGroupNameWithCidUid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDockerApi_ContainerPathToHostPath(t *testing.T) {
	type args struct {
		cid  string
		path string
	}
	tests := []struct {
		name         string
		args         args
		containerPid int
		want         string
		wantErr      bool
	}{
		{"alp/tmp", args{"c32f07b91e9cc82517e2b2eaa5808fca6c64ba359c52926532b565b5fac2f92f", "/tmp"}, 17335, "/var/lib/docker/overlay2/939adafd908a46ede6b9e142d91686c0272f51c492931e970d8c7127b56d9c9c/merged/tmp", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDockerApi()
			if err != nil {
				t.Errorf("cant create DockerAPI: %v", err)
			}
			if ok := d.IsContainerProcess(tt.containerPid); !ok {
				t.Errorf("pid %v is not container process", tt.containerPid)
			}
			got, err := d.ContainerPathToHostPath(tt.args.cid, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("DockerApi.ContainerPathToHostPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DockerApi.ContainerPathToHostPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
