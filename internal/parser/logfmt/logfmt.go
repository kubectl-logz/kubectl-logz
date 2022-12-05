package logfmt

import (
	"github.com/kr/logfmt"
	"github.com/kubectl-logz/kubectl-logz/internal/types"
)

type lfentry struct {
	Time     string
	Level    string
	ThreadId string
	Msg      string
}

func Unmarshal(data []byte, entry *types.Entry) error {
	e := &lfentry{}
	err := logfmt.Unmarshal(data, e)
	if err != nil {
		return err
	}
	entry.Time = types.ParseTime(e.Time)
	entry.Level = types.ParseLevel(e.Level)
	entry.ThreadID = e.ThreadId
	entry.Msg = e.Msg
	return nil
}
