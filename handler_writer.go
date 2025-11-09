package log

import (
	"fmt"
	"io"
	"strings"
)

// WriterHandler is a basic text writer handler.
type WriterHandler struct {
	w io.Writer
}

func (h *WriterHandler) Handle(r Record) error {
	// Minimal text line: timestamp level [prefix] message key=val ...
	b := &strings.Builder{}
	if ts := formatTimestamp(r.Time, r.Flags); ts != "" {
		b.WriteString(ts)
		b.WriteByte(' ')
	}
	b.WriteString(r.Level.String())
	if r.Prefix != "" {
		b.WriteByte(' ')
		b.WriteByte('[')
		b.WriteString(r.Prefix)
		b.WriteByte(']')
	}
	if r.Message != "" {
		b.WriteByte(' ')
		b.WriteString(r.Message)
	}
	for _, a := range r.Attrs {
		b.WriteByte(' ')
		b.WriteString(a.Key)
		b.WriteByte('=')
		b.WriteString(fmt.Sprint(a.Value))
	}
	b.WriteByte('\n')
	_, err := io.WriteString(h.w, b.String())
	return err
}
