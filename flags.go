package log

import stdlog "log"

// Re-export stdlib log flags for drop-in compatibility.
const (
	Ldate         = stdlog.Ldate
	Ltime         = stdlog.Ltime
	Lmicroseconds = stdlog.Lmicroseconds
	Llongfile     = stdlog.Llongfile
	Lshortfile    = stdlog.Lshortfile
	LUTC          = stdlog.LUTC
	Lmsgprefix    = stdlog.Lmsgprefix
	LstdFlags     = stdlog.LstdFlags
)
