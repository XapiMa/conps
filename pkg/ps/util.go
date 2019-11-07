package ps

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Cmdline(filePath string) (string, error) {
	cmdPath := filepath.Join(filePath, "cmdline")
	f, err := os.Open(cmdPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(b), nil
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
	envPath := filepath.Join(filePath, "environ")
	f, err := os.Open(envPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	envs := bytes.Split(b, []byte{0})
	m := make(map[string]string)
	for _, envString := range envs {
		if len(envString) == 0 {
			break
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
		return "", err
	}
	return realPath, nil
}

func readSymlink(filePath string) (string, error) {
	info, err := os.Lstat(filePath)
	if err != nil {
		return "", err
	}
	if info.Mode()&os.ModeSymlink != os.ModeSymlink {
		return "", fmt.Errorf("%v is not symlink.", filePath)
	}
	realPath, err := os.Readlink(filePath)
	if err != nil {
		return "", err
	}
	return realPath, nil
}
