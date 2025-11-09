package log

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
)

// JSONHandler writes each record as a single JSON object per line.
// Fields: time, level, msg, optional prefix, attrs map, optional source.
type JSONHandler struct {
	mu  sync.Mutex
	enc *json.Encoder
}

// NewJSONHandler creates a JSONHandler writing to w (defaults to stderr if nil).
func NewJSONHandler(w io.Writer) *JSONHandler {
	if w == nil {
		w = os.Stderr
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return &JSONHandler{enc: enc}
}

func (h *JSONHandler) Handle(r Record) error {
	m := map[string]any{
		"time":  formatTimestamp(r.Time, r.Flags),
		"level": r.Level.String(),
		"msg":   r.Message,
	}
	if r.Prefix != "" {
		m["prefix"] = r.Prefix
	}
	if len(r.Attrs) > 0 {
		attrs := make(map[string]any, len(r.Attrs))
		for _, a := range r.Attrs {
			attrs[a.Key] = a.Value
		}
		m["attrs"] = attrs
	}
	if r.PC != 0 {
		if fn := runtime.FuncForPC(r.PC); fn != nil {
			if file, line := fn.FileLine(r.PC); file != "" {
				m["source"] = fmt.Sprintf("%s:%d", file, line)
			}
		}
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.enc.Encode(m)
}
