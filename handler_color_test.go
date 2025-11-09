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

func TestColoredHandler_Defaults_AutoNoTTY_NoColor(t *testing.T) {
	var buf bytes.Buffer
	// NO_COLOR not set; writer is not a TTY (bytes.Buffer), so Auto disables
	os.Unsetenv("NO_COLOR")
	h := NewColoredWriterHandler(&buf, ColorOptions{Mode: ColorAuto})
	r := Record{Level: LevelInfo, Message: "x", Flags: LstdFlags}
	assert.NoError(t, h.Handle(r))
	s := buf.String()
	assert.NotContains(t, s, "\x1b[")
}

func TestColoredHandler_DefaultToggles_LevelColored(t *testing.T) {
	var buf bytes.Buffer
	h := NewColoredWriterHandler(&buf, ColorOptions{Mode: ColorOn}) // no toggles set -> level colored by default
	r := Record{Level: LevelInfo, Message: "hi", Flags: LstdFlags}
	assert.NoError(t, h.Handle(r))
	s := buf.String()
	assert.Contains(t, s, "\x1b[")
}

func TestColoredHandler_NoPaletteMatch_NoColorOnTokens(t *testing.T) {
	var buf bytes.Buffer
	pal := map[Level]string{LevelDebug: "\x1b[35m"} // no entry for INFO
	h := NewColoredWriterHandler(&buf, ColorOptions{Mode: ColorOn, Palette: pal, ColorLevel: true, ColorPrefix: true, ColorMessage: true, ColorAttrs: true})
	r := Record{Level: LevelInfo, Message: "m", Prefix: "p", Flags: LstdFlags, Attrs: []Attr{{Key: "k", Value: "v"}}}
	assert.NoError(t, h.Handle(r))
	s := buf.String()
	// No palette match -> no ANSI sequences should appear
	assert.NotContains(t, s, "\x1b[")
}

func TestColoredHandler_Auto_WithStdoutNonTTY(t *testing.T) {
	// Ensure NO_COLOR is not set so we exercise the isatty branch; on CI this is non-TTY
	_ = os.Unsetenv("NO_COLOR")
	h := NewColoredWriterHandler(os.Stdout, ColorOptions{Mode: ColorAuto})
	// Use a buffer to capture output by swapping os.Stdout is not feasible here; just invoke Handle to execute code paths
	// We don't assert output; this test is for coverage of the isatty branch where writer is *os.File.
	r := Record{Level: LevelInfo, Message: "x", Flags: LstdFlags}
	_ = h.Handle(r)
}

func TestColoredHandler_PrefixMessageAttrs_ColorOn(t *testing.T) {
	var buf bytes.Buffer
	h := NewColoredWriterHandler(&buf, ColorOptions{Mode: ColorOn, ColorLevel: false, ColorPrefix: true, ColorMessage: true, ColorAttrs: true})
	r := Record{Level: LevelError, Message: "msg", Prefix: "p", Flags: LstdFlags, Attrs: []Attr{{Key: "k", Value: 1}}}
	assert.NoError(t, h.Handle(r))
	s := buf.String()
	assert.Contains(t, s, "\x1b[") // has color
	assert.Contains(t, s, "[p]")   // prefix present
	assert.Contains(t, s, "\x1b[31mmsg\x1b[0m")
	assert.Contains(t, s, "\x1b[31mk=1\x1b[0m")
}

func TestColoredHandler_ColorOff_NoAnsi(t *testing.T) {
	var buf bytes.Buffer
	h := NewColoredWriterHandler(&buf, ColorOptions{Mode: ColorOff, ColorLevel: true, ColorPrefix: true, ColorMessage: true, ColorAttrs: true})
	r := Record{Level: LevelInfo, Message: "msg", Prefix: "p", Flags: LstdFlags, Attrs: []Attr{{Key: "k", Value: "v"}}}
	assert.NoError(t, h.Handle(r))
	s := buf.String()
	assert.NotContains(t, s, "\x1b[")
}

func TestColoredHandler_CustomPaletteUsed(t *testing.T) {
	var buf bytes.Buffer
	pal := map[Level]string{LevelInfo: "\x1b[34m"}
	h := NewColoredWriterHandler(&buf, ColorOptions{Mode: ColorOn, Palette: pal, ColorLevel: true})
	r := Record{Level: LevelInfo, Message: "info", Flags: LstdFlags}
	assert.NoError(t, h.Handle(r))
	s := buf.String()
	assert.Contains(t, s, "\x1b[34m")
}
