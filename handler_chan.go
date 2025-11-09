package log

import "fmt"

// StringChanHandler sends formatted log lines to a string channel.
type StringChanHandler struct {
	C chan<- string
}

func (h *StringChanHandler) Handle(r Record) error {
	ts := formatTimestamp(r.Time, r.Flags)
	if ts != "" {
		h.C <- fmt.Sprintf("%s %s %s", ts, r.Level.String(), r.Message)
		return nil
	}
	h.C <- fmt.Sprintf("%s %s", r.Level.String(), r.Message)
	return nil
}
