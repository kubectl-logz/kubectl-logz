package types

import (
	"fmt"
	"strings"
	"time"
)

type Entry struct {
	Time     time.Time `json:"time"`
	Level    Level     `json:"level"`
	Msg      string    `json:"msg,omitempty"`
	ThreadID string    `json:"threadId,omitempty"`
}

func (r Entry) String() string {
	elems := []string{
		fmt.Sprintf("time=%v", r.Time.Format(time.RFC3339)),
		fmt.Sprintf("level=%v", r.Level),
	}
	if r.Msg != "" {
		elems = append(elems, fmt.Sprintf("msg=%q", r.Msg))
	}
	return strings.Join(
		elems, " ")
}

func (r Entry) Valid() bool {
	return !r.Time.IsZero() && r.Level != ""
}
