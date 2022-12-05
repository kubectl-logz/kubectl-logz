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
	return types.Entry{Msg: string(line)}
}

func Parse(lines <-chan []byte, entries chan<- types.Entry) {
	for line := range lines {
		entries <- unmarshall(line)
	}
}
