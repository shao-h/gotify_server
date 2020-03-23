package stream

import (
	"sync"
	"sync/atomic"
)

// once is an object that will perform exactly one action.
type once struct {
	// done indicates whether the action has been performed.
	// It is first in the struct because it is used in the hot path.
	// The hot path is inlined at every call site.
	// Placing done first allows more compact instructions on some architectures (amd64/x86),
	// and fewer instructions (to calculate offset) on other architectures.
	done uint32
	m    sync.Mutex
}

// Do calls the function f if and only if Do is being called for the
// first time for this instance of Once.
func (o *once) Do(f func()) {
	// Do does not guarantee that when it returns, f has finished.
	if atomic.LoadUint32(&o.done) == 1 {
		return
	}

	if o.mayExecute() {
		f()
	}
}

func (o *once) mayExecute() bool {
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		atomic.StoreUint32(&o.done, 1)
		return true
	}
	return false
}
