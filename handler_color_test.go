package log

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColoredHandler_LevelOnly_ColorOn(t *testing.T) {
	var buf bytes.Buffer
	h := NewColoredWriterHandler(&buf, ColorOptions{Mode: ColorOn, ColorLevel: true})
	r := Record{Level: LevelWarn, Message: "hello", Flags: LstdFlags}
	assert.NoError(t, h.Handle(r))
	s := buf.String()
	// Level token colored, message plain
	assert.Contains(t, s, "\x1b[")
	assert.Contains(t, s, "WARN")
	assert.Contains(t, s, " hello")
}

func TestColoredHandler_Auto_DisabledWithNO_COLOR(t *testing.T) {
	var buf bytes.Buffer
	os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")
	h := NewColoredWriterHandler(&buf, ColorOptions{Mode: ColorAuto, ColorLevel: true})
	r := Record{Level: LevelInfo, Message: "m", Flags: LstdFlags}
	assert.NoError(t, h.Handle(r))
	s := buf.String()
	assert.NotContains(t, s, "\x1b[")
}

func TestSetColoredOutputHelper(t *testing.T) {
	withStdReset(t, func() {
		SetColoredOutput(LevelInfo, ColorOptions{Mode: ColorOn, ColorLevel: true})
		Info("ok")
	})
}
