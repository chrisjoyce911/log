// Package log is a drop-in replacement for the standard library "log" package
// with level-aware, multi-output routing and a structured logging surface inspired
// by "log/slog".
//
// Goals
//
//  1. Familiarity: Keep the core API of the stdlib log (Print, Printf, Println,
//     Fatal, Panic, SetFlags, SetPrefix, New) so adopting this package only
//     requires changing the import path to this module.
//  2. Levels: Provide conventional levels (Debug, Info, Warn, Error) and
//     helpers like Debug, Info, Warn, Error methods and functions similar to slog.
//  3. Multi-output routing: Configure which levels go to which outputs, such as
//     console (stdout/stderr), files, or channels. You can add multiple outputs
//     and set different minimum levels per output.
//  4. Extensibility: A simple Handler interface allows custom outputs.
//
// Quick start
//
//		import (
//		    stdlog "log"
//		    mylog "chrisjoyce911/log"
//		    "os"
//		)
//
//		func main() {
//		    // Drop-in std log usage
//		    mylog.Println("standard logger")
//		    mylog.SetFlags(mylog.LstdFlags | mylog.Lmicroseconds)
//		    mylog.Println("with micro")
//		    mylog.SetFlags(mylog.LstdFlags | mylog.Lshortfile)
//		    mylog.Println("with file/line (formatting placeholders)")
//
//		    // Create a new logger like stdlib log.New, wired to stdout
//		    l := mylog.New(os.Stdout, "my:", mylog.LstdFlags)
//		    l.Println("from mylog")
//
//	    // Configure multi-output routing
//	    mylog.AddWriter(mylog.LevelInfo, os.Stdout)   // info+ to stdout
//	    mylog.AddWriter(mylog.LevelError, os.Stderr)  // error+ to stderr
//
//	    // JSON output
//	    jh := mylog.NewJSONHandler(os.Stderr)
//	    mylog.AddHandler(mylog.LevelInfo, jh)
//
//	    mylog.Info("hello", "key", "val")
//	    mylog.Error("boom", "err", "example")
//
//		    // Interop: you can still use the stdlib log alongside if you alias imports
//		    stdlog.Println("stdlib log still works (different import)")
//		}
//
// # Notes
//
//   - Flags and file/line formatting are provided as placeholders initially; the
//     constants are re-exported for drop-in compatibility, but formatting can be
//     refined over time.
//   - Structured attributes (key-value pairs) are accepted by the slog-like APIs
//     but are formatted minimally in the default text writer.
//   - Handlers can implement arbitrary formats (text, JSON) and destinations.
package log
