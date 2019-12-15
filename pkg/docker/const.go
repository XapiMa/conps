package docker

import (
	"fmt"
)

var (
	proc                     = "/proc"
	PidIsNotInContainerError = fmt.Errorf("this process is not in container")
	UnknownCidError          = fmt.Errorf("unkown cid")
)
