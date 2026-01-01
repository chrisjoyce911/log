package log

import "strconv"

// Level represents the severity of a log record.
//
// Lower numeric values indicate less severe events; higher values are more severe.
// Use LevelAll to capture everything and LevelOff to disable all logging for a sink.
type Level int

const (
	// Control levels
	LevelAll Level = -1000 // captures everything
	LevelOff Level = 1000  // disables all

	// Diagnostic levels
	LevelTrace   Level = -8
	LevelVerbose Level = -6
	LevelDebug   Level = -4
	LevelDetail  Level = -2

	// Operational levels
	LevelInfo     Level = 0
	LevelNotice   Level = 2
	LevelWarn     Level = 4
	LevelError    Level = 8
	LevelCritical Level = 10
	LevelAlert    Level = 12
	LevelFatal    Level = 14
	LevelPanic    Level = 16
)

func (l Level) String() string {
	switch l {
	case LevelAll:
		return "ALL     "
	case LevelOff:
		return "OFF     "
	case LevelTrace:
		return "TRACE   "
	case LevelVerbose:
		return "VERBOSE "
	case LevelDebug:
		return "DEBUG   "
	case LevelDetail:
		return "DETAIL  "
	case LevelInfo:
		return "INFO    "
	case LevelNotice:
		return "NOTICE  "
	case LevelWarn:
		return "WARN    "
	case LevelError:
		return "ERROR   "
	case LevelCritical:
		return "CRITICAL"
	case LevelAlert:
		return "ALERT   "
	case LevelFatal:
		return "FATAL   "
	case LevelPanic:
		return "PANIC.  "
	default:
		return "LEVEL(" + strconv.Itoa(int(l)) + ")"
	}
}
