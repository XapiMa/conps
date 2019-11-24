package monitor

import (
	"fmt"

	"github.com/xapima/conps/pkg/ps"
	"github.com/xapima/conps/pkg/util"
)

type cache map[int]cacheItem
type cacheItem struct {
	f FilterItem
	l loggingItem
}

type pidPPidCache map[int]*pidItem
type pidItem struct {
	pid           int
	ppid          int
	childrenPids  map[int]struct{}
	containerName string
	containerID   string
	// if checked : true , else : false
	checkedIsContainer bool
}

func (c cache) in(pid int) bool {
	_, ok := c[pid]
	return ok
}

func (c cache) add(item cacheItem) error {
	pid := item.f.pid
	if _, ok := c[pid]; ok {
		return fmt.Errorf("pid:%d is already cached as %v. item is %v", pid, c[pid], item)
	}
	c[pid] = item
	return nil
}

func (c cache) del(pids []int) error {
	notCachePid := []int{}
	for _, pid := range pids {
		if _, ok := c[pid]; !ok {
			notCachePid = append(notCachePid, pid)
		} else {
			delete(c, pid)
		}
	}
	if len(notCachePid) != 0 {
		return fmt.Errorf("pids %v is not cached", notCachePid)
	}
	return nil
}

func (c pidPPidCache) add(pid int) error {
	if pid == 0 {
		return nil
	}
	if _, ok := c[pid]; ok {
		// default ppid : -1
		if c[pid].ppid != -1 {
			// 同じpidのppidが後で変わることがあるのか？
			// daemon化の場合は？
			return nil
		}
	}
	ppid, err := ps.PPid(proc, pid)
	if err != nil {
		return util.ErrorWrapFunc(err)
	}
	c[pid] = newPidItem()
	c[pid].pid = pid
	c[pid].ppid = ppid
	if _, ok := c[ppid]; !ok {
		c[ppid] = newPidItem()
	}
	c[ppid].childrenPids[pid] = struct{}{}

	if err := c.add(ppid); err != nil {
		return err
	}
	return nil
}

func newPidItem() *pidItem {
	// set default ppid : -1
	return &pidItem{ppid: -1, childrenPids: make(map[int]struct{}), checkedIsContainer: false}
}

// func (c cacheCount) in(t *FilterItem) bool {
// 	count, ok := c[t]
// 	if ok {
// 		if count != 0 {
// 			return true
// 		}
// 	}
// 	return false
// }

// func (c cacheCount) is(t *Target) bool {
// 	_, ok := c[t]
// 	return ok
// }

// func (c cacheCount) add(t *Target) error {
// 	if _, ok := c[t]; !ok {
// 		return fmt.Errorf("%v is not key of cacheCount", *t)
// 	}
// 	// もしも 0 → 1 ならメッセージを出す
// 	c[t]++
// 	return nil
// }
// func (c cacheCount) del(t *Target) error {
// 	if c[t] == 0 {
// 		return fmt.Errorf("%v %v is not cached", c, *c[t])
// 	}
// 	c[t]--
// 	if c[t] == 0 {
// 		delete(c, t)
// 		// メッセージを出す
// 	}
// }
