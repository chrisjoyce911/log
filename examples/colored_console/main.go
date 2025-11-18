package main

import (
	log "github.com/chrisjoyce911/log"
)

func main() {
	// Colors are on by default. Use ColorOff to disable, or ColorAuto for TTY detection.
	log.SetColoredOutput(log.LevelAll, log.ColorOptions{})
	log.SetFlags(log.Ldate | log.Ltime)

	log.SetPrefix("demo")

	// Demonstrate all non-terminating levels with default palette.
	log.Trace("TRACE sample", "step", 1)
	log.Verbose("VERBOSE sample")
	log.Debug("DEBUG sample", "k", "v")
	log.Detail("DETAIL sample")
	log.Info("INFO sample")
	log.Notice("NOTICE sample")
	log.Warn("WARN sample")
	log.Error("ERROR sample")
	log.Critical("CRITICAL sample")
	log.Alert("ALERT sample")

	// Show per-part coloring: color level + prefix + message + attrs.
	log.SetColoredOutput(log.LevelAll, log.ColorOptions{
		ColorLevel:   true,
		ColorPrefix:  true,
		ColorMessage: true,
		ColorAttrs:   true,
	})
	log.SetPrefix("colored-parts")
	log.Info("message colored too", "user", "alice", "age", 30)

	// Note: FATAL and PANIC are not demonstrated since they terminate the program.
}
