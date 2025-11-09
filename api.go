package log

import (
	"fmt"
	"os"
)

// exitFunc allows tests to intercept process termination.
var exitFunc = os.Exit

// Print, Printf, Println implement stdlib's surface and log at Info level.
func Print(v ...any)                 { std.logf(LevelInfo, "%s", fmt.Sprint(v...)) }
func Printf(format string, v ...any) { std.logf(LevelInfo, format, v...) }
func Println(v ...any)               { std.logf(LevelInfo, "%s", trimNL(fmt.Sprintln(v...))) }

// Fatal variants log at Fatal and then exit(1).
func Fatal(v ...any) {
	std.logf(LevelFatal, "%s", fmt.Sprint(v...))
	exitFunc(1)
}
func Fatalf(format string, v ...any) {
	std.logf(LevelFatal, format, v...)
	exitFunc(1)
}
func Fatalln(v ...any) {
	std.logf(LevelFatal, "%s", trimNL(fmt.Sprintln(v...)))
	exitFunc(1)
}

// Panic variants log at Panic then panic.
func Panic(v ...any)                 { doPanic(fmt.Sprint(v...)) }
func Panicf(format string, v ...any) { doPanic(fmt.Sprintf(format, v...)) }
func Panicln(v ...any)               { doPanic(trimNL(fmt.Sprintln(v...))) }

func doPanic(msg string) {
	std.logf(LevelPanic, "%s", msg)
	panic(msg)
}

// Slog-like helpers on the default logger.
func Debug(msg string, kv ...any)    { std.logStructured(LevelDebug, msg, kv...) }
func Info(msg string, kv ...any)     { std.logStructured(LevelInfo, msg, kv...) }
func Warn(msg string, kv ...any)     { std.logStructured(LevelWarn, msg, kv...) }
func Error(msg string, kv ...any)    { std.logStructured(LevelError, msg, kv...) }
func Trace(msg string, kv ...any)    { std.logStructured(LevelTrace, msg, kv...) }
func Verbose(msg string, kv ...any)  { std.logStructured(LevelVerbose, msg, kv...) }
func Detail(msg string, kv ...any)   { std.logStructured(LevelDetail, msg, kv...) }
func Notice(msg string, kv ...any)   { std.logStructured(LevelNotice, msg, kv...) }
func Critical(msg string, kv ...any) { std.logStructured(LevelCritical, msg, kv...) }
func Alert(msg string, kv ...any)    { std.logStructured(LevelAlert, msg, kv...) }

// Logger method shims for stdlib-like API
func (l *Logger) Print(v ...any)                 { l.logf(LevelInfo, "%s", fmt.Sprint(v...)) }
func (l *Logger) Printf(format string, v ...any) { l.logf(LevelInfo, format, v...) }
func (l *Logger) Println(v ...any)               { l.logf(LevelInfo, "%s", trimNL(fmt.Sprintln(v...))) }

// Level helpers on Logger
func (l *Logger) Debug(msg string, kv ...any)    { l.logStructured(LevelDebug, msg, kv...) }
func (l *Logger) Info(msg string, kv ...any)     { l.logStructured(LevelInfo, msg, kv...) }
func (l *Logger) Warn(msg string, kv ...any)     { l.logStructured(LevelWarn, msg, kv...) }
func (l *Logger) Error(msg string, kv ...any)    { l.logStructured(LevelError, msg, kv...) }
func (l *Logger) Trace(msg string, kv ...any)    { l.logStructured(LevelTrace, msg, kv...) }
func (l *Logger) Verbose(msg string, kv ...any)  { l.logStructured(LevelVerbose, msg, kv...) }
func (l *Logger) Detail(msg string, kv ...any)   { l.logStructured(LevelDetail, msg, kv...) }
func (l *Logger) Notice(msg string, kv ...any)   { l.logStructured(LevelNotice, msg, kv...) }
func (l *Logger) Critical(msg string, kv ...any) { l.logStructured(LevelCritical, msg, kv...) }
func (l *Logger) Alert(msg string, kv ...any)    { l.logStructured(LevelAlert, msg, kv...) }

// Provide the formatMsg used in logger_core without importing fmt there.
func formatMsg(format string, v ...any) string { return fmt.Sprintf(format, v...) }
