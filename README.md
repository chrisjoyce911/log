# github.com/chrisjoyce911/log

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/chrisjoyce911/log?tab=overview)
[![codecov](https://codecov.io/gh/chrisjoyce911/log/graph/badge.svg)](https://codecov.io/gh/chrisjoyce911/log)
[![Go Report Card](https://goreportcard.com/badge/github.com/chrisjoyce911/log)](https://goreportcard.com/report/github.com/chrisjoyce911/log)

Drop-in replacement for the standard library `log` with:

- Levels: TRACE, VERBOSE, DEBUG, DETAIL, INFO, NOTICE, WARN, ERROR, CRITICAL, ALERT, FATAL, PANIC
- Multi-output routing: send different minimum levels to stdout, files, channels, JSON, etc.
- Handlers: text (writer), JSON, channel, and colored console
- Stdlib compatibility: `Print*`, `Fatal*`, `Panic*`, `SetFlags`, `SetPrefix`, `New`, and flags re-exported
- HTTP middleware: colorized request/access logs with optional body preview

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

Colors are enabled by default. Use `ColorOff` to disable or `ColorAuto` for TTY detection (honors NO_COLOR).

```go
// Default: colors always on
log.SetColoredOutput(log.LevelInfo, log.ColorOptions{})
log.Info("hello", "user", "alice")

// Disable colors
log.SetColoredOutput(log.LevelInfo, log.ColorOptions{Mode: log.ColorOff})

// Auto mode (TTY only)
log.SetColoredOutput(log.LevelInfo, log.ColorOptions{Mode: log.ColorAuto})

// Customize which parts are colored
log.SetColoredOutput(log.LevelDebug, log.ColorOptions{
  ColorLevel:   true,   // default
  ColorPrefix:  true,
  ColorMessage: true,
  ColorAttrs:   true,
})
log.SetPrefix("colored")
log.Warn("parts colored", "k", 1)
```

## HTTP logging middleware

Colorized method and path (default: always on), severity from status (INFO <400, WARN 4xx, ERROR 5xx), optional body preview for POST/PUT/PATCH.

```go
mux := http.NewServeMux()
// register handlers...

// Colorized console (colors on by default)
log.SetColoredOutput(log.LevelDebug, log.ColorOptions{})
log.SetFlags(log.Ldate | log.Ltime)

// HTTP logging with body preview (2KB cap) and query in display path
h := log.HTTPLogging(mux, &log.HTTPLogOptions{
  IncludeQuery: true,
  LogPostBody:  true,
  MaxBodyBytes: 2048,
})

_ = http.ListenAndServe(":8080", h)
```

## Formatted logging

Use f-variants for printf-style logging.

```go
log.Infof("service %s started on port %d", "catalog", 8080)
log.Debugf("loaded %d features", 12)
```

See `examples/format_demo` for more.

## Testing helpers

Useful for unit/integration tests.

```go
// Prevent os.Exit in Fatal* during tests
log.SetExitFunc(func(code int) { /* capture */ })
defer log.SetExitFunc(nil)

// Deterministic time for new/default loggers
log.SetNowFunc(func() time.Time { return fixed })
defer log.SetNowFunc(nil)

// Silence default output during tests
log.SetTestingMode(true)
defer log.SetTestingMode(false)
```

## Examples

- `examples/colored_console`: show all levels and per-part coloring
- `examples/mylog_console`: basic console usage
- `examples/mylog_file_text`: write text logs to file (with helper)
- `examples/mylog_file_json`: write JSON logs to file
- `examples/all`: combined demo
- `examples/levels_demo`: how min-level filters work
- `examples/http_logger`: HTTP middleware with colored method/path and optional body logging
- `examples/format_demo`: demonstrate formatted logging with f-variants (Infof, Debugf, etc.)

## CI

GitHub Actions workflow at `.github/workflows/ci.yml` runs tests on pushes/PRs to `main`, caches modules, excludes `./examples` from coverage, and uploads coverage (Codecov supported).

## Test Coverage

[![Codecov sunburst](https://codecov.io/gh/chrisjoyce911/log/graphs/sunburst.svg?token=AI9ZF6PN5N)](https://codecov.io/gh/chrisjoyce911/log)
