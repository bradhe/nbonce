package nbonce

import (
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
	Resettable bool

	// Flags to indicate what state things are in.
	started int32

	wg sync.WaitGroup
}

// Schedules a function to run one time in a goroutine if it hasn't been
// scheduled already. If it has then it will return and not re-schedule the
// function.
func (self *NonblockingOnce) Do(f func()) {
	if atomic.CompareAndSwapInt32(&self.started, notstarted, started) {
		self.wg.Add(1)
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

	self.wg.Wait()
}

func (self *NonblockingOnce) doFunc(f func()) {
	// Schedule cleanup for after this is fully completed.
	defer func(self *NonblockingOnce) {
		self.wg.Done()

		if self.Resettable {
			atomic.StoreInt32(&self.started, notstarted)
		}
	}(self)

	f()
}
