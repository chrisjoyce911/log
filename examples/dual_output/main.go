package main

import (
	log "github.com/chrisjoyce911/log"
)

func main() {
	// Set up dual output: CLI (stdout) with colors and file (plain text)
	// No need to defer f.Close() - log.Close() handles it

	// Colored output to CLI (stdout)
	log.SetColoredOutput(log.LevelDebug, log.ColorOptions{})

	// Add plain text file output
	_, err := log.AddFileWriter(log.LevelInfo, "./logs/ac-sync.log")
	if err != nil {
		log.Fatal("failed to open log file: ", err)
	}

	// Clean up all files on exit
	defer log.Close()

	log.Info("goes to both stdout and file", "user", "alice")
	log.Warn("warning message", "code", 42)
	log.Error("error occurred", "reason", "timeout")

	// Add another file for debug logs
	_, err = log.AddFileWriter(log.LevelDebug, "./logs/debug.log")
	if err != nil {
		log.Fatal("failed to open debug log: ", err)
	}

	log.Debug("debug message only in debug.log and stdout")
	log.Info("info message in all outputs")
}
