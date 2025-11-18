package log

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestHTTPLogging_Info_Warn_Error_and_Query(t *testing.T) {
	withStdReset(t, func() {
		var buf bytes.Buffer
		SetOutput(&buf)
		SetFlags(Ldate | Ltime)

		// 200 OK
		ok := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})
		h := HTTPLogging(ok, &HTTPLogOptions{Mode: ColorOn, IncludeQuery: true})
		req := httptest.NewRequest(http.MethodGet, "/a?q=1", nil)
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)
		out := buf.String()
		if !strings.Contains(out, "INFO") || !strings.Contains(out, "?q=1") {
			t.Fatalf("expected INFO access with query, got: %s", out)
		}

		// 404 -> WARN
		buf.Reset()
		notFound := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		})
		h2 := HTTPLogging(notFound, &HTTPLogOptions{Mode: ColorOn, IncludeQuery: true})
		rr2 := httptest.NewRecorder()
		h2.ServeHTTP(rr2, httptest.NewRequest(http.MethodGet, "/nf", nil))
		out2 := buf.String()
		if !strings.Contains(out2, "WARN") {
			t.Fatalf("expected WARN access, got: %s", out2)
		}

		// 500 -> ERROR
		buf.Reset()
		boom := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		h3 := HTTPLogging(boom, &HTTPLogOptions{Mode: ColorOn})
		rr3 := httptest.NewRecorder()
		h3.ServeHTTP(rr3, httptest.NewRequest(http.MethodGet, "/e", nil))
		out3 := buf.String()
		if !strings.Contains(out3, "ERROR") {
			t.Fatalf("expected ERROR access, got: %s", out3)
		}
	})
}

func TestHTTPLogging_BodyLoggingAndTruncate(t *testing.T) {
	withStdReset(t, func() {
		var buf bytes.Buffer
		SetOutput(&buf)
		SetFlags(Ldate | Ltime)

		echo := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("ok"))
		})
		// Small cap to force truncation
		h := HTTPLogging(echo, &HTTPLogOptions{Mode: ColorOn, LogPostBody: true, MaxBodyBytes: 5, IncludeQuery: true})

		req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader("0123456789"))
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)
		out := buf.String()
		if !strings.Contains(out, "body=01234") || !strings.Contains(out, "(truncated)") {
			t.Fatalf("expected body preview with truncation, got: %s", out)
		}
	})
}

func TestHTTPLogging_BodyLogging_NoTruncate(t *testing.T) {
	withStdReset(t, func() {
		var buf bytes.Buffer
		SetOutput(&buf)
		SetFlags(Ldate | Ltime)

		h := HTTPLogging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}), &HTTPLogOptions{Mode: ColorOn, LogPostBody: true, MaxBodyBytes: 1024})

		req := httptest.NewRequest(http.MethodPut, "/echo", strings.NewReader("abc"))
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)
		out := buf.String()
		if !strings.Contains(out, "body=abc") || strings.Contains(out, "truncated") {
			t.Fatalf("expected full body without truncation, got: %s", out)
		}
	})
}

func TestHTTPLogging_ColorModes(t *testing.T) {
	withStdReset(t, func() {
		var buf bytes.Buffer
		SetOutput(&buf)
		SetFlags(Ldate | Ltime)

		h := HTTPLogging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}), &HTTPLogOptions{Mode: ColorOn})
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest(http.MethodHead, "/head", nil))
		if !strings.Contains(buf.String(), "\x1b[") { // colored method/path
			t.Fatalf("expected ANSI color in output")
		}

		// Auto mode disabled by NO_COLOR
		buf.Reset()
		_ = os.Setenv("NO_COLOR", "1")
		defer os.Unsetenv("NO_COLOR")
		h2 := HTTPLogging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}), &HTTPLogOptions{Mode: ColorAuto})
		rr2 := httptest.NewRecorder()
		h2.ServeHTTP(rr2, httptest.NewRequest(http.MethodGet, "/auto", nil))
		if strings.Contains(buf.String(), "\x1b[") {
			t.Fatalf("did not expect ANSI color in Auto mode with NO_COLOR set")
		}
	})
}

func TestHTTPLogging_DefaultIncludeQueryFalseWhenNilOpts(t *testing.T) {
	withStdReset(t, func() {
		var buf bytes.Buffer
		SetOutput(&buf)
		SetFlags(Ldate | Ltime)

		h := HTTPLogging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}), nil)
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/p?q=1", nil))
		// Expect no query present in message when opts=nil (current default)
		if strings.Contains(buf.String(), "?q=1") {
			t.Fatalf("did not expect query in path by default")
		}
	})
}

func TestHTTPLogging_ColorOff_NoAnsi(t *testing.T) {
	withStdReset(t, func() {
		var buf bytes.Buffer
		SetOutput(&buf)
		SetFlags(Ldate | Ltime)
		h := HTTPLogging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}), &HTTPLogOptions{Mode: ColorOff, IncludeQuery: true})
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/x?y=1", nil))
		if strings.Contains(buf.String(), "\x1b[") {
			t.Fatalf("did not expect ANSI color when ColorOff")
		}
	})
}

func TestHTTPLogging_DefaultMode_ColorsOn(t *testing.T) {
	withStdReset(t, func() {
		var buf bytes.Buffer
		SetOutput(&buf)
		SetFlags(Ldate | Ltime)
		// Default (zero value) should be ColorOn
		h := HTTPLogging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}), &HTTPLogOptions{})
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/default", nil))
		if !strings.Contains(buf.String(), "\x1b[") {
			t.Fatalf("expected ANSI color with default mode (ColorOn)")
		}
	})
}

func TestHTTPLogging_MethodColorBranches(t *testing.T) {
	withStdReset(t, func() {
		var buf bytes.Buffer
		SetOutput(&buf)
		SetFlags(Ldate | Ltime)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		h := HTTPLogging(handler, &HTTPLogOptions{Mode: ColorOn})
		// PUT
		h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPut, "/p", nil))
		// DELETE
		h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodDelete, "/d", nil))
		// PATCH
		h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPatch, "/pt", nil))
		// HEAD (default branch)
		h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodHead, "/h", nil))
		_ = buf.String()
	})
}

func TestHTTPLogging_StatusDefaultOKWhenNoWrite(t *testing.T) {
	withStdReset(t, func() {
		var buf bytes.Buffer
		SetOutput(&buf)
		SetFlags(Ldate | Ltime)
		// Handler does not write headers/body -> status defaults to 200
		h := HTTPLogging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}), &HTTPLogOptions{Mode: ColorOn})
		h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/nowrite", nil))
		if !strings.Contains(buf.String(), "status=200") {
			t.Fatalf("expected default status=200 when handler writes nothing")
		}
	})
}
