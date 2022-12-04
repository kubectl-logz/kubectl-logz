package fields

import (
	"strings"

	"github.com/kubectl-logz/kubectl-logz/internal/types"
)

func Unmarshal(l []byte, r *types.Entry) {
	var msg []string
	for _, s := range strings.Fields(string(l)) {
		v := strings.Trim(s, "[]")
		if r.Time.IsZero() {
			r.Time = types.ParseTime(v)
			if !r.Time.IsZero() {
				continue
			}
		}
		if r.Level.IsZero() {
			r.Level = types.ParseLevel(v)
			if !r.Level.IsZero() {
				continue
			}
		}
		msg = append(msg, s)
	}
	if !r.Time.IsZero() && r.Level.IsZero() {
		r.Level = "info"
	}
	r.Msg = strings.Join(msg, " ")
}
