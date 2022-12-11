package parser

import (
	"strings"

	"github.com/kubectl-logz/kubectl-logz/internal/parser/fields"
	"github.com/kubectl-logz/kubectl-logz/internal/parser/json"
	"github.com/kubectl-logz/kubectl-logz/internal/parser/logfmt"
	"github.com/kubectl-logz/kubectl-logz/internal/parser/unstructured"
	"github.com/kubectl-logz/kubectl-logz/internal/types"
)

type unmarshaller = func([]byte, *types.Entry) error

var unmarshalers = []unmarshaller{
	logfmt.Unmarshal,
	json.Unmarshal,
	fields.Unmarshaler(fields.Braces),
	fields.Unmarshaler(strings.Fields),
	unstructured.Unmarshal,
}

func unmarshall(line []byte) types.Entry {
	for _, u := range unmarshalers {
		entry := types.Entry{}
		err := u(line, &entry)
		if err == nil && !entry.Time.IsZero() {
			return entry
		}
	}
	return types.Entry{Level: "info", Msg: string(line)}
}

func Parse(lines <-chan []byte, entries chan<- types.Entry) {
	// last is used to merge multi-line entries together
	var last types.Entry
	for line := range lines {
		entry := unmarshall(line)
		if entry.Time.IsZero() && !last.Time.IsZero() {
			last.Msg = last.Msg + "\n" + entry.Msg
		} else {
			if !last.IsZero() {
				entries <- last
			}
			last = entry
		}
	}
	entries <- last
}
