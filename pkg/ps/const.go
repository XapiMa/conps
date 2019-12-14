package ps

import "fmt"

var (
	PidNameSpaceNotFoundError = fmt.Errorf("pid name space is not found")

	proc = "/proc"
)
