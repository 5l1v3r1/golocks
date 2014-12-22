package golocks

import (
	"testing"
	"time"
)

func TestSkipWait(t *testing.T) {
	lock := NewWaitMutex()
	waitingChannel := make(chan struct{})
	channel := make(chan bool)

	go func() {
		lock.Lock()
		waitingChannel <- struct{}{}
		res := lock.WaitTime(time.Hour)
		if res {
			lock.Unlock()
		}
		channel <- res
	}()

	<-waitingChannel
	lock.Lock()
	lock.Skip()
	lock.Unlock()

	select {
	case val := <-channel:
		if !val {
			t.Fatal("Wait returned false after SkipWait")
		}
	case <-time.After(time.Second):
		t.Fatal("Wait timed out after SkipWait call.")
	}
}

func TestStopWait(t *testing.T) {
	lock := NewWaitMutex()
	waitingChannel := make(chan struct{})
	channel := make(chan bool)

	go func() {
		lock.Lock()
		waitingChannel <- struct{}{}
		res := lock.WaitTime(time.Hour)
		if res {
			lock.Unlock()
		}
		channel <- res
	}()

	<-waitingChannel
	lock.Lock()
	lock.Stop()

	select {
	case val := <-channel:
		if val {
			t.Fatal("Wait returned true after Stop")
		}
	case <-time.After(time.Second):
		t.Fatal("Wait timed out after SkipWait call.")
	}
}
