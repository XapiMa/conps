package ps

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type PasswdValues struct {
	Name     string
	Password string
	Uid      int32
	Gid      int32
	Home     string
	Comment  string
	Shell    string
}

type PasswdMap map[int32]PasswdValues

type GroupValues struct {
	Name     string
	Password string
	Gid      int32
	Users    []string
}

type GroupMap map[int32]GroupValues

func GetPasswdMap(rootPath string) (PasswdMap, error) {
	passwdPath := filepath.Join(rootPath, "etc", "passwd")
	fp, err := os.Open(passwdPath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	m := make(PasswdMap)
	for scanner.Scan() {
		if err := m.addLine(scanner.Text()); err != nil {
			return nil, err
		}
	}
	return m, nil
}

func parsePasswdLine(line string) (PasswdValues, error) {
	p := strings.Split(line, ":")
	if len(p) != 7 {
		return PasswdValues{}, fmt.Errorf("Not passwd line : %s", line)
	}
	uid, err := strconv.ParseInt(p[2], 10, 32)
	if err != nil {
		return PasswdValues{}, err
	}
	gid, err := strconv.ParseInt(p[3], 10, 32)
	if err != nil {
		return PasswdValues{}, err
	}
	return PasswdValues{Name: p[0], Password: p[1], Uid: int32(uid), Gid: int32(gid), Home: p[4], Comment: p[5], Shell: p[6]}, nil
}

func (m PasswdMap) addLine(line string) error {
	v, err := parsePasswdLine(line)
	if err != nil {
		return err
	}
	m[v.Uid] = v
	return nil
}

func GetGroupMap(rootPath string) (GroupMap, error) {
	groupPath := filepath.Join(rootPath, "etc", "group")
	fp, err := os.Open(groupPath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	m := make(GroupMap)
	for scanner.Scan() {
		if err := m.addLine(scanner.Text()); err != nil {
			return nil, err
		}
	}
	return m, nil
}

func parseGroupLine(line string) (GroupValues, error) {
	p := strings.Split(line, ":")
	if len(p) != 4 {
		return GroupValues{}, fmt.Errorf("Not group line : %s", line)
	}
	gid, err := strconv.ParseInt(p[2], 10, 32)
	if err != nil {
		return GroupValues{}, err
	}
	users := []string{}
	if p[3] != "" {
		users = strings.Split(p[3], ",")
	}
	return GroupValues{Name: p[0], Password: p[1], Gid: int32(gid), Users: users}, nil
}

func (m GroupMap) addLine(line string) error {
	v, err := parseGroupLine(line)
	if err != nil {
		return err
	}
	m[v.Gid] = v
	return nil
}
