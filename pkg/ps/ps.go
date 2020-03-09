package ps

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/xapima/conps/pkg/util"
)

type StatusValues struct {
	Name      string
	Umask     int32
	State     string
	Tgid      int32
	Pid       int32
	PPid      int32
	TracerPid int32
	Uids      []int32
	Gids      []int32
	Threads   int32
}

// /proc/[number]/task/[number] 以下にも同様の調査をする

func Cmdline(filePath string) ([]string, error) {
	cmdPath := filepath.Join(filePath, "cmdline")
	f, err := os.Open(cmdPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	cmds := bytes.Split(b, []byte{0})
	s := make([]string, len(cmds))
	count := 0
	for _, args := range cmds {
		if len(args) == 0 {
			continue
		}
		s[count] = string(args)
		count++
	}
	return s[:count], nil
}
func Cwd(filePath string) (string, error) {
	cwdPath := filepath.Join(filePath, "cwd")
	realPath, err := readSymlink(cwdPath)
	if err != nil {
		return "", err
	}
	return realPath, nil
}
func Env(filePath string) (map[string]string, error) {
	log.Debug("in Env")
	// envPath := filepath.Join(filePath, "status")
	envPath := filepath.Join(filePath, "environ")
	log.Debugf("do Open: %v", envPath)
	f, err := os.Open(envPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	log.Debug("do ReadAll")
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	log.Debugf("environ_line : %v", string(b))
	envs := bytes.Split(b, []byte{0})
	m := make(map[string]string)
	for _, envString := range envs {
		if len(envString) == 0 {
			continue
		}
		envStrings := strings.Split(string(envString), "=")
		m[envStrings[0]] = envStrings[1]
	}
	return m, nil
}

func Exe(filePath string) (string, error) {
	exePath := filepath.Join(filePath, "exe")
	realPath, err := readSymlink(exePath)
	if err != nil {
		return "", util.ErrorWrapFunc(err)
	}
	return realPath, nil
}

func Fd(filePath string) (map[string]string, error) {
	fdDir := filepath.Join(filePath, "fd")
	fds, err := ioutil.ReadDir(fdDir)
	if err != nil {
		return nil, err
	}
	fdm := make(map[string]string)
	for _, fd := range fds {
		fdm[fd.Name()], err = readSymlink(filepath.Join(fdDir, fd.Name()))
		if err != nil {
			return nil, util.ErrorWrapFunc(err)
		}
	}
	return fdm, nil
}

func Status(filePath string) (StatusValues, error) {
	statusPath := filepath.Join(filePath, "status")
	fp, err := os.Open(statusPath)
	if err != nil {
		return StatusValues{}, util.ErrorWrapFunc(err)
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	sv := StatusValues{}
	for scanner.Scan() {
		line := scanner.Text()
		tabParts := strings.SplitN(line, "\t", 2)
		if len(tabParts) < 2 {
			continue
		}
		value := tabParts[1]
		switch strings.TrimRight(tabParts[0], ":") {
		case "Name":
			sv.Name = strings.Trim(value, " \t")
			// if len(sv.name) >= 15 {
			// 	cmdlineSlice, err := sv.CmdlineSlice()
			// 	if err != nil {
			// 		return err
			// 	}
			// 	if len(cmdlineSlice) > 0 {
			// 		extendedName := filepath.Base(cmdlineSlice[0])
			// 		if strings.HasPrefix(extendedName, sv.name) {
			// 			sv.name = extendedName
			// 		} else {
			// 			sv.name = cmdlineSlice[0]
			// 		}
			// 	}
			// }
		case "Umask":
			v, err := strconv.ParseInt(value, 8, 32)
			if err != nil {
				return StatusValues{}, err
			}
			sv.Umask = int32(v)
		case "State":
			sv.State = value[0:1]
		case "Tgid":
			v, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return StatusValues{}, err
			}
			sv.Tgid = int32(v)
		case "Pid":
			v, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return StatusValues{}, err
			}
			sv.Pid = int32(v)
		case "PPid", "Ppid":
			v, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return StatusValues{}, err
			}
			sv.PPid = int32(v)
		case "TracerPid":
			vv, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return StatusValues{}, err
			}
			sv.TracerPid = int32(vv)
		case "Uid":
			sv.Uids, err = parseTSInt32(value)
			if err != nil {
				return StatusValues{}, err
			}
		case "Gid":
			sv.Gids, err = parseTSInt32(value)
			if err != nil {
				return StatusValues{}, err
			}
		case "Threads":
			vv, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return StatusValues{}, err
			}
			sv.Threads = int32(vv)
		}
	}
	return sv, nil
}

func PPid(pid int) (int, error) {
	ppid, err := PPidWithProcPath(proc, pid)
	if err != nil {
		return 0, util.ErrorWrapFunc(err)
	}
	return ppid, nil
}

func PPidWithProcPath(proc string, pid int) (int, error) {
	statusPath := filepath.Join(proc, strconv.Itoa(pid), "status")
	fp, err := os.Open(statusPath)
	if err != nil {
		return 0, util.ErrorWrapFunc(err)
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		tabParts := strings.SplitN(line, "\t", 2)
		if len(tabParts) < 2 {
			continue
		}
		value := tabParts[1]
		switch strings.TrimRight(tabParts[0], ":") {
		case "PPid", "Ppid":
			v, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return 0, err
			}
			return int(v), nil
		}
	}
	return 0, fmt.Errorf("PPid not found")
}

func GetPidNameSpace(proc string, pid int) (string, error) {
	cgroupPath := filepath.Join(proc, strconv.Itoa(pid), "cgroup")
	fp, err := os.Open(cgroupPath)
	if err != nil {
		return "", util.ErrorWrapFunc(err)
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		colonParts := strings.SplitN(line, ":", 3)
		if len(colonParts) < 3 {
			continue
		}

		switch colonParts[1] {
		case "pids":
			return colonParts[2], nil
		}
	}
	return "", PidNameSpaceNotFoundError
}
