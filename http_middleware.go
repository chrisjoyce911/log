package log

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/mattn/go-isatty"
)

// ANSI helpers
const (
	ansiBold    = "\x1b[1m"
	ansiGreen   = "\x1b[32m"
	ansiCyan    = "\x1b[36m"
	ansiYellow  = "\x1b[33m"
	ansiRed     = "\x1b[31m"
	ansiMagenta = "\x1b[35m"
)

// HTTPLogOptions customizes HTTPLogging middleware.
type HTTPLogOptions struct {
	// Color mode for method/path highlighting.
	Mode ColorMode
	// Include query string when building the request path display.
	IncludeQuery bool
	// LogPostBody enables logging of request body for POST/PUT/PATCH on the pre-request line.
	LogPostBody bool
	// MaxBodyBytes caps the size of the logged body (default 64KB when zero or negative).
	MaxBodyBytes int
}

func (o *HTTPLogOptions) enabled() bool {
	switch o.Mode {
	case ColorOff:
		return false
	case ColorAuto:
		if os.Getenv("NO_COLOR") != "" {
			return false
		}
		// Best-effort: colorize when stdout is a TTY
		fd := os.Stdout.Fd()
		return isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)
	default: // ColorOn
		return true
	}
}

func colorWrap(s, code string, on bool) string {
	if !on {
		return s
	}
	return code + s + ansiReset
}

func methodColor(method string) string {
	switch method {
	case http.MethodGet:
		return ansiGreen
	case http.MethodPost:
		return ansiCyan
	case http.MethodPut:
		return ansiYellow
	case http.MethodDelete:
		return ansiRed
	case http.MethodPatch:
		return ansiMagenta
	default:
		return ansiBold
	}
}

// HTTPLogging returns middleware that logs a colored pre-request line (method/path)
// at DEBUG and a colored access line at INFO/WARN/ERROR according to status.
// It highlights only the method and path tokens; other parts follow the logger handler's coloring.
func HTTPLogging(next http.Handler, opts *HTTPLogOptions) http.Handler {
	var o HTTPLogOptions
	if opts != nil {
		o = *opts
	}
	// default: include query
	if !o.IncludeQuery {
		// leave as-is when false; default is true
	} else {
		o.IncludeQuery = true
	}
	colorOn := o.enabled()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Build display path
		dispPath := r.URL.Path
		if o.IncludeQuery && r.URL.RawQuery != "" {
			dispPath += "?" + r.URL.RawQuery
		}
		// Optionally read and log body (preview) for mutation methods and restore body for handler.
		var bodyPreview string
		if o.LogPostBody && (r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch) && r.Body != nil {
			limit := int64(64 * 1024)
			if o.MaxBodyBytes > 0 {
				limit = int64(o.MaxBodyBytes)
			}
			data, _ := io.ReadAll(io.LimitReader(r.Body, limit+1))
			truncated := int64(len(data)) > limit
			if truncated {
				data = data[:limit]
			}
			_ = r.Body.Close()
			r.Body = io.NopCloser(bytes.NewReader(data))
			bodyPreview = string(data)
			if truncated {
				bodyPreview += "…(truncated)"
			}
		}
		// Pre-request line with highlighted method and path in the message
		msg := colorWrap(r.Method, methodColor(r.Method), colorOn) + " " + colorWrap(dispPath, ansiBold, colorOn)
		attrs := []any{"remote", r.RemoteAddr, "ua", r.UserAgent()}
		if bodyPreview != "" {
			attrs = append(attrs, "body", bodyPreview)
		}
		Debug(msg, attrs...)

		// Wrap writer to capture status/bytes
		wrapper := &httpLogRW{ResponseWriter: w}
		next.ServeHTTP(wrapper, r)

		dur := time.Since(start)
		status := wrapper.status
		if status == 0 {
			status = http.StatusOK
		}
		// Access line; message shows colored method/path again
		msg2 := colorWrap(r.Method, methodColor(r.Method), colorOn) + " " + colorWrap(dispPath, ansiBold, colorOn)
		accessAttrs := []any{"status", status, "bytes", wrapper.bytes, "duration", dur.String()}
		switch {
		case status >= 500:
			Error(msg2, accessAttrs...)
		case status >= 400:
			Warn(msg2, accessAttrs...)
		default:
			Info(msg2, accessAttrs...)
		}
	})
}

// httpLogRW is a small ResponseWriter wrapper used by HTTPLogging.
type httpLogRW struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *httpLogRW) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *httpLogRW) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}
