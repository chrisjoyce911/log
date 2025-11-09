package main

import (
	"fmt"
	"io"
	"os"

	log "github.com/chrisjoyce911/log"
)

// Run executes the demo, writing section headers to out and routing logs to logOut.
func Run(out io.Writer, logOut io.Writer) {
	levels := []struct {
		name string
		lvl  log.Level
	}{
		{"ALL", log.LevelAll},
		{"TRACE", log.LevelTrace},
		{"VERBOSE", log.LevelVerbose},
		{"DEBUG", log.LevelDebug},
		{"DETAIL", log.LevelDetail},
		{"INFO", log.LevelInfo},
		{"NOTICE", log.LevelNotice},
		{"WARN", log.LevelWarn},
		{"ERROR", log.LevelError},
		{"CRITICAL", log.LevelCritical},
		{"ALERT", log.LevelAlert},
		// FATAL and PANIC are excluded because they terminate the program.
		{"OFF", log.LevelOff},
	}

	// Use a predictable timestamp format for readability.
	log.SetFlags(log.Ldate | log.Ltime)

	for _, cfg := range levels {
		fmt.Fprintf(out, "\n=== Min level: %s ===\n", cfg.name)

		// Route logs only to stdout with the min level under test.
		log.SetOutput(io.Discard)
		target := logOut
		if target == nil {
			target = os.Stdout
		}
		log.AddWriter(cfg.lvl, target)

		log.Println("Log Level set to ", cfg.name)

		log.Println("[Println] hello")
		log.Println("[Println] hello")
		log.Debug("[Debug] hello")
		log.Trace("[Trace] hello")
		log.Verbose("[Verbose] hello")
		log.Detail("[Detail] hello")
		log.Info("[Info] hello")
		log.Notice("[Notice] hello")
		log.Warn("[Warn] hello")
		log.Error("[Error] hello")
		log.Critical("[Critical] hello")
		log.Alert("[Alert] hello")
	}
}

func main() {
	Run(os.Stdout, os.Stdout)
}
