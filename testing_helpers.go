package log

import (
	"io"
	"os"
	"time"
)

// SetExitFunc sets the function used to exit the process for Fatal variants.
// Pass nil to restore default os.Exit.
func SetExitFunc(f func(int)) {
	if f == nil {
		exitFunc = os.Exit
		return
	}
	exitFunc = f
}

// SetNowFunc sets the time source for newly created loggers and updates the
// default logger. Pass nil to restore time.Now.
func SetNowFunc(fn func() time.Time) {
	if fn == nil {
		globalNow = time.Now
	} else {
		globalNow = fn
	}
	// Update default logger's now func as well
	std.mu.Lock()
	std.now = globalNow
	std.mu.Unlock()
}

// SetTestingMode toggles behaviors helpful during tests.
// When on, default logger writes to io.Discard to keep test output clean.
// When off, default logger writes to os.Stderr.
func SetTestingMode(on bool) {
	if on {
		SetOutput(io.Discard)
		return
	}
	SetOutput(os.Stderr)
}
