package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLevelString(t *testing.T) {
	assert.Equal(t, "ALL     ", LevelAll.String())
	assert.Equal(t, "OFF     ", LevelOff.String())
	assert.Equal(t, "TRACE   ", LevelTrace.String())
	assert.Equal(t, "VERBOSE ", LevelVerbose.String())
	assert.Equal(t, "DEBUG   ", LevelDebug.String())
	assert.Equal(t, "DETAIL  ", LevelDetail.String())
	assert.Equal(t, "INFO    ", LevelInfo.String())
	assert.Equal(t, "NOTICE  ", LevelNotice.String())
	assert.Equal(t, "WARN    ", LevelWarn.String())
	assert.Equal(t, "ERROR   ", LevelError.String())
	assert.Equal(t, "CRITICAL", LevelCritical.String())
	assert.Equal(t, "ALERT   ", LevelAlert.String())
	assert.Equal(t, "FATAL   ", LevelFatal.String())
	assert.Equal(t, "PANIC   ", LevelPanic.String())
	// Unknown
	assert.Equal(t, "LEVEL(123)", Level(123).String())
}
