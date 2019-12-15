package docker

import (
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

// func TestDockerApi_containerUpDown(t *testing.T) {
// 	type args struct {
// 		repoTag       string
// 		containerName string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{"alpine:3.7", args{"alpine:3.7", "teat_alpine_3.7"}, false},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			d, err := NewDockerApi()
// 			if err != nil {
// 				t.Errorf("cant create DockerAPI: %v", err)
// 			}
// 			containerID, err := d.containerUp(tt.args.repoTag, tt.args.containerName)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("DockerApi.containerUp() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			t.Error("container Uped!")
// 			// ContainerDownがうまく動いていない気がする
// 			if err := d.containerDown(containerID); (err != nil) != tt.wantErr {
// 				t.Errorf("DockerApi.containerDown() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
