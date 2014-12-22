package stoplock

import "sync"

// Mutex is a stoppable mutual-exclusion lock
type Mutex struct {
	mutex   sync.Mutex
	stopped bool
	onStop  chan struct{}
}

// NewMutex creates a new Mutex object which is not stopped or locked.
func NewMutex() *Mutex {
	return &Mutex{sync.Mutex{}, false, make(chan struct{})}
}

// Lock seizes a mutex.
// If the mutex was stopped, this returns false and the mutex is not seized.
func (m *Mutex) Lock() bool {
	m.mutex.Lock()
	if m.stopped {
		m.mutex.Unlock()
		return false
	}
	return true
}

// OnStop returns a channel that will be closed when the mutex is stopped.
// The caller does not need to seize a lock before calling OnStop() on it.
func (m *Mutex) OnStop() chan struct{} {
	return m.onStop
}

// Stop stops a mutex.
// The receiver must be locked.
// After this returns, the receiver will no longer be locked.
func (m *Mutex) Stop() {
	if m.stopped {
		panic("Stop() called on stopped lock.")
	}
	m.stopped = true
	close(m.onStop)
	m.mutex.Unlock()
}

// Unlock unlocks a seized mutex.
func (m *Mutex) Unlock() {
	m.mutex.Unlock()
}
