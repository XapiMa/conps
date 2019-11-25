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

func TestAdd(t *testing.T) {
	m, err := NewMonitor()
	if err != nil {
		t.Errorf("cant new monitor: %v\n", err)
	}
	err = m.pidppid.add(1)
	if err != nil {
		t.Error(err)
	}
}
