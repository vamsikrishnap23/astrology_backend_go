package ephemeris

import "sync"

// Mu ensures thread-safe access to the underlying C Swiss Ephemeris library,
// which is not thread-safe by default when multiple goroutines call it concurrently.
var Mu sync.Mutex
