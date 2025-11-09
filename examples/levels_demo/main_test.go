package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	log "github.com/chrisjoyce911/log"
	"github.com/stretchr/testify/assert"
)

func section(out string, name string) string {
	start := strings.Index(out, "=== Min level: "+name+" ===")
	if start < 0 {
		return ""
	}
	rest := out[start+1:]
	// find next header
	idx := strings.Index(rest, "=== Min level:")
	if idx < 0 {
		return rest
	}
	return rest[:idx]
}

func TestLevelsDemoRun_CoversAllAndFilters(t *testing.T) {
	var buf bytes.Buffer
	// Silence default outputs if any
	log.SetTestingMode(true)
	defer log.SetTestingMode(false)

	Run(&buf, &buf)
	out := buf.String()

	// Basic guards: contains ALL and OFF sections
	assert.Contains(t, out, "=== Min level: ALL ===")
	assert.Contains(t, out, "=== Min level: OFF ===")

	// In ALL section, low and high levels should appear
	all := section(out, "ALL")
	assert.Contains(t, all, "DEBUG [Debug] hello")
	assert.Contains(t, all, "ALERT [Alert] hello")

	// In OFF section, no log lines should appear (only header printed outside this slice)
	off := section(out, "OFF")
	assert.NotContains(t, off, "[Println] hello")
	assert.NotContains(t, off, "[Debug] hello")
	assert.NotContains(t, off, "[Alert] hello")
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() { os.Stdout = old }()
	fn()
	_ = w.Close()
	b, _ := io.ReadAll(r)
	return string(b)
}

func TestMainFunction_CoversMain(t *testing.T) {
	log.SetTestingMode(true)
	defer log.SetTestingMode(false)
	out := captureStdout(t, func() { main() })
	assert.Contains(t, out, "=== Min level:")
}

func TestRunNilLogWriter_UsesStdout(t *testing.T) {
	log.SetTestingMode(true)
	defer log.SetTestingMode(false)
	// Capture stdout because Run will fallback to os.Stdout when log writer is nil
	out := captureStdout(t, func() {
		var buf bytes.Buffer
		Run(&buf, nil)
		// The header goes to buf; ensure it printed
		assert.Contains(t, buf.String(), "=== Min level:")
	})
	// And there should be some logger output to stdout as well
	assert.Contains(t, out, "[Println] hello")
}
