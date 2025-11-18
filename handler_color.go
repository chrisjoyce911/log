package log

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
)

const ansiReset = "\x1b[0m"

// Default level -> ANSI color palette.
var defaultPalette = map[Level]string{
	LevelTrace:    "\x1b[2m",     // dim
	LevelVerbose:  "\x1b[37m",    // gray
	LevelDebug:    "\x1b[36m",    // cyan
	LevelDetail:   "\x1b[32m",    // green
	LevelInfo:     "\x1b[32m",    // green
	LevelNotice:   "\x1b[36m",    // cyan
	LevelWarn:     "\x1b[33m",    // yellow
	LevelError:    "\x1b[31m",    // red
	LevelCritical: "\x1b[35m",    // magenta
	LevelAlert:    "\x1b[91m",    // bright red
	LevelFatal:    "\x1b[97;41m", // white on red
	LevelPanic:    "\x1b[97;41m", // white on red
}

// ColorMode controls whether coloring is enabled.
type ColorMode int

const (
	ColorOn   ColorMode = iota // default: always enable colors
	ColorOff                   // disable colors
	ColorAuto                  // enable if TTY and NO_COLOR not set
)

// ColorOptions configures the colored console handler.
// By default only the level token is colored.
type ColorOptions struct {
	Mode         ColorMode
	Palette      map[Level]string
	ColorLevel   bool
	ColorPrefix  bool
	ColorMessage bool
	ColorAttrs   bool // color key=value as a whole
}

// ColoredWriterHandler writes text like WriterHandler but with ANSI coloring.
type ColoredWriterHandler struct {
	w       io.Writer
	opts    ColorOptions
	enabled bool
}

// NewColoredWriterHandler constructs a colored handler writing to w.
// By default (Mode=ColorOn) colors are always enabled. Use ColorOff to disable, or ColorAuto for TTY detection.
func NewColoredWriterHandler(w io.Writer, opts ColorOptions) *ColoredWriterHandler {
	if opts.Palette == nil {
		opts.Palette = defaultPalette
	}
	// defaults: only level colored
	if !opts.ColorLevel && !opts.ColorPrefix && !opts.ColorMessage && !opts.ColorAttrs {
		opts.ColorLevel = true
	}
	enabled := true // default: ColorOn
	switch opts.Mode {
	case ColorOff:
		enabled = false
	case ColorAuto:
		// Check NO_COLOR and TTY
		enabled = false // default to false for auto mode
		if os.Getenv("NO_COLOR") == "" {
			// Best effort: detect if writer is a file with FD and is a terminal
			if f, ok := w.(*os.File); ok {
				fd := f.Fd()
				enabled = isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)
			}
		}
	}
	return &ColoredWriterHandler{w: w, opts: opts, enabled: enabled}
}

func (h *ColoredWriterHandler) Handle(r Record) error {
	b := &strings.Builder{}
	if ts := formatTimestamp(r.Time, r.Flags); ts != "" {
		b.WriteString(ts)
		b.WriteByte(' ')
	}

	// Level token
	lvl := r.Level.String()
	if h.enabled && h.opts.ColorLevel {
		if c, ok := h.opts.Palette[r.Level]; ok {
			lvl = c + lvl + ansiReset
		}
	}
	b.WriteString(lvl)

	// Prefix
	if r.Prefix != "" {
		b.WriteByte(' ')
		if h.enabled && h.opts.ColorPrefix {
			if c, ok := h.opts.Palette[r.Level]; ok {
				b.WriteString(c)
				b.WriteByte('[')
				b.WriteString(r.Prefix)
				b.WriteByte(']')
				b.WriteString(ansiReset)
			} else {
				b.WriteByte('[')
				b.WriteString(r.Prefix)
				b.WriteByte(']')
			}
		} else {
			b.WriteByte('[')
			b.WriteString(r.Prefix)
			b.WriteByte(']')
		}
	}

	// Message
	if r.Message != "" {
		b.WriteByte(' ')
		if h.enabled && h.opts.ColorMessage {
			if c, ok := h.opts.Palette[r.Level]; ok {
				b.WriteString(c)
				b.WriteString(r.Message)
				b.WriteString(ansiReset)
			} else {
				b.WriteString(r.Message)
			}
		} else {
			b.WriteString(r.Message)
		}
	}

	// Attrs
	for _, a := range r.Attrs {
		b.WriteByte(' ')
		if h.enabled && h.opts.ColorAttrs {
			if c, ok := h.opts.Palette[r.Level]; ok {
				b.WriteString(c)
				b.WriteString(a.Key)
				b.WriteByte('=')
				b.WriteString(fmt.Sprint(a.Value))
				b.WriteString(ansiReset)
				continue
			}
		}
		b.WriteString(a.Key)
		b.WriteByte('=')
		b.WriteString(fmt.Sprint(a.Value))
	}

	b.WriteByte('\n')
	_, err := io.WriteString(h.w, b.String())
	return err
}

// SetColoredOutput is a convenience to send colored logs to stdout.
// It replaces the default logger's output with a colored handler and discards plain text.
func SetColoredOutput(minLevel Level, opts ColorOptions) {
	// Silence default text output
	SetOutput(io.Discard)
	// Attach colored handler to stdout
	AddHandler(minLevel, NewColoredWriterHandler(os.Stdout, opts))
}
