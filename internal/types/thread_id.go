package types

import (
	"regexp"
	"strings"
)

var threadIdRx = regexp.MustCompile("\\S*thread\\S*")

func ParseThreadId(s string) string {
	return strings.Trim(threadIdRx.FindString(s), "\"")
}
