package docker

import (
	"fmt"
)

var (
	proc                         = "/proc"
	ThisPidIsNotINContainerError = fmt.Errorf("this process is not in container")
)
