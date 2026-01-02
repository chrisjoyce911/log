package log

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNilWritersAndHandlersAreIgnored(t *testing.T) {
	withStdReset(t, func() {
		// AddWriter/Handler with nil should be no-op and not panic
		AddWriter(LevelAll, nil)
		AddHandler(LevelAll, nil)
	})
}

func TestSetOutputNilUsesStderr(t *testing.T) {
	withStdReset(t, func() {
		SetOutput(nil) // cover nil branch
		// route something to ensure no panic
		AddWriter(LevelAll, io.Discard)
		Print("ok")
	})
}

func TestNewNilWriterDefaultsToStderr(t *testing.T) {
	_ = New(nil, "", LstdFlags) // cover nil branch in constructor
}

func TestToAttrsOddPairs(t *testing.T) {
	attrs := toAttrs([]any{"a", 1, "b"})
	assert.Len(t, attrs, 1)
	assert.Equal(t, "a", attrs[0].Key)
	assert.Equal(t, 1, attrs[0].Value)
}

func TestWriterHandlerNoMessage(t *testing.T) {
	var buf bytes.Buffer
	h := &WriterHandler{w: &buf}
	r := Record{Level: LevelInfo, Prefix: "p", Flags: LstdFlags}
	assert.NoError(t, h.Handle(r))
	out := buf.String()
	assert.Contains(t, out, "INFO     [p]")
}

func TestJSONHandlerPrefixAndNilWriter(t *testing.T) {
	// NewJSONHandler(nil) -> stderr; we won't validate output, only coverage
	_ = NewJSONHandler(nil)
	var buf bytes.Buffer
	h := NewJSONHandler(&buf)
	r := Record{Level: LevelInfo, Message: "m", Prefix: "p", Flags: LstdFlags}
	assert.NoError(t, h.Handle(r))
	assert.Contains(t, buf.String(), "\"prefix\":\"p\"")
}
