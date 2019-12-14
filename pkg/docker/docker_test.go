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

func TestDockerApi_PidFromCid(t *testing.T) {
	type args struct {
		cid string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"alp", args{"c32f07b91e9cc82517e2b2eaa5808fca6c64ba359c52926532b565b5fac2f92f"}, 17317, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDockerApi()
			if err != nil {
				t.Errorf("cant create DockerAPI: %v", err)
			}
			got, err := d.PidFromCid(tt.args.cid)
			if (err != nil) != tt.wantErr {
				t.Errorf("DockerApi.PidFromCid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DockerApi.PidFromCid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDockerApi_CidFromPid(t *testing.T) {
	type args struct {
		pid int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"alp", args{17317}, "c32f07b91e9cc82517e2b2eaa5808fca6c64ba359c52926532b565b5fac2f92f", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDockerApi()
			if err != nil {
				t.Errorf("cant create DockerAPI: %v", err)
			}
			got, err := d.CidFromPid(tt.args.pid)
			if (err != nil) != tt.wantErr {
				t.Errorf("DockerApi.CidFromPid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DockerApi.CidFromPid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDockerApi_NamesFromCid(t *testing.T) {

	type args struct {
		cid string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"alp", args{"c32f07b91e9cc82517e2b2eaa5808fca6c64ba359c52926532b565b5fac2f92f"}, []string{"alp"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDockerApi()
			if err != nil {
				t.Errorf("cant create DockerAPI: %v", err)
			}
			got, err := d.NamesFromCid(tt.args.cid)
			if (err != nil) != tt.wantErr {
				t.Errorf("DockerApi.NameFromCid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DockerApi.NameFromCid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDockerApi_CidFromName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"alp", args{"alp"}, "c32f07b91e9cc82517e2b2eaa5808fca6c64ba359c52926532b565b5fac2f92f", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDockerApi()
			if err != nil {
				t.Errorf("cant create DockerAPI: %v", err)
			}
			got, err := d.CidFromName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("DockerApi.CidFromName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DockerApi.CidFromName() = %v, want %v", got, tt.want)
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
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"alp_/tmp", args{"2d8937a0f35a307a15048fadfc89702908a03977af6edceeff1081fa36086f6a", "/tmp"}, "/var/lib/docker/overlay2/94c12e447a02e00aee3e4e9c7e5fcb65b33f091d2e2b28fd119b123784cd6ee5/merged/tmp", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDockerApi()
			if err != nil {
				t.Errorf("cant create DockerAPI: %v", err)
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

func TestDockerApi_GetUserNameFromCidUid(t *testing.T) {
	type args struct {
		cid string
		uid int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"alp_root", args{"c32f07b91e9cc82517e2b2eaa5808fca6c64ba359c52926532b565b5fac2f92f", 0}, "root", false},
		{"alp_nobody", args{"c32f07b91e9cc82517e2b2eaa5808fca6c64ba359c52926532b565b5fac2f92f", 65534}, "nobody", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDockerApi()
			if err != nil {
				t.Errorf("cant create DockerAPI: %v", err)
			}
			got, err := d.GetUserNameFromCidUid(tt.args.cid, tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("DockerApi.GetUserNameFromCidUid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DockerApi.GetUserNameFromCidUid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDockerApi_GetGroupNameFromCidUid(t *testing.T) {
	type args struct {
		name string
		gid  int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"alp_root", args{"alp", 0}, "root", false},
		{"alp_ping", args{"alp", 999}, "ping", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDockerApi()
			if err != nil {
				t.Errorf("cant create DockerAPI: %v", err)
			}
			cid, err := d.CidFromName(tt.args.name)
			if err != nil {
				t.Errorf("cant get container ID: %v", tt.args.name)
			}
			got, err := d.GetGroupNameFromCidUid(cid, tt.args.gid)
			if (err != nil) != tt.wantErr {
				t.Errorf("DockerApi.GetGroupNameFromCidUid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DockerApi.GetGroupNameFromCidUid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDockerApi_GetContainerIdFromPid(t *testing.T) {

	type args struct {
		pid int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"alp", args{17317}, "c32f07b91e9cc82517e2b2eaa5808fca6c64ba359c52926532b565b5fac2f92f", false},
		{"notContainerPid", args{17317}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDockerApi()
			if err != nil {
				t.Errorf("cant create DockerAPI: %v", err)
			}
			got, err := d.GetDockerIdFromPid(tt.args.pid)
			if (err != nil) != tt.wantErr {
				t.Errorf("DockerApi.GetDockerIdFromPid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DockerApi.GetDockerIdFromPid() = %v, want %v", got, tt.want)
			}
		})
	}
}
