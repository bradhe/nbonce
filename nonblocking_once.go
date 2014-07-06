package nbonce

import (
	"time"
	"sync"
	"sync/atomic"
)

const (
	// indicates that the function has not yet been started.
	notstarted		= 0

	// indicates that the function has been started.
	started				= 1

	// indicates that the function is not yet complete.
	notfinished		= 0

	// indicates that the function is indeed complete.
	finished			= 1
)

type NonblockingOnce struct {
	Resetable bool

	// Flags to indicate what state things are in.
	started, finished int32

	lock sync.Mutex
}

// Schedules a function to run one time in a goroutine if it hasn't been
// scheduled already. If it has then it will return and not re-schedule the
// function.
func (self *NonblockingOnce) Do(f func()) {
	if atomic.CompareAndSwapInt32(&self.started, notstarted, started) {
		go self.doFunc(f)
	}
}

// If the function is running, waits for it to complete before returning. If it
// hasn't been started yet at all then it will just return immediately.
func (self *NonblockingOnce) Wait() {
	// If it wasn't started there is nothing to wait for.
	if atomic.LoadInt32(&self.started) == notstarted {
		return
	}

	for {
		self.lock.Lock()
		if self.finished == finished {
			break
		}
		self.lock.Unlock()

		// Not sure if this is strictly nescessary, just don't really want to burn
		// up the CPU ya know?
		time.Sleep(0)
	}
}

func (self *NonblockingOnce) doFunc(f func()) {
	atomic.StoreInt32(&self.finished, notfinished)

	// We do this to ensure that it always performs this action, even if f panics
	// for whatever reason but doesn't kill this fun. There might be a better way
	// to do this, and might not be nescessary at all.
	defer func() {
		atomic.StoreInt32(&self.finished, finished)
	}()

	// off to the races!
	f()

	if self.Resetable {
		atomic.StoreInt32(&self.started, notstarted)
	}
}
