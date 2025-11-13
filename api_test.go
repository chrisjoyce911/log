package log

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func withStdReset(t *testing.T, fn func()) {
	t.Helper()
	oldFlags := std.flags
	oldOutputs := std.outputs
	oldPrefix := std.prefix
	defer func() {
		std.flags = oldFlags
		std.outputs = oldOutputs
		std.prefix = oldPrefix
	}()
	fn()
}

func TestPrintAndLevelsToWriter(t *testing.T) {
	withStdReset(t, func() {
		var buf bytes.Buffer
		SetOutput(io.Discard)
		AddWriter(LevelAll, &buf)
		SetFlags(Ldate | Ltime)
		SetPrefix("pfx")

		Print("pmsg")
		Printf("%s", "pmsg2")
		Println("pmsg3")
		Debug("dmsg")
		Debugf("df: %s", "d")
		Trace("tmsg")
		Tracef("tf: %s", "t")
		Verbose("vmsg")
		Verbosef("vf: %s", "v")
		Detail("demsg")
		Detailf("def: %s", "de")
		Info("imsg")
		Infof("if: %s", "i")
		Notice("nmsg")
		Noticef("nf: %s", "n")
		Warn("wmsg")
		Warnf("wf: %s", "w")
		Error("emsg")
		Errorf("ef: %s", "e")
		Critical("cmsg")
		Criticalf("cf: %s", "c")
		Alert("amsg")
		Alertf("af: %s", "a")

		out := buf.String()
		assert.Contains(t, out, "INFO [pfx] pmsg")
		assert.Contains(t, out, "INFO [pfx] pmsg2")
		assert.Contains(t, out, "INFO [pfx] pmsg3")
		assert.Contains(t, out, "DEBUG [pfx] dmsg")
		assert.Contains(t, out, "TRACE [pfx] tmsg")
		assert.Contains(t, out, "VERBOSE [pfx] vmsg")
		assert.Contains(t, out, "DETAIL [pfx] demsg")
		assert.Contains(t, out, "INFO [pfx] imsg")
		assert.Contains(t, out, "NOTICE [pfx] nmsg")
		assert.Contains(t, out, "WARN [pfx] wmsg")
		assert.Contains(t, out, "ERROR [pfx] emsg")
		assert.Contains(t, out, "CRITICAL [pfx] cmsg")
		assert.Contains(t, out, "ALERT [pfx] amsg")
	})
}

func TestLoggerInstanceAndWith(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf, "myp:", Ldate|Ltime)
	l.With("ignored", 1) // placeholder coverage
	l.Println("hello")
	l.Print("p")
	l.Printf("%s", "pf")
	l.Debug("d")
	l.Debugf("df: %s", "d")
	l.Trace("t")
	l.Tracef("tf: %s", "t")
	l.Verbose("v")
	l.Verbosef("vf: %s", "v")
	l.Detail("de")
	l.Detailf("def: %s", "de")
	l.Info("i")
	l.Infof("if: %s", "i")
	l.Notice("n")
	l.Noticef("nf: %s", "n")
	l.Warn("w")
	l.Warnf("wf: %s", "w")
	l.Error("e")
	l.Errorf("ef: %s", "e")
	l.Critical("c")
	l.Criticalf("cf: %s", "c")
	l.Alert("a")
	l.Alertf("af: %s", "a")
	assert.Contains(t, buf.String(), "INFO [myp:] hello")
}

func TestAddHandlerAndJSONOutput(t *testing.T) {
	withStdReset(t, func() {
		var buf bytes.Buffer
		SetOutput(io.Discard)
		jh := NewJSONHandler(&buf)
		AddHandler(LevelInfo, jh)
		SetFlags(Ldate | Ltime)
		Info("json test", "key", "val", "age", 42)

		lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
		assert.Len(t, lines, 1)
		var m map[string]any
		assert.NoError(t, json.Unmarshal([]byte(lines[0]), &m))
		assert.Equal(t, "INFO", m["level"])
		assert.Equal(t, "json test", m["msg"])
		attrs, _ := m["attrs"].(map[string]any)
		assert.Equal(t, "val", attrs["key"])
		assert.Equal(t, float64(42), attrs["age"])
	})
}

func TestJSONHandlerSourceIncludesWhenShortfile(t *testing.T) {
	withStdReset(t, func() {
		var buf bytes.Buffer
		SetOutput(io.Discard)
		AddHandler(LevelInfo, NewJSONHandler(&buf))
		SetFlags(LstdFlags | Lshortfile)
		Info("src test")
		var m map[string]any
		_ = json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &m)
		_, has := m["source"]
		assert.True(t, has)
	})
}

func TestStringChanHandler(t *testing.T) {
	ch := make(chan string, 2)
	h := &StringChanHandler{C: ch}
	r := Record{Level: LevelInfo, Message: "m", Flags: LstdFlags}
	assert.NoError(t, h.Handle(r))
	str := <-ch
	assert.Contains(t, str, "INFO m")

	r2 := Record{Level: LevelInfo, Message: "m2", Flags: 0}
	assert.NoError(t, h.Handle(r2))
	str2 := <-ch
	assert.Equal(t, "INFO m2", str2)
}

func TestDefaultAndHelpers(t *testing.T) {
	assert.NotNil(t, Default())
	assert.Equal(t, "x-2", formatMsg("%s-%d", "x", 2))
	assert.Equal(t, "no-nl", trimNL("no-nl"))
}

func TestWriterHandlerFormatting(t *testing.T) {
	var buf bytes.Buffer
	h := &WriterHandler{w: &buf}
	r := Record{Level: LevelWarn, Message: "hi", Prefix: "p", Flags: LstdFlags, Attrs: []Attr{{Key: "k", Value: 1}}}
	assert.NoError(t, h.Handle(r))
	s := buf.String()
	assert.Contains(t, s, "WARN [p] hi k=1")
}

func TestFatalInterceptsExit(t *testing.T) {
	withStdReset(t, func() {
		var code int
		orig := exitFunc
		exitFunc = func(c int) { code = c }
		defer func() { exitFunc = orig }()

		SetOutput(io.Discard)
		AddWriter(LevelAll, io.Discard)
		Fatal("fatal msg")
		Fatalln("fatal ln")
		Fatalf("%s", "fatal f")

		assert.Equal(t, 1, code)
	})
}

func TestPanicVariants(t *testing.T) {
	withStdReset(t, func() {
		SetOutput(io.Discard)
		AddWriter(LevelAll, io.Discard)
		assert.Panics(t, func() { Panic("pmsg") })
		assert.Panics(t, func() { Panicln("pmsg") })
		assert.Panics(t, func() { Panicf("%s", "pmsg") })
	})
}
