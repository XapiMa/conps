package docker

import (
	"testing"
)

func TestDockerApi_NameFromCid(t *testing.T) {

	type args struct {
		cid string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"alp", args{"2e9e95413319081bc36b474be02917096e0e6858b9c66fa9204596485b34d9b6"}, "alp", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDockerApi()
			if err != nil {
				t.Errorf("cant new DockerApid: %v", err)
			}
			got, err := d.NameFromCid(tt.args.cid)
			if (err != nil) != tt.wantErr {
				t.Errorf("DockerApi.NameFromCid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DockerApi.NameFromCid() = %v, want %v", got, tt.want)
			}
		})
	}
}
