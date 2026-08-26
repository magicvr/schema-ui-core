package handler

import "sync"

// keyedMutex serializes check-then-act pairs PER KEY instead of globally
// (W13 B-3 · GOAL-013 A-001): a single process-wide mutex made every user's
// upload-quota scan and avatar image processing wait on every other user's.
// Keys are owner ids, so the map is bounded by the account count — no eviction
// needed. The W7/W11 invariants are unchanged: the SAME owner's check+write
// pair still holds one lock end-to-end.
type keyedMutex struct {
	mu    sync.Mutex
	locks map[string]*sync.Mutex
}

func newKeyedMutex() *keyedMutex {
	return &keyedMutex{locks: make(map[string]*sync.Mutex)}
}

// lock acquires the per-key mutex and returns its release func.
func (k *keyedMutex) lock(key string) func() {
	k.mu.Lock()
	l, ok := k.locks[key]
	if !ok {
		l = &sync.Mutex{}
		k.locks[key] = l
	}
	k.mu.Unlock()
	l.Lock()
	return l.Unlock
}
