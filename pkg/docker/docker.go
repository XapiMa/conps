package docker

import (
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"github.com/xapima/conps/pkg/util"
	"golang.org/x/net/context"
)

type CidInspect map[string]types.ContainerJSON
type CidSet map[string]struct{}
type CidName map[string]string
type PidCid map[int]string

type DockerApi struct {
	cli        *client.Client
	cidinspect CidInspect
	cidset     CidSet
	cidname    CidName
	Pidcid     PidCid
}

func NewDockerApi() (*DockerApi, error) {
	d := DockerApi{}
	var err error
	d.cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return &DockerApi{}, util.ErrorWrapFunc(err)
	}
	d.cidinspect = make(CidInspect)
	d.cidset = make(CidSet)
	d.cidname = make(CidName)
	d.AddNewContainer()
	return &d, nil
}

func (d DockerApi) get1mContainer() ([]types.Container, error) {
	containers, err := d.cli.ContainerList(context.Background(), types.ContainerListOptions{Since: "1m"})
	if err != nil {
		return nil, util.ErrorWrapFunc(err)
	}
	return containers, nil
}

func (d DockerApi) getAllContainer() ([]types.Container, error) {
	containers, err := d.cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, util.ErrorWrapFunc(err)
	}
	return containers, nil
}

// ContainerExecInspect() で事足りるかもしれない
func (d DockerApi) getInspectFromContainer(container types.Container) (types.ContainerJSON, error) {
	cJson, err := d.cli.ContainerInspect(context.Background(), container.ID)
	if err != nil {
		return types.ContainerJSON{}, util.ErrorWrapFunc(err)
	}
	return cJson, nil
}

func (d DockerApi) AddNewContainer() error {
	var cs []types.Container
	var err error
	if len(d.cidset) != 0 {
		cs, err = d.get1mContainer()
	} else {
		cs, err = d.getAllContainer()
	}
	if err != nil {
		return util.ErrorWrapFunc(err)
	}

	for _, c := range cs {
		if _, ok := d.cidset[c.ID]; !ok {
			d.cidset[c.ID] = struct{}{}
			cjson, err := d.getInspectFromContainer(c)
			if err != nil {
				return util.ErrorWrapFunc(err)
			}
			d.cidinspect[c.ID] = cjson
		}
		// add cid to map[name]cid
		for _, name := range c.Names {
			// trim start character "/"
			d.cidname[c.ID] = name[1:]
		}
	}
	log.WithField("DockerApi PidCid len", len(d.Pidcid)).Debug()

	return nil
}

func (d DockerApi) PidFromCid(cid string) (int32, error) {
	if inspect, ok := d.cidinspect[cid]; ok {
		return int32(inspect.State.Pid), nil
	} else {
		return 0, util.ErrorWrapFunc(fmt.Errorf("unknown container id : %s", cid))
	}
}

func (d DockerApi) CidFromPid(pid int) (string, error) {
	if cid, ok := d.Pidcid[pid]; ok {
		return cid, nil
	} else {
		return "", util.ErrorWrapFunc(fmt.Errorf("unknown name"))
	}
}
func (d DockerApi) NameFromCid(cid string) (string, error) {
	if name, ok := d.cidname[cid]; !ok {
		return "", util.ErrorWrapFunc(fmt.Errorf("unkown cid: %v", cid))
	} else {
		return name, nil
	}
}

// func (d DockerApi) PidCid() PidCid {
// 	return d.pidcid
// }
