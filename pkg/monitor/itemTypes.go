package monitor

import (
	"github.com/xapima/conps/pkg/ps"

	"github.com/docker/docker/api/types"
)

type messageType string

const (
	ErrorMessage   = messageType("ERROR")
	WarningMessage = messageType("WARNIG")
	ActiveMessage  = messageType("ACTIVE")
)

type loggingItem struct {
	cmdline      string
	cwd          string
	env          map[string]string
	exe          string
	fd           map[string]string
	statusValues ps.StatusValues
}

type FilterItem struct {
	exec string
	cmd  string
	open []OpenFile
	user string
	pid  int
}

type OpenFile struct {
	Path string `json:"path"`
	Fd   int    `json:"fd"`
}

type containerInfo types.ContainerJSON
