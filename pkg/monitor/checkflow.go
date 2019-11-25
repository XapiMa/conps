package monitor

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/xapima/conps/pkg/ps"
	"github.com/xapima/conps/pkg/util"
)

func (m *Monitor) check() error {
	// "/usr/bin/containerd" の Pid を調べる． これをPPidにしているプロセスをみつける．これがコンテナ．

	if err := m.findContainerd(); err != nil {
		return util.ErrorWrapFunc(err)
	}
	log.WithField("containerdPid", m.containerdPid).Debug()
	if err := m.initialWatch(); err != nil {
		return util.ErrorWrapFunc(err)
	}
	if err := m.setContainerPid(); err != nil {
		return util.ErrorWrapFunc(err)
	}

	fmt.Println("Checking processes infomretions")

	return nil
}

// initialWache set pidppid and ContainerProcInformation and FdDirectory
// if you add the directory to watcher, watcher notify information of item in the directory.
func (m *Monitor) initialWatch() error {
	proc := "/proc"
	m.watcher.Add(proc)
	fileinfos, err := ioutil.ReadDir(proc)
	if err != nil {
		return util.ErrorWrapFunc(err)
	}
	for _, fi := range fileinfos {
		if isNum(fi.Name()) {
			pid, _ := strconv.Atoi(fi.Name())
			if err := m.pidppid.add(pid); err != nil {
				return util.ErrorWrapFunc(err)
			}
			if _, err := m.isContainerProc(pid); err != nil {
				return util.ErrorWrapFunc(err)
			}
			path := filepath.Join(proc, fi.Name())
			if err := m.addProcNumber(path); err != nil {
				return util.ErrorWrapFunc(err)
			}
		}
	}
	return nil
}

func (m *Monitor) addProcNumber(path string) error {
	path = filepath.Clean(path)
	m.watcher.Add(path)
	if err := m.addFd(path); err != nil {
		return util.ErrorWrapFunc(err)
	}
	return nil
}

func (m *Monitor) addFd(path string) error {
	fdDir := filepath.Join(path, "fd")
	m.watcher.Add(fdDir)
	// fileinfos, err := ioutil.ReadDir(fdDir)
	// if err != nil {
	// 	return util.ErrorWrapFunc(err)
	// }
	// for _, fi := range fileinfos {
	// 	if isNum(fi.Name()) {
	// 		if isExist(path) {
	// 			m.watcher.Add(path)
	// 		}
	// 	}
	// }
	return nil
}

// 処理をしている途中で containerdが死んだ場合の処理を後で適切な場所に挿入する
// if m.containerdPid == -1{
//
// }

func (m *Monitor) setContainerPid() error {
	m.containerApi.AddNewContainer()
	for pid, cid := range m.containerApi.PidCid() {
		name, err := m.containerApi.NameFromCid(cid)
		if err != nil {
			return util.ErrorWrapFunc(err)
		}
		m.pidppid.addCidName(pid, cid, name)
		log.WithFields(log.Fields{"pid": pid, "cid": cid}).Debug("containerPid")
	}
	return nil
}

func (m *Monitor) isContainerProc(pid int) (bool, error) {
	_, name, err := m.getPPidContainerCidName(pid)
	if err != nil {
		return false, util.ErrorWrapFunc(err)
	}
	return name != "", nil
}

func (m *Monitor) getPPidContainerCidName(pid int) (string, string, error) {
	cid, name, err := m.getPPidContainerCidNameRec(pid, true)
	if err != nil {
		return "", "", util.ErrorWrapFunc(err)
	}
	return cid, name, nil
}

func (m *Monitor) getPPidContainerCidNameRec(pid int, first bool) (string, string, error) {
	pc, ok := m.pidppid[pid]
	if !ok {
		if err := m.pidppid.add(pid); err != nil {
			return "", "", util.ErrorWrapFunc(err)
		}
		pc, _ = m.pidppid[pid]
	} else if pc.ppid == -1 {
		if err := m.pidppid.add(pid); err != nil {
			return "", "", util.ErrorWrapFunc(err)
		}
		pc, _ = m.pidppid[pid]
	}
	if pc.checkedIsContainer {
		return pc.containerID, pc.containerName, nil
	} else {
		var name string
		var cid string
		var err error
		if pc.ppid == m.containerdPid {
			if first {
				if err := m.setContainerPid(); err != nil {
					return "", "", util.ErrorWrapFunc(err)
				}
				cid, name, err := m.getPPidContainerCidNameRec(pid, false)
				if err != nil {
					return "", "", util.ErrorWrapFunc(err)
				}
				return cid, name, nil
			}
		} else {
			cid, name, err = m.getPPidContainerCidName(pc.ppid)
			if err != nil {
				return "", "", util.ErrorWrapFunc(err)
			}
		}
		pc.checkedIsContainer = true
		pc.containerName = name
		pc.containerID = cid
		return cid, name, nil
	}
}

func isContainerd(filePath string) bool {
	real, err := ps.Exe(filePath)
	if err != nil {
		return false
	}
	return real == containerdPath
}

func (m *Monitor) findContainerd() error {
	proc := "/proc"
	fileinfos, err := ioutil.ReadDir(proc)
	if err != nil {
		return util.ErrorWrapFunc(err)
	}
	for _, fi := range fileinfos {
		if !isNum(fi.Name()) {
			continue
		}
		if ok := isContainerd(filepath.Join(proc, fi.Name())); ok {
			m.containerdPid, err = strconv.Atoi(fi.Name())
			if err != nil {
				return util.ErrorWrapFunc(err)
			}
			return nil
		}
	}
	return fmt.Errorf("containerd not found")
}

func (m *Monitor) setPidCid() error {
	if err := m.containerApi.AddNewContainer(); err != nil {
		return util.ErrorWrapFunc(err)
	}
	for pid, cid := range m.containerApi.PidCid() {
		if _, ok := m.pidppid[pid]; !ok {
			m.pidppid[pid] = newPidItem()
		}
		m.pidppid[pid].pid = pid
		m.pidppid[pid].ppid = m.containerdPid
		m.pidppid[pid].containerID = cid
		if name, err := m.containerApi.NameFromCid(cid); err != nil {
			return util.ErrorWrapFunc(err)
		} else {
			m.pidppid[pid].containerName = name
		}
		m.pidppid[pid].checkedIsContainer = true
	}
	return nil
}
