package docker

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"github.com/xapima/conps/pkg/ps"
	"github.com/xapima/conps/pkg/util"
	"golang.org/x/net/context"
)

type CidInspect map[string]types.ContainerJSON
type CidSet map[string]struct{}
type NameCid map[string]string
type PidCid map[int]string
type CidPasswdMap map[string]ps.PasswdMap
type CidGroupMap map[string]ps.GroupMap

type DockerApi struct {
	cli          *client.Client
	cidinspect   CidInspect
	namecid      NameCid
	cidPasswdMap CidPasswdMap
	cidGroupMap  CidGroupMap
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
	d.namecid = make(NameCid)
	d.pidcid = make(PidCid)
	d.cidPasswdMap = make(CidPasswdMap)
	d.cidGroupMap = make(CidGroupMap)
	return &d, nil
}

func (d *DockerApi) IsContainerProcess(pid int) bool {
	if _, ok := d.pidcid[pid]; ok {
		return true
	}
	err := d.setCidWithPid(pid)
	if err != nil {
		return false
	}
	return true
}

func (d *DockerApi) setCidWithPid(pid int) error {
	cid, err := d.getCidFromCgroupWithPid(pid)
	if err != nil {
		return util.ErrorWrapFunc(err)
	}
	if err := d.setCid(cid); err != nil {
		return util.ErrorWrapFunc(err)
	}
	d.pidcid[pid] = cid
	return nil
}

func (d DockerApi) getCidFromCgroupWithPid(pid int) (string, error) {
	nameSpace, err := ps.GetPidNameSpace(proc, pid)
	if err == ps.PidNameSpaceNotFoundError {
		log.Debug("PidNameSpaceNotFound error cause")
		return "", PidIsNotInContainerError
	} else if err != nil {
		return "", util.ErrorWrapFunc(err)
	}
	slashParts := strings.Split(nameSpace, "/")
	if len(slashParts) < 3 {
		return "", PidIsNotInContainerError
	}
	log.WithFields(log.Fields{"slashPart0": slashParts[0], "slashPart1": slashParts[1]}).Debug("SlashParts")
	if slashParts[1] == "docker" {
		return slashParts[2], nil
	}
	return "", PidIsNotInContainerError

}

func (d *DockerApi) setCid(cid string) error {
	if _, ok := d.cidinspect[cid]; ok {
		return nil
	}

	cjson, err := d.getInspectWithCid(cid)
	if err != nil {
		return util.ErrorWrapFunc(err)
	}
	d.cidinspect[cid] = cjson
	name := cjson.ContainerJSONBase.Name
	// trim start character "/"
	d.namecid[name[1:]] = cid
	if err := d.setPasswdMapWithCid(cid); err != nil {
		return util.ErrorWrapFunc(err)
	}
	if err := d.setGroupMap(cid); err != nil {
		return util.ErrorWrapFunc(err)
	}

	return nil
}

func (d DockerApi) getInspectWithCid(cid string) (types.ContainerJSON, error) {
	if cjson, ok := d.cidinspect[cid]; ok {
		return cjson, nil
	}
	cJson, err := d.cli.ContainerInspect(context.Background(), cid)
	if err != nil {
		return types.ContainerJSON{}, util.ErrorWrapFunc(err)
	}
	return cJson, nil
}

func (d *DockerApi) setPasswdMapWithCid(cid string) error {
	containerRoot, err := d.ContainerPathToHostPath(cid, "/")
	if err != nil {
		return util.ErrorWrapFunc(err)
	}
	m, err := ps.GetPasswdMap(containerRoot)
	if err != nil {
		return util.ErrorWrapFunc(err)
	}
	d.cidPasswdMap[cid] = m
	return nil
}

func (d *DockerApi) setGroupMap(cid string) error {
	containerRoot, err := d.ContainerPathToHostPath(cid, "/")
	if err != nil {
		return util.ErrorWrapFunc(err)
	}
	m, err := ps.GetGroupMap(containerRoot)
	if err != nil {
		return util.ErrorWrapFunc(err)
	}
	d.cidGroupMap[cid] = m
	return nil
}

func (d DockerApi) GetPidWithCid(cid string) (int, error) {
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

func (d DockerApi) GetCidWithPid(pid int) (string, error) {
	if cid, ok := d.pidcid[pid]; ok {
		return cid, nil
	}
	return "", util.ErrorWrapFunc(fmt.Errorf("unknown name"))
}

func (d DockerApi) GetNameWithCid(cid string) (string, error) {
	if cjson, ok := d.cidinspect[cid]; !ok {
		return "", util.ErrorWrapFunc(fmt.Errorf("unkown cid: %v", cid))
	} else {
		// trim start character "/"
		return cjson.ContainerJSONBase.Name[1:], nil
	}
}

func (d *DockerApi) GetCidWithName(name string) (string, error) {
	if cid, ok := d.namecid[name]; ok {
		return cid, nil
	}
	return "", util.ErrorWrapFunc(fmt.Errorf("unknown container name: %v", name))
}

func (d DockerApi) PidCid() PidCid {
	return d.pidcid
}

func (d *DockerApi) GetUserNameWithCidUid(cid string, uid int) (string, error) {
	if _, ok := d.cidinspect[cid]; !ok {
		return "", util.ErrorWrapFunc(unknownCid(cid))
	}
	m := d.cidPasswdMap[cid]
	if item, ok := m[int32(uid)]; !ok {
		return "", fmt.Errorf("unknown uid: %v in container: %v", uid, cid)
	} else {
		return item.Name, nil
	}
}

func (d *DockerApi) GetGroupNameWithCidGid(cid string, gid int) (string, error) {
	if _, ok := d.cidinspect[cid]; !ok {
		return "", util.ErrorWrapFunc(unknownCid(cid))
	}
	m := d.cidGroupMap[cid]
	if item, ok := m[int32(gid)]; !ok {
		return "", fmt.Errorf("unknown gid: %v in container: %v", gid, cid)
	} else {
		return item.Name, nil
	}
}

func (d *DockerApi) ContainerPathToHostPath(cid string, path string) (string, error) {
	if _, ok := d.cidinspect[cid]; !ok {
		return "", util.ErrorWrapFunc(unknownCid(cid))
	}
	merged := d.cidinspect[cid].ContainerJSONBase.GraphDriver.Data["MergedDir"]
	return filepath.Clean(filepath.Join(merged, path)), nil
}
