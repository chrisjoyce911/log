# github.com/chrisjoyce911/log

Drop-in replacement for the standard library `log` with:

- Levels: TRACE, VERBOSE, DEBUG, DETAIL, INFO, NOTICE, WARN, ERROR, CRITICAL, ALERT, FATAL, PANIC
- Multi-output routing: send different minimum levels to stdout, files, channels, JSON, etc.
- Handlers: text (writer), JSON, channel, and colored console
- Stdlib compatibility: `Print*`, `Fatal*`, `Panic*`, `SetFlags`, `SetPrefix`, `New`, and flags re-exported

## Quick start

```go
import (
  log "github.com/chrisjoyce911/log"
)

func main() {
  log.Println("hello")
  log.Info("structured", "key", "val")
}
```

## Multi-output routing

```go
log.AddWriter(log.LevelInfo, os.Stdout)  // INFO+ to stdout
log.AddWriter(log.LevelError, os.Stderr) // ERROR+ to stderr
log.AddHandler(log.LevelInfo, log.NewJSONHandler(os.Stderr)) // JSON to stderr
```

## JSON logging

```go
jh := log.NewJSONHandler(os.Stdout)
log.AddHandler(log.LevelInfo, jh)
log.Info("created", "id", 123)
```

## File helpers

```go
// Ensures parent directories exist and truncates existing file
f, _ := log.SetOutputFile("logs/app.log")
defer f.Close()
log.Info("to file")

// Attach JSON file
jf, _ := log.AddJSONFile(log.LevelInfo, "logs/app.jsonl")
defer jf.Close()
log.Info("json to file", "user", "alice")
```

## Colored console output

Color levels automatically when the output is a terminal (honors NO_COLOR). Or force colors on.

```go
// Auto mode (TTY only)
log.SetColoredOutput(log.LevelInfo, log.ColorOptions{Mode: log.ColorAuto})
log.Info("hello", "user", "alice")

// Force colors on and color specific parts
log.SetColoredOutput(log.LevelDebug, log.ColorOptions{
  Mode:         log.ColorOn,
  ColorLevel:   true,   // default
  ColorPrefix:  true,
  ColorMessage: true,
  ColorAttrs:   true,
})
log.SetPrefix("colored")
log.Warn("parts colored", "k", 1)
```

## Examples

- `examples/colored_console`: show all levels and per-part coloring
- `examples/mylog_console`: basic console usage
- `examples/mylog_file_text`: write text logs to file (with helper)
- `examples/mylog_file_json`: write JSON logs to file
- `examples/all`: combined demo
- `examples/levels_demo`: how min-level filters work
