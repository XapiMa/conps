package monitor

import (
	"testing"
)

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
