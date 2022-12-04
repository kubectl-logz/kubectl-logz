package fields

import (
	"strings"

	"github.com/kubectl-logz/kubectl-logz/internal/types"
)

func Unmarshaler(splitter func(string) []string) func([]byte, any) error {
	return func(l []byte, v any) error {
		r := v.(*types.Entry)
		var msg []string
		for _, s := range splitter(string(l)) {
			v := strings.Trim(s, "[]")
			if r.Time.IsZero() {
				r.Time = types.ParseTime(v)
				if !r.Time.IsZero() {
					continue
				}
			}
			if r.Level.IsZero() {
				r.Level = types.ParseLevel(strings.TrimSpace(v))
				if !r.Level.IsZero() {
					continue
				}
			}
			if threadID := types.ParseThreadId(v); r.ThreadID == "" {
				r.ThreadID = threadID
				if r.ThreadID != "" {
					continue
				}
			}
			msg = append(msg, s)
		}
		if !r.Time.IsZero() && r.Level.IsZero() {
			r.Level = "info"
		}
		r.Msg = strings.Join(msg, " ")
		return nil
	}
}
