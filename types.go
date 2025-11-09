package log

import "time"

// Attr is a simple key/value attribute similar to slog.Attr.
type Attr struct {
	Key   string
	Value any
}

// Record is a lightweight log record passed to Handlers.
type Record struct {
	Time    time.Time
	Level   Level
	Message string
	Prefix  string
	Attrs   []Attr
	PC      uintptr
	Flags   int
}

// Handler consumes log Records. Custom handlers may implement JSON output,
// remote shipping, buffering, etc.
type Handler interface {
	Handle(r Record) error
}
