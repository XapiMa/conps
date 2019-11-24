package monitor

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func isExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
func isNum(s string) bool {
	if _, err := strconv.Atoi(s); err != nil {
		return false
	}
	return true
}

func isFdDir(path string) bool {
	path = filepath.Clean(path)
	p := strings.Split(path, "/")
	if len(p) != 4 {
		return false
	}
	if p[1] == "proc" && isNum(p[2]) && p[3] == "fd" {
		return true
	}
	return false
}

func isFd(path string) bool {
	path = filepath.Clean(path)
	p := strings.Split(path, "/")
	if len(p) != 5 {
		return false
	}
	if p[1] == "proc" && isNum(p[2]) && p[3] == "fd" && isNum(p[4]) {
		return true
	}
	return false
}

func isProcDir(path string) bool {
	path = filepath.Clean(path)
	dir, base := filepath.Split(path)

	if filepath.Clean(dir) != "/proc" {
		return false
	}
	return isNum(base)
}

func getPidFromProcPath(path string) (int, error) {
	path = filepath.Clean(path)
	p := strings.Split(path, "/")
	if len(p) < 3 {
		return 0, fmt.Errorf("%s is not proc path", path)
	}
	if p[1] != "proc" {
		return 0, fmt.Errorf("%s is not proc path", path)
	}
	pid, err := strconv.Atoi(p[2])
	if err != nil {
		return 0, fmt.Errorf("%s is not proc path", path)
	}
	return pid, nil
}

// func inContainerProc(pid int) (bool, error) {

// }
