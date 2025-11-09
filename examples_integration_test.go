package log

import (
	"bufio"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func runInDir(t *testing.T, dir string, args ...string) (string, string, error) {
	t.Helper()
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), "", err
}

func TestExampleStdlog(t *testing.T) {
	dir := filepath.Join("examples", "stdlog")
	out, _, err := runInDir(t, dir, "go", "run", ".")
	assert.NoError(t, err, out)
	assert.Contains(t, out, "hello from std log")
	assert.Contains(t, out, "user:")
}

func TestExampleMylogConsole(t *testing.T) {
	dir := filepath.Join("examples", "mylog_console")
	out, _, err := runInDir(t, dir, "go", "run", ".")
	assert.NoError(t, err, out)
	assert.Contains(t, out, "hello from my log (console)")
	assert.Contains(t, out, "user:")
}

func TestExampleMylogFileText(t *testing.T) {
	dir := filepath.Join("examples", "mylog_file_text")
	defer os.Remove(filepath.Join(dir, "logs", "app.log"))
	out, _, err := runInDir(t, dir, "go", "run", ".")
	assert.NoError(t, err, out)
	b, err := os.ReadFile(filepath.Join(dir, "logs", "app.log"))
	assert.NoError(t, err)
	content := string(b)
	assert.Contains(t, content, "hello from my log (file text)")
	assert.Contains(t, content, "user:")
}

func TestExampleMylogFileJSON(t *testing.T) {
	dir := filepath.Join("examples", "mylog_file_json")
	defer os.Remove(filepath.Join(dir, "app.jsonl"))
	out, _, err := runInDir(t, dir, "go", "run", ".")
	assert.NoError(t, err, out)
	f, err := os.Open(filepath.Join(dir, "app.jsonl"))
	assert.NoError(t, err)
	defer f.Close()
	s := bufio.NewScanner(f)
	var lines []string
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	assert.GreaterOrEqual(t, len(lines), 1)
	var m map[string]any
	assert.NoError(t, json.Unmarshal([]byte(lines[0]), &m))
	assert.Equal(t, "INFO", m["level"])
	assert.True(t, strings.Contains(lines[0], "\"msg\":"))
}

func TestExampleAll(t *testing.T) {
	dir := filepath.Join("examples", "all")
	defer os.Remove(filepath.Join(dir, "all_app.log"))
	defer os.Remove(filepath.Join(dir, "all_app.jsonl"))
	out, _, err := runInDir(t, dir, "go", "run", ".")
	assert.NoError(t, err, out)
	assert.Contains(t, out, "[stdlog] hello")
	assert.Contains(t, out, "[mylog-console] hello")
	_, err = os.Stat(filepath.Join(dir, "all_app.log"))
	assert.NoError(t, err)
	_, err = os.Stat(filepath.Join(dir, "all_app.jsonl"))
	assert.NoError(t, err)
}

func TestExampleLevelsDemo(t *testing.T) {
	dir := filepath.Join("examples", "levels_demo")
	out, _, err := runInDir(t, dir, "go", "run", ".")
	assert.NoError(t, err, out)
	assert.Contains(t, out, "=== Min level: ALL ===")
	assert.Contains(t, out, "=== Min level: OFF ===")
}
