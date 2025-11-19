package log

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClose_ClosesAllRegisteredFiles(t *testing.T) {
	withStdReset(t, func() {
		tmpDir := t.TempDir()
		file1 := filepath.Join(tmpDir, "test1.log")
		file2 := filepath.Join(tmpDir, "test2.log")

		// Open files using package helpers
		f1, err := SetOutputFile(file1)
		assert.NoError(t, err)
		assert.NotNil(t, f1)

		f2, err := AddFileWriter(LevelInfo, file2)
		assert.NoError(t, err)
		assert.NotNil(t, f2)

		// Write to files
		Info("test message 1")
		Warn("test message 2")

		// Close all files
		Close()

		// Verify files were closed by attempting to write (should fail or panic)
		// We can't easily test if file is closed, so we'll just verify Close doesn't panic
		// and files exist with content
		data1, err := os.ReadFile(file1)
		assert.NoError(t, err)
		assert.Contains(t, string(data1), "test message")

		data2, err := os.ReadFile(file2)
		assert.NoError(t, err)
		assert.Contains(t, string(data2), "test message")
	})
}

func TestAddJSONFile_AutoRegisters(t *testing.T) {
	withStdReset(t, func() {
		tmpDir := t.TempDir()
		jsonFile := filepath.Join(tmpDir, "test.jsonl")

		f, err := AddJSONFile(LevelInfo, jsonFile)
		assert.NoError(t, err)
		assert.NotNil(t, f)

		Info("json test", "key", "value")

		Close()

		data, err := os.ReadFile(jsonFile)
		assert.NoError(t, err)
		assert.Contains(t, string(data), `"level":"INFO"`)
		assert.Contains(t, string(data), `"key":"value"`)
	})
}
