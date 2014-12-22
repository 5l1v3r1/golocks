package golocks

import "time"

type WaitMutex struct {
	*Mutex
	skipWait chan struct{}
}

// NewWaitMutex creates a new, unlocked WaitMutex.
func NewWaitMutex() *WaitMutex {
	return &WaitMutex{NewMutex(), nil}
}

// Skip stops the waiting thread early.
// The receiving mutex must be locked when Skip() is called.
// Returns true if and only if a thread was woken up.
func (m *WaitMutex) Skip() bool {
	if m.skipWait == nil {
		return false
	}
	close(m.skipWait)
	return true
}

// Wait waits for a function to complete.
// The receiving mutex must be locked when Wait() is called.
// The supplied function should not assume that it has ownership of the lock.
// If the mutex is stopped while waiting, false is returned and the caller
// loses ownership of the lock.
// If the function completes successfully or if the wait is skipped, Wait()
// returns true and the lock is re-locked.
// It is invalid to Wait() from multiple Goroutines at once.
func (m *WaitMutex) Wait(f func()) bool {
	if m.skipWait != nil {
		panic("Wait() called on mutex that is already waiting.")
	}
	m.skipWait = make(chan struct{})
	ch := make(chan struct{})
	go func() {
		f()
		close(ch)
	}()
	m.Unlock()
	select {
	case <-ch:
	case <-m.skipWait:
	case <-m.OnStop():
		return false
	}
	if m.Lock() {
		m.skipWait = nil
		return true
	}
	return false
}

// WaitDuration has the same semantics as Wait() except that it waits for a
// given duration instead of calling a function.
func (m *WaitMutex) WaitTime(duration time.Duration) bool {
	return m.Wait(func() {
		time.Sleep(duration)
	})
}
