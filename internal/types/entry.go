package types

import (
	"fmt"
	"strings"
	"time"
)

type Entry struct {
	Time     time.Time `json:"time,omitempty"`
	Level    Level     `json:"level,omitempty"`
	Msg      string    `json:"msg,omitempty"`
	ThreadID string    `json:"threadId,omitempty"`
}

func (e Entry) String() string {
	var elems []string
	if !e.Time.IsZero() {
		elems = append(elems, fmt.Sprintf("time=%s", e.Time))
	}
	if e.Level != "" {
		elems = append(elems, fmt.Sprintf("level=%v", e.Level))
	}
	if e.ThreadID != "" {
		elems = append(elems, fmt.Sprintf("threadId=%q", e.ThreadID))
	}
	if e.Msg != "" {
		elems = append(elems, fmt.Sprintf("msg=%q", e.Msg))
	}
	return strings.Join(elems, " ")
}
