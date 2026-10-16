package ephemeris

import (
	"runtime"
	"sync"

	"github.com/tejzpr/go-swisseph"
)

var mu sync.Mutex

// Lock must be held around every call into the Swiss Ephemeris, and released with Unlock.
//
// The C library is not thread-safe, so calls are serialised with a mutex. It also keeps its state (the
// ephemeris path, the sidereal mode) per OS thread, and Go may move a goroutine to a different thread
// between two cgo calls. Lock therefore pins the calling goroutine to its OS thread until Unlock, and
// re-applies the ephemeris path on that thread, so every call inside the section sees the same state.
// Callers still set the sidereal mode (swisseph.SetSidMode) themselves before using the ayanamsa.
//
// Sections must not be nested: Lock is not reentrant.
func Lock() {
	mu.Lock()
	runtime.LockOSThread()
	if EphePath != "" {
		swisseph.SetEphePath(EphePath)
	}
}

// Unlock releases the lock taken by Lock and unpins the goroutine.
func Unlock() {
	runtime.UnlockOSThread()
	mu.Unlock()
}
