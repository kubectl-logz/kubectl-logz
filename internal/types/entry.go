package types

import (
	"fmt"
	"strings"
	"time"
)

type Entry struct {
	Time     time.Time `json:"time"`
	Level    Level     `json:"level"`
	Msg      string    `json:"msg"`
	ThreadID string    `json:"threadId,omitempty"`
}

func (e Entry) String() string {
	var elems = []string{
		fmt.Sprintf("time=%s", e.Time.Format(time.RFC3339)),
		fmt.Sprintf("level=%v", e.Level),
		fmt.Sprintf("msg=%q", e.Msg),
	}
	if e.ThreadID != "" {
		elems = append(elems, fmt.Sprintf("threadId=%q", e.ThreadID))
	}
	return strings.Join(elems, " ")
}

func (e Entry) IsZero() bool {
	return e.Time.IsZero() && e.Level.IsZero() && e.ThreadID == "" && e.Msg == ""
}

func (e Entry) Validate() error {
	if e.Time.IsZero() {
		return fmt.Errorf("time is zero")
	}
	if e.Level.IsZero() {
		return fmt.Errorf("levelis zero")
	}
	if e.Msg == "" {
		return fmt.Errorf("msg is zero")
	}
	return nil
}
