package log

import (
	"fmt"
	"strings"
	"time"
)

// formatTimestamp renders a timestamp using stdlib log flags for compatibility.
// Matches behaviors of LUTC, Ldate, Ltime, and Lmicroseconds.
func formatTimestamp(t time.Time, flags int) string {
	if flags&LUTC != 0 {
		t = t.UTC()
	}
	haveDate := flags&Ldate != 0
	haveTime := flags&Ltime != 0
	if !haveDate && !haveTime {
		return ""
	}
	var b strings.Builder
	if haveDate {
		y, m, d := t.Date()
		write4(&b, y)
		b.WriteByte('/')
		write2(&b, int(m))
		b.WriteByte('/')
		write2(&b, d)
	}
	if haveTime {
		if haveDate {
			b.WriteByte(' ')
		}
		h, m, s := t.Clock()
		write2(&b, h)
		b.WriteByte(':')
		write2(&b, m)
		b.WriteByte(':')
		write2(&b, s)
		if flags&Lmicroseconds != 0 {
			us := t.Nanosecond() / 1_000
			b.WriteByte('.')
			writeN(&b, us, 6)
		}
	}
	return b.String()
}

func write2(b *strings.Builder, v int) {
	b.WriteByte(byte('0' + (v/10)%10))
	b.WriteByte(byte('0' + v%10))
}

func write4(b *strings.Builder, v int) {
	b.WriteByte(byte('0' + (v/1000)%10))
	b.WriteByte(byte('0' + (v/100)%10))
	b.WriteByte(byte('0' + (v/10)%10))
	b.WriteByte(byte('0' + v%10))
}

func writeN(b *strings.Builder, v int, n int) {
	// write v zero-padded to n digits
	s := fmt.Sprintf("%0*d", n, v)
	b.WriteString(s)
}

func trimNL(s string) string {
	if strings.HasSuffix(s, "\n") {
		return strings.TrimSuffix(s, "\n")
	}
	return s
}
