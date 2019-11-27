package ps

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}
func TestGetPasswdMap(t *testing.T) {
	type args struct {
		rootPath string
	}
	tests := []struct {
		name    string
		args    args
		want    PasswdMap
		wantErr bool
	}{
		{"in container", args{"/"}, PasswdMap{
			0:     PasswdValues{Name: "root", Password: "x", Uid: 0, Gid: 0, Home: "root", Comment: "/root", Shell: "/bin/ash"},
			1:     PasswdValues{Name: "bin", Password: "x", Uid: 1, Gid: 1, Home: "bin", Comment: "/bin", Shell: "/sbin/nologin"},
			2:     PasswdValues{Name: "daemon", Password: "x", Uid: 2, Gid: 2, Home: "daemon", Comment: "/sbin", Shell: "/sbin/nologin"},
			3:     PasswdValues{Name: "adm", Password: "x", Uid: 3, Gid: 4, Home: "adm", Comment: "/var/adm", Shell: "/sbin/nologin"},
			4:     PasswdValues{Name: "lp", Password: "x", Uid: 4, Gid: 7, Home: "lp", Comment: "/var/spool/lpd", Shell: "/sbin/nologin"},
			5:     PasswdValues{Name: "sync", Password: "x", Uid: 5, Gid: 0, Home: "sync", Comment: "/sbin", Shell: "/bin/sync"},
			6:     PasswdValues{Name: "shutdown", Password: "x", Uid: 6, Gid: 0, Home: "shutdown", Comment: "/sbin", Shell: "/sbin/shutdown"},
			7:     PasswdValues{Name: "halt", Password: "x", Uid: 7, Gid: 0, Home: "halt", Comment: "/sbin", Shell: "/sbin/halt"},
			8:     PasswdValues{Name: "mail", Password: "x", Uid: 8, Gid: 12, Home: "mail", Comment: "/var/spool/mail", Shell: "/sbin/nologin"},
			9:     PasswdValues{Name: "news", Password: "x", Uid: 9, Gid: 13, Home: "news", Comment: "/usr/lib/news", Shell: "/sbin/nologin"},
			10:    PasswdValues{Name: "uucp", Password: "x", Uid: 10, Gid: 14, Home: "uucp", Comment: "/var/spool/uucppublic", Shell: "/sbin/nologin"},
			11:    PasswdValues{Name: "operator", Password: "x", Uid: 11, Gid: 0, Home: "operator", Comment: "/root", Shell: "/sbin/nologin"},
			13:    PasswdValues{Name: "man", Password: "x", Uid: 13, Gid: 15, Home: "man", Comment: "/usr/man", Shell: "/sbin/nologin"},
			14:    PasswdValues{Name: "postmaster", Password: "x", Uid: 14, Gid: 12, Home: "postmaster", Comment: "/var/spool/mail", Shell: "/sbin/nologin"},
			16:    PasswdValues{Name: "cron", Password: "x", Uid: 16, Gid: 16, Home: "cron", Comment: "/var/spool/cron", Shell: "/sbin/nologin"},
			21:    PasswdValues{Name: "ftp", Password: "x", Uid: 21, Gid: 21, Home: "", Comment: "/var/lib/ftp", Shell: "/sbin/nologin"},
			22:    PasswdValues{Name: "sshd", Password: "x", Uid: 22, Gid: 22, Home: "sshd", Comment: "/dev/null", Shell: "/sbin/nologin"},
			25:    PasswdValues{Name: "at", Password: "x", Uid: 25, Gid: 25, Home: "at", Comment: "/var/spool/cron/atjobs", Shell: "/sbin/nologin"},
			31:    PasswdValues{Name: "squid", Password: "x", Uid: 31, Gid: 31, Home: "Squid", Comment: "/var/cache/squid", Shell: "/sbin/nologin"},
			33:    PasswdValues{Name: "xfs", Password: "x", Uid: 33, Gid: 33, Home: "X Font Server", Comment: "/etc/X11/fs", Shell: "/sbin/nologin"},
			35:    PasswdValues{Name: "games", Password: "x", Uid: 35, Gid: 35, Home: "games", Comment: "/usr/games", Shell: "/sbin/nologin"},
			70:    PasswdValues{Name: "postgres", Password: "x", Uid: 70, Gid: 70, Home: "", Comment: "/var/lib/postgresql", Shell: "/bin/sh"},
			85:    PasswdValues{Name: "cyrus", Password: "x", Uid: 85, Gid: 12, Home: "", Comment: "/usr/cyrus", Shell: "/sbin/nologin"},
			89:    PasswdValues{Name: "vpopmail", Password: "x", Uid: 89, Gid: 89, Home: "", Comment: "/var/vpopmail", Shell: "/sbin/nologin"},
			123:   PasswdValues{Name: "ntp", Password: "x", Uid: 123, Gid: 123, Home: "NTP", Comment: "/var/empty", Shell: "/sbin/nologin"},
			209:   PasswdValues{Name: "smmsp", Password: "x", Uid: 209, Gid: 209, Home: "smmsp", Comment: "/var/spool/mqueue", Shell: "/sbin/nologin"},
			405:   PasswdValues{Name: "guest", Password: "x", Uid: 405, Gid: 100, Home: "guest", Comment: "/dev/null", Shell: "/sbin/nologin"},
			65534: PasswdValues{Name: "nobody", Password: "x", Uid: 65534, Gid: 65534, Home: "nobody", Comment: "/", Shell: "/sbin/nologin"},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPasswdMap(tt.args.rootPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPasswdMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// printPasswdMap(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPasswdMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func printPasswdMap(got PasswdMap) {
	strs := []string{}
	for k, v := range got {
		strs = append(strs, fmt.Sprintf("%d:PasswdValues{Name:\"%s\",Password:\"%s\",Uid:%d,Gid:%d,Home:\"%s\",Comment:\"%s\",Shell:\"%s\"},\n", k, v.Name, v.Password, v.Uid, v.Gid, v.Home, v.Comment, v.Shell))
	}
	sort.Slice(strs, func(i, j int) bool {
		a, _ := strconv.ParseInt(strings.Split(strs[i], ":")[0], 10, 32)
		b, _ := strconv.ParseInt(strings.Split(strs[j], ":")[0], 10, 32)
		return a < b
	})
	for _, str := range strs {
		fmt.Print(str)
	}
}

func Test_getGroupMap(t *testing.T) {
	type args struct {
		rootPath string
	}
	tests := []struct {
		name    string
		args    args
		want    GroupMap
		wantErr bool
	}{
		{"in container", args{"/"}, GroupMap{
			0:     GroupValues{Name: "root", Password: "x", Gid: 0, Users: []string{"root"}},
			1:     GroupValues{Name: "bin", Password: "x", Gid: 1, Users: []string{"root", "bin", "daemon"}},
			2:     GroupValues{Name: "daemon", Password: "x", Gid: 2, Users: []string{"root", "bin", "daemon"}},
			3:     GroupValues{Name: "sys", Password: "x", Gid: 3, Users: []string{"root", "bin", "adm"}},
			4:     GroupValues{Name: "adm", Password: "x", Gid: 4, Users: []string{"root", "adm", "daemon"}},
			5:     GroupValues{Name: "tty", Password: "x", Gid: 5, Users: []string{}},
			6:     GroupValues{Name: "disk", Password: "x", Gid: 6, Users: []string{"root", "adm"}},
			7:     GroupValues{Name: "lp", Password: "x", Gid: 7, Users: []string{"lp"}},
			8:     GroupValues{Name: "mem", Password: "x", Gid: 8, Users: []string{}},
			9:     GroupValues{Name: "kmem", Password: "x", Gid: 9, Users: []string{}},
			10:    GroupValues{Name: "wheel", Password: "x", Gid: 10, Users: []string{"root"}},
			11:    GroupValues{Name: "floppy", Password: "x", Gid: 11, Users: []string{"root"}},
			12:    GroupValues{Name: "mail", Password: "x", Gid: 12, Users: []string{"mail"}},
			13:    GroupValues{Name: "news", Password: "x", Gid: 13, Users: []string{"news"}},
			14:    GroupValues{Name: "uucp", Password: "x", Gid: 14, Users: []string{"uucp"}},
			15:    GroupValues{Name: "man", Password: "x", Gid: 15, Users: []string{"man"}},
			16:    GroupValues{Name: "cron", Password: "x", Gid: 16, Users: []string{"cron"}},
			17:    GroupValues{Name: "console", Password: "x", Gid: 17, Users: []string{}},
			18:    GroupValues{Name: "audio", Password: "x", Gid: 18, Users: []string{}},
			19:    GroupValues{Name: "cdrom", Password: "x", Gid: 19, Users: []string{}},
			20:    GroupValues{Name: "dialout", Password: "x", Gid: 20, Users: []string{"root"}},
			21:    GroupValues{Name: "ftp", Password: "x", Gid: 21, Users: []string{}},
			22:    GroupValues{Name: "sshd", Password: "x", Gid: 22, Users: []string{}},
			23:    GroupValues{Name: "input", Password: "x", Gid: 23, Users: []string{}},
			25:    GroupValues{Name: "at", Password: "x", Gid: 25, Users: []string{"at"}},
			26:    GroupValues{Name: "tape", Password: "x", Gid: 26, Users: []string{"root"}},
			27:    GroupValues{Name: "video", Password: "x", Gid: 27, Users: []string{"root"}},
			28:    GroupValues{Name: "netdev", Password: "x", Gid: 28, Users: []string{}},
			30:    GroupValues{Name: "readproc", Password: "x", Gid: 30, Users: []string{}},
			31:    GroupValues{Name: "squid", Password: "x", Gid: 31, Users: []string{"squid"}},
			33:    GroupValues{Name: "xfs", Password: "x", Gid: 33, Users: []string{"xfs"}},
			34:    GroupValues{Name: "kvm", Password: "x", Gid: 34, Users: []string{"kvm"}},
			35:    GroupValues{Name: "games", Password: "x", Gid: 35, Users: []string{}},
			42:    GroupValues{Name: "shadow", Password: "x", Gid: 42, Users: []string{}},
			70:    GroupValues{Name: "postgres", Password: "x", Gid: 70, Users: []string{}},
			80:    GroupValues{Name: "cdrw", Password: "x", Gid: 80, Users: []string{}},
			85:    GroupValues{Name: "usb", Password: "x", Gid: 85, Users: []string{}},
			89:    GroupValues{Name: "vpopmail", Password: "x", Gid: 89, Users: []string{}},
			100:   GroupValues{Name: "users", Password: "x", Gid: 100, Users: []string{"games"}},
			123:   GroupValues{Name: "ntp", Password: "x", Gid: 123, Users: []string{}},
			200:   GroupValues{Name: "nofiles", Password: "x", Gid: 200, Users: []string{}},
			209:   GroupValues{Name: "smmsp", Password: "x", Gid: 209, Users: []string{"smmsp"}},
			245:   GroupValues{Name: "locate", Password: "x", Gid: 245, Users: []string{}},
			300:   GroupValues{Name: "abuild", Password: "x", Gid: 300, Users: []string{}},
			406:   GroupValues{Name: "utmp", Password: "x", Gid: 406, Users: []string{}},
			999:   GroupValues{Name: "ping", Password: "x", Gid: 999, Users: []string{}},
			65533: GroupValues{Name: "nogroup", Password: "x", Gid: 65533, Users: []string{}},
			65534: GroupValues{Name: "nobody", Password: "x", Gid: 65534, Users: []string{}},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGroupMap(tt.args.rootPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("getGroupMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// printGroupMap(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getGroupMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func printGroupMap(got GroupMap) {
	strs := []string{}
	for k, v := range got {
		users := ""
		for _, name := range v.Users {
			users += "\"" + name + "\","
		}
		if len(users) != 0 {
			users = users[:len(users)-1]
		}
		strs = append(strs, fmt.Sprintf("%d:GroupValues{Name:\"%s\",Password:\"%s\",Gid:%d,Users:[]string{%s}},\n", k, v.Name, v.Password, v.Gid, users))
	}
	sort.Slice(strs, func(i, j int) bool {
		a, _ := strconv.ParseInt(strings.Split(strs[i], ":")[0], 10, 32)
		b, _ := strconv.ParseInt(strings.Split(strs[j], ":")[0], 10, 32)
		return a < b
	})
	for _, str := range strs {
		fmt.Print(str)
	}
}
