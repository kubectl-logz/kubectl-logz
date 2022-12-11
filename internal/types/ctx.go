package types

import (
	"fmt"
)

type Ctx struct {
	Hostname string `json:"hostname"`
}

func (c Ctx) String() string {
	return fmt.Sprintf("host=%s", c.Hostname)
}
