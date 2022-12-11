package unstructured

import (
	"regexp"

	"github.com/kubectl-logz/kubectl-logz/internal/types"
)

var levelRx = regexp.MustCompile("\\b(error|warn|warning|info|debug|ERROR|WARN|WARNING|INFO|DEBUG)(ing)?\\b")

func Unmarshal(data []byte, entry *types.Entry) error {
	s := string(data)
	entry.Time = types.ParseTime(s)
	if find := levelRx.FindSubmatch(data); len(find) > 0 {
		entry.Level = types.ParseLevel(string(find[0]))
	} else {
		entry.Level = "info"
	}
	entry.ThreadID = types.ParseThreadId(s)
	entry.Msg = s
	return nil
}
