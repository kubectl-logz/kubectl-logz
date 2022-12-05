package fields

import (
	"fmt"
	"strings"

	"github.com/kubectl-logz/kubectl-logz/internal/types"
)

func Unmarshaler(splitter func(string) []string) func(data []byte, entry *types.Entry) error {
	return func(data []byte, entry *types.Entry) error {
		var msg []string
		for _, s := range splitter(string(data)) {
			v := strings.Trim(s, "[]")
			if entry.Time.IsZero() {
				entry.Time = types.ParseTime(v)
				if !entry.Time.IsZero() {
					continue
				}
			}
			if entry.Level.IsZero() {
				entry.Level = types.ParseLevel(strings.TrimSpace(v))
				if !entry.Level.IsZero() {
					continue
				}
			}
			if threadID := types.ParseThreadId(v); entry.ThreadID == "" {
				entry.ThreadID = threadID
				if entry.ThreadID != "" {
					continue
				}
			}
			msg = append(msg, s)
		}
		if entry.Time.IsZero() {
			return fmt.Errorf("could not parse")
		}
		entry.Msg = strings.Join(msg, " ")
		return nil
	}
}
