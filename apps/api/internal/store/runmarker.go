package store

import (
	"errors"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

// Nested Run detection (R1 v1.4 A-008 F-004).
//
// The kernel.Tx callback signature is func(Tx) error and does not receive a
// context, so a ctx-value marker cannot reliably flow into a nested Run. The
// contract allows "调用栈 / goroutine 局部" as the per-callback mechanism; we
// use a goroutine-local counter guarded by a mutex. This detects a nested Run
// made from within the *same goroutine / callback stack* while leaving the
// concurrent Runs that the postgres pool is allowed to carry untouched (a
// Store-level or process-level flag would wrongly forbid those).
var (
	runMu        sync.Mutex
	active       = map[uint64]int{}
	errNestedRun = errors.New("store: nested Run is forbidden (R1 v1.4 §2)")
)

// goroutineID returns the current goroutine id via runtime.Stack (the first
// line has the form "goroutine N [..]"). Used only at Run entry/exit, not per
// query.
func goroutineID() uint64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	line := string(buf[:n])
	if !strings.HasPrefix(line, "goroutine ") {
		return 0
	}
	line = line[len("goroutine "):]
	if i := strings.IndexByte(line, ' '); i >= 0 {
		line = line[:i]
	}
	id, _ := strconv.ParseUint(line, 10, 64)
	return id
}

// enterRun marks the current goroutine as inside one Run callback. It fails
// closed when this goroutine is already inside a Run (nested call).
func enterRun() error {
	id := goroutineID()
	runMu.Lock()
	defer runMu.Unlock()
	if active[id] > 0 {
		return errNestedRun
	}
	active[id]++
	return nil
}

// leaveRun unmarks the current goroutine; it must pair with a successful
// enterRun.
func leaveRun() {
	id := goroutineID()
	runMu.Lock()
	defer runMu.Unlock()
	if active[id] > 1 {
		active[id]--
		return
	}
	delete(active, id)
}
