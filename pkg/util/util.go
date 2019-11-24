package util

import (
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

func ErrorWrapFunc(err error) error {
	programCounter, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(programCounter)
	funcNames := strings.Split(fn.Name(), ".")
	return errors.Wrap(err, funcNames[len(funcNames)-1])
}
