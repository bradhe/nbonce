package nbonce

import (
	"time"
	"testing"
)

func TestOnlySchedulesFunctionsOnce(t *testing.T) {
	var ran int

	f := func() { ran += 1 }

	once := &NonblockingOnce{}

	for i := 0; i < 10000; i++ {
		once.Do(f)
	}

	once.Wait()

	if ran != 1 {
		t.Fatalf("expected ran to be 1, got %d", ran)
	}
}

func TestWaitReturnsIfOnceWasNeverCalled(t *testing.T) {
	once := &NonblockingOnce{}

	timeout := make(chan bool, 1)

	go func() {
		once.Wait()
		timeout <- true
	}()

	var worked bool

	select {
	case res := <-timeout:
		worked = res
	case <-time.After(time.Second * 1):
		worked = false
	}

	if !worked {
		t.Fatalf("expected worked to be true, got false")
	}
}

func TestSchedulesMultipleWhenResettable(t *testing.T) {
	var ran int

	f := func() { ran += 1 }

	once := &NonblockingOnce{}
	once.Resettable = true

	once.Do(f)
	once.Wait()

	once.Do(f)
	once.Wait()

	if ran != 2 {
		t.Fatalf("expected ran to be 2, got %d", ran)
	}
}

func TestDoesNotRescheduleWhenNotResettable(t *testing.T) {
	var ran int

	f := func() { ran += 1 }

	once := &NonblockingOnce{}
	once.Resettable = false

	once.Do(f)
	once.Wait()

	once.Do(f)
	once.Wait()

	if ran != 1 {
		t.Fatalf("expected ran to be 1, got %d", ran)
	}
}
