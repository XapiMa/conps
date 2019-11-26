package docker

import (
	"reflect"
	"testing"
)

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
		{"alp", args{"2e9e95413319081bc36b474be02917096e0e6858b9c66fa9204596485b34d9b6"}, []string{"alp"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDockerApi()
			if err != nil {
				t.Errorf("cant new DockerApid: %v", err)
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
				t.Errorf("cant New Docker API: %v", err)
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
