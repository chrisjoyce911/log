package log

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetOutputFileAndAddFileWriterAndJSON(t *testing.T) {
	withStdReset(t, func() {
		tdir := t.TempDir()
		p1 := filepath.Join(tdir, "a", "b", "text.log")
		p2 := filepath.Join(tdir, "x", "y", "writer.log")
		p3 := filepath.Join(tdir, "j", "k", "json.log")

		// SetOutputFile
		f1, err := SetOutputFile(p1)
		assert.NoError(t, err)
		defer f1.Close()
		SetFlags(Ldate | Ltime)
		Println("hello text")
		b1, err := os.ReadFile(p1)
		assert.NoError(t, err)
		assert.Contains(t, string(b1), "hello text")

		// AddFileWriter
		SetOutput(io.Discard)
		f2, err := AddFileWriter(LevelInfo, p2)
		assert.NoError(t, err)
		defer f2.Close()
		Info("hello writer")
		b2, err := os.ReadFile(p2)
		assert.NoError(t, err)
		assert.Contains(t, string(b2), "hello writer")

		// AddJSONFile
		SetOutput(io.Discard)
		f3, err := AddJSONFile(LevelInfo, p3)
		assert.NoError(t, err)
		defer f3.Close()
		Info("hello json", "k", 1)
		ff, err := os.Open(p3)
		assert.NoError(t, err)
		defer ff.Close()
		s := bufio.NewScanner(ff)
		assert.True(t, s.Scan())
		var m map[string]any
		assert.NoError(t, json.Unmarshal([]byte(s.Text()), &m))
		assert.Equal(t, "INFO", m["level"])
		attrs, _ := m["attrs"].(map[string]any)
		assert.Equal(t, float64(1), attrs["k"])
		assert.False(t, strings.Contains(string(b1), "\n\n"))
	})
}

func TestOpenFileAppendAndRootPath(t *testing.T) {
	// Test append behavior
	tdir := t.TempDir()
	ap := filepath.Join(tdir, "append", "file.log")
	f, err := OpenFileAppend(ap, 0o644)
	assert.NoError(t, err)
	_, _ = f.WriteString("a")
	_ = f.Close()
	f2, err := OpenFileAppend(ap, 0o644)
	assert.NoError(t, err)
	_, _ = f2.WriteString("b")
	_ = f2.Close()
	b, err := os.ReadFile(ap)
	assert.NoError(t, err)
	assert.Equal(t, "ab", string(b))

	// Test ensureDirForFile branch for current directory (no parent dir)
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	_ = os.Chdir(tdir)
	f3, err := OpenFileTruncate("rootfile.log", 0o644)
	assert.NoError(t, err)
	_ = f3.Close()
	_, err = os.Stat(filepath.Join(tdir, "rootfile.log"))
	assert.NoError(t, err)
}

func TestFileHelpersErrorPaths(t *testing.T) {
	tdir := t.TempDir()
	// Create a file that will be used as a parent path to trigger MkdirAll error
	parentFile := filepath.Join(tdir, "parent")
	_ = os.WriteFile(parentFile, []byte("x"), 0o644)

	// OpenFileTruncate error
	_, err := OpenFileTruncate(filepath.Join(parentFile, "child.log"), 0o644)
	assert.Error(t, err)

	// OpenFileAppend error
	_, err = OpenFileAppend(filepath.Join(parentFile, "child.log"), 0o644)
	assert.Error(t, err)

	// SetOutputFile error
	_, err = SetOutputFile(filepath.Join(parentFile, "child.log"))
	assert.Error(t, err)

	// AddFileWriter error
	_, err = AddFileWriter(LevelInfo, filepath.Join(parentFile, "child.log"))
	assert.Error(t, err)

	// AddJSONFile error
	_, err = AddJSONFile(LevelInfo, filepath.Join(parentFile, "child.log"))
	assert.Error(t, err)
}
