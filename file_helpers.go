package log

import (
	"os"
	"path/filepath"
)

// ensureDirForFile creates parent directories for the given file path if needed.
func ensureDirForFile(path string) error {
	dir := filepath.Dir(path)
	if dir == "." || dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}

// OpenFileTruncate opens the file for writing, creating parent directories if required,
// and truncates any existing file.
func OpenFileTruncate(path string, perm os.FileMode) (*os.File, error) {
	if err := ensureDirForFile(path); err != nil {
		return nil, err
	}
	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
}

// OpenFileAppend opens the file for appending, creating parent directories if required.
func OpenFileAppend(path string, perm os.FileMode) (*os.File, error) {
	if err := ensureDirForFile(path); err != nil {
		return nil, err
	}
	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, perm)
}

// SetOutputFile ensures the file exists (creating parent dirs) and sets the default
// logger's output to that file, truncating any existing content. It returns the open *os.File.
// The file is registered for automatic cleanup when log.Close() is called.
func SetOutputFile(path string) (*os.File, error) {
	f, err := OpenFileTruncate(path, 0o644)
	if err != nil {
		return nil, err
	}
	SetOutput(f)
	registerFile(f)
	return f, nil
}

// AddFileWriter ensures the file exists and attaches it as a writer for the given minimum level
// on the default logger, truncating any existing content. It returns the open *os.File.
// The file is registered for automatic cleanup when log.Close() is called.
func AddFileWriter(minLevel Level, path string) (*os.File, error) {
	f, err := OpenFileTruncate(path, 0o644)
	if err != nil {
		return nil, err
	}
	AddWriter(minLevel, f)
	registerFile(f)
	return f, nil
}

// AddJSONFile ensures the file exists and attaches a JSON handler for the given minimum level
// on the default logger, truncating any existing content. It returns the open *os.File.
// The file is registered for automatic cleanup when log.Close() is called.
func AddJSONFile(minLevel Level, path string) (*os.File, error) {
	f, err := OpenFileTruncate(path, 0o644)
	if err != nil {
		return nil, err
	}
	AddHandler(minLevel, NewJSONHandler(f))
	registerFile(f)
	return f, nil
}
