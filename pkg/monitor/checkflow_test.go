package monitor

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func init() {
	log.SetLevel(log.DebugLevel)
}
func TestInit(t *testing.T) {
	m, err := NewMonitor()
	if err != nil {
		t.Errorf("cant new monitor: %v\n", err)
	}
	err = m.check()
	if err != nil {
		t.Errorf("cant check: %v\n", err)
	}
}

func Test0(t *testing.T) {
	m, err := NewMonitor()
	if err != nil {
		t.Errorf("cant new monitor: %v\n", err)
	}
	t.Errorf("pidppid[0] : %v\n", *m.pidppid[0])
}
