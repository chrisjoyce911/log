package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

// output wraps a handler with a minimum level threshold.
type output struct {
	h   Handler
	min Level
}

// Logger is a leveled, multi-output logger with a stdlib-like surface.
type Logger struct {
	mu      sync.Mutex
	prefix  string
	flags   int
	outputs []output
	now     func() time.Time
}

// globalNow allows tests to control time used by new loggers.
var globalNow = time.Now

// New creates a new Logger that writes to w for all levels by default,
// keeping stdlib's constructor shape for drop-in adoption.
func New(w io.Writer, prefix string, flag int) *Logger {
	l := &Logger{
		prefix: prefix,
		flags:  flag,
		now:    globalNow,
	}
	if w == nil {
		w = os.Stderr
	}
	l.outputs = []output{{h: &WriterHandler{w: w}, min: LevelDebug}} // default min: debug
	return l
}

// Default returns the package-level standard Logger.
func Default() *Logger { return std }

var std = New(os.Stderr, "", LstdFlags)

// SetFlags sets the output flags on the default logger.
func SetFlags(flag int) { std.SetFlags(flag) }

// SetPrefix sets the output prefix on the default logger.
func SetPrefix(prefix string) { std.SetPrefix(prefix) }

// SetOutput replaces all outputs on the default logger with a single writer.
// This mirrors stdlib log's SetOutput notion while still allowing later additions
// via AddWriter/AddHandler.
func SetOutput(w io.Writer) { std.SetOutput(w) }

// AddWriter attaches a writer for messages at minLevel and above on the default logger.
func AddWriter(minLevel Level, w io.Writer) { std.AddWriter(minLevel, w) }

// AddHandler attaches a custom Handler for messages at minLevel and above on the default logger.
func AddHandler(minLevel Level, h Handler) { std.AddHandler(minLevel, h) }

func (l *Logger) SetFlags(flag int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.flags = flag
}

func (l *Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix
}

func (l *Logger) SetOutput(w io.Writer) {
	if w == nil {
		w = os.Stderr
	}
	l.mu.Lock()
	l.outputs = []output{{h: &WriterHandler{w: w}, min: LevelDebug}}
	l.mu.Unlock()
}

func (l *Logger) AddWriter(minLevel Level, w io.Writer) {
	if w == nil {
		return
	}
	l.AddHandler(minLevel, &WriterHandler{w: w})
}

func (l *Logger) AddHandler(minLevel Level, h Handler) {
	if h == nil {
		return
	}
	l.mu.Lock()
	l.outputs = append(l.outputs, output{h: h, min: minLevel})
	l.mu.Unlock()
}

// With returns a shallow copy of the Logger with additional attributes applied
// to every record. Placeholder: attributes are included in records but only
// minimally formatted by the default WriterHandler.
func (l *Logger) With(kv ...any) *Logger {
	nl := &Logger{
		prefix:  l.prefix,
		flags:   l.flags,
		outputs: l.outputs,
		now:     l.now,
	}
	return nl
}

// Internal helpers used by API and methods
func (l *Logger) logStructured(level Level, msg string, kv ...any) {
	attrs := toAttrs(kv)
	l.dispatch(level, msg, attrs)
}

func (l *Logger) logf(level Level, format string, v ...any) {
	l.dispatch(level, fmt.Sprintf(format, v...), nil)
}

func (l *Logger) dispatch(level Level, msg string, attrs []Attr) {
	l.mu.Lock()
	prefix := l.prefix
	flags := l.flags
	now := l.now
	outs := append([]output(nil), l.outputs...)
	l.mu.Unlock()

	r := Record{
		Time:    now(),
		Level:   level,
		Message: msg,
		Prefix:  prefix,
		Attrs:   attrs,
		Flags:   flags,
	}

	// Optionally capture source info (placeholder: only when flag hints it)
	if flags&(Llongfile|Lshortfile) != 0 {
		if pc, _, _, ok := runtime.Caller(3); ok {
			r.PC = pc
		}
	}

	for _, o := range outs {
		if level >= o.min {
			_ = o.h.Handle(r)
		}
	}
}

func toAttrs(kv []any) []Attr {
	if len(kv) == 0 {
		return nil
	}
	n := len(kv) / 2
	attrs := make([]Attr, 0, n)
	for i := 0; i+1 < len(kv); i += 2 {
		k, _ := kv[i].(string)
		attrs = append(attrs, Attr{Key: k, Value: kv[i+1]})
	}
	return attrs
}

// no-op
