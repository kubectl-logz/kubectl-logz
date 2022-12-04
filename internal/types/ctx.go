package types

import (
	"fmt"
	"strings"
)

type Ctx struct {
	Host string `json:"host"`
	PID  int    `json:"pid,omitempty"`
}

func (c Ctx) String() string {
	elems := []string{
		fmt.Sprintf("host=%s", c.Host),
	}
	if c.PID > 0 {
		elems = append(elems, fmt.Sprintf("pid=%d", c.PID))
	}
	return strings.Join(elems, " ")
}
