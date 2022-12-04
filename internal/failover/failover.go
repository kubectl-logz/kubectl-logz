package failover

import (
	"regexp"

	"github.com/kubectl-logz/kubectl-logz/internal/types"
)

var levelRx = regexp.MustCompile("\\b(error|warn|warning|info|debug|ERROR|WARN|WARNING|INFO|DEBUG)(ing)?\\b")

func Unmarshal(l []byte, v interface{}) error {
	r := v.(*types.Entry)
	s := string(l)
	r.Time = types.ParseTime(s)
	if find := levelRx.FindSubmatch(l); len(find) > 0 {
		r.Level = types.ParseLevel(string(find[0]))
	} else {
		r.Level = "info"
	}
	r.ThreadID = types.ParseThreadId(s)
	r.Msg = s
	return nil
}
