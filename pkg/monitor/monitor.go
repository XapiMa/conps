package monitor

import (
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/xapima/conps/pkg/docker"
	"github.com/xapima/conps/pkg/util"
)

type Monitor struct {
	cache         cache
	pidppid       pidPPidCache
	watcher       *fsnotify.Watcher
	containerApi  *docker.DockerApi
	containerdPid int
	message       chan message
}
type message struct {
	m    messageType
	c    containerInfo
	l    loggingItem
	time time.Time
}

// NewMonitor create new monitor object
func NewMonitor() (*Monitor, error) {
	m := new(Monitor)
	m.cache = make(cache)
	m.pidppid = make(pidPPidCache)
	m.pidppid[0] = &pidItem{pid: 0, ppid: 0, childrenPids: make(map[int]struct{}), containerID: "", containerNames: []string{}, checkedIsContainer: true}
	if containerApi, err := docker.NewDockerApi(); err != nil {
		return nil, util.ErrorWrapFunc(err)
	} else {
		m.containerApi = containerApi
	}
	m.containerdPid = -1
	m.message = make(chan message, 100)
	if err := m.initialWatcher(); err != nil {
		return m, util.ErrorWrapFunc(err)
	}
	return m, nil
}

func (m *Monitor) initialWatcher() error {
	var err error
	m.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return util.ErrorWrapFunc(err)
	}
	return nil
}

// Monitor is processes monitoring function
func (m *Monitor) Monitor() error {
	return m.monitor()
}

func (m *Monitor) monitor() error {
	return m.check()
}
