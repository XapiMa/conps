package docker

import (
	"fmt"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"github.com/xapima/conps/pkg/ps"
	"github.com/xapima/conps/pkg/util"
	"golang.org/x/net/context"
)

type CidInspect map[string]types.ContainerJSON
type CidSet map[string]struct{}
type CidNameSet map[string]map[string]struct{}
type PidCid map[int]string
type CidPasswdMap map[string]ps.PasswdMap

type DockerApi struct {
	cli          *client.Client
	cidinspect   CidInspect
	cidset       CidSet
	cidnameSet   CidNameSet
	cidPasswdMap CidPasswdMap
	pidcid       PidCid
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
	d.cidnameSet = make(CidNameSet)
	d.pidcid = make(PidCid)
	d.cidPasswdMap = make(CidPasswdMap)
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
		// add Names to map[cid]NameSet
		if _, ok := d.cidnameSet[c.ID]; !ok {
			d.cidnameSet[c.ID] = make(map[string]struct{})
		}
		for _, name := range c.Names {
			// trim start character "/"
			d.cidnameSet[c.ID][name[1:]] = struct{}{}
		}
		var err error
		var pid int
		pid, err = d.PidFromCid(c.ID)
		if err != nil {
			return util.ErrorWrapFunc(err)
		}
		d.pidcid[pid] = c.ID

	}
	for k, v := range d.pidcid {
		log.WithFields(log.Fields{"pid": k, "cid": v}).Debug("DockerApi PidCid")
	}

	return nil
}

func (d DockerApi) PidFromCid(cid string) (int, error) {
	inspect, ok := d.cidinspect[cid]
	if ok {
		pid := inspect.State.Pid
		cpid, err := ps.PPid(pid)
		if err != nil {
			return 0, util.ErrorWrapFunc(err)
		}
		return cpid, nil
	}

	return 0, util.ErrorWrapFunc(fmt.Errorf("unknown container id : %s", cid))
}

func (d DockerApi) CidFromPid(pid int) (string, error) {
	if cid, ok := d.pidcid[pid]; ok {
		return cid, nil
	}
	return "", util.ErrorWrapFunc(fmt.Errorf("unknown name"))
}

func (d DockerApi) NamesFromCid(cid string) ([]string, error) {
	if err := d.AddNewContainer(); err != nil {
		return nil, util.ErrorWrapFunc(err)
	}

	nameSet, ok := d.cidnameSet[cid]
	if !ok {
		return nil, util.ErrorWrapFunc(fmt.Errorf("unkown cid: %v", cid))
	}
	names := make([]string, 0, len(nameSet))
	for k, _ := range nameSet {
		names = append(names, k)
	}
	return names, nil
}

func (d *DockerApi) CidFromName(name string) (string, error) {
	if err := d.AddNewContainer(); err != nil {
		return "", util.ErrorWrapFunc(err)
	}
	for cid, nameSet := range d.cidnameSet {
		for containerName, _ := range nameSet {
			if containerName == name {
				return cid, nil
			}
		}
	}
	return "", util.ErrorWrapFunc(fmt.Errorf("unknown container name: %v", name))
}

func (d DockerApi) PidCid() PidCid {
	return d.pidcid
}

func (d *DockerApi) GetUserNameFromCidUid(cid string, uid int) (string, error) {
	if _, ok := d.cidPasswdMap[cid]; !ok {
		if err := d.setCidPasswdMap(cid); err != nil {
			return "", util.ErrorWrapFunc(err)
		}
	}
	m := d.cidPasswdMap[cid]
	if item, ok := m[int32(uid)]; !ok {
		return "", fmt.Errorf("unknown uid :%v in container: %v", uid, cid)
	} else {
		return item.Name, nil
	}
}

func (d *DockerApi) setCidPasswdMap(cid string) error {
	containerRoot, err := d.ContainerPathToHostPath(cid, "/")
	if err != nil {
		return util.ErrorWrapFunc(err)
	}
	m, err := ps.GetUidNameMap(containerRoot)
	if err != nil {
		return util.ErrorWrapFunc(err)
	}
	d.cidPasswdMap[cid] = m
	return nil
}

func (d *DockerApi) ContainerPathToHostPath(cid string, path string) (string, error) {
	if err := d.AddNewContainer(); err != nil {
		return "", util.ErrorWrapFunc(err)
	}
	merged := d.cidinspect[cid].ContainerJSONBase.GraphDriver.Data["MergedDir"]
	return filepath.Clean(filepath.Join(merged, path)), nil
}
