package docker

import (
	"fmt"
)

var (
	proc                         = "/proc"
	ThisPidIsNotInContainerError = fmt.Errorf("this process is not in container")
)
