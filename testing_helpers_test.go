package log

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetExitFuncAndRestore(t *testing.T) {
	var code int
	SetExitFunc(func(c int) { code = c })
	defer SetExitFunc(nil)
	Fatal("trigger exit")
	assert.Equal(t, 1, code)
}

func TestSetNowFuncAffectsNewAndDefault(t *testing.T) {
	fixed := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	SetNowFunc(func() time.Time { return fixed })
	defer SetNowFunc(nil)
	l := New(nil, "", LstdFlags)
	if got := l.now(); !got.Equal(fixed) {
		t.Fatalf("expected fixed time, got %v", got)
	}
	// default logger should also use fixed now
	if got := std.now(); !got.Equal(fixed) {
		t.Fatalf("expected fixed time on std, got %v", got)
	}
}

func TestSetTestingModeToggle(t *testing.T) {
	SetTestingMode(true)
	SetTestingMode(false)
}
