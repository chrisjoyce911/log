package log

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatTimestampCombinations(t *testing.T) {
	// Fixed time for deterministic output
	tm := time.Date(2025, 11, 8, 12, 34, 56, 123456000, time.FixedZone("X", 9*3600))

	// No flags -> empty
	assert.Equal(t, "", formatTimestamp(tm, 0))

	// Date only
	assert.Equal(t, "2025/11/08", formatTimestamp(tm, Ldate))

	// Time only (local zone)
	assert.Equal(t, "12:34:56", formatTimestamp(tm, Ltime))

	// Date + Time
	assert.Equal(t, "2025/11/08 12:34:56", formatTimestamp(tm, Ldate|Ltime))

	// With microseconds
	assert.Equal(t, "2025/11/08 12:34:56.123456", formatTimestamp(tm, Ldate|Ltime|Lmicroseconds))

	// With UTC
	out := formatTimestamp(tm, Ldate|Ltime|LUTC)
	assert.True(t, strings.HasPrefix(out, "2025/11/08 "))
}
