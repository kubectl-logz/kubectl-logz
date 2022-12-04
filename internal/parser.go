package internal

import (
	"strings"

	"github.com/kr/logfmt"
	"github.com/kubectl-logz/kubectl-logz/internal/failover"
	"github.com/kubectl-logz/kubectl-logz/internal/fields"
	"github.com/kubectl-logz/kubectl-logz/internal/types"
	"k8s.io/apimachinery/pkg/util/json"
)

type unmarshaller = func([]byte, any) error

var unmarshalers = []unmarshaller{
	logfmt.Unmarshal,
	json.Unmarshal,
	fields.Unmarshaler(fields.Braces),
	fields.Unmarshaler(strings.Fields),
	failover.Unmarshal,
}

func parse(l []byte) types.Entry {
	r := types.Entry{}
	for _, u := range unmarshalers {
		_ = u(l, &r)
		if r.Valid() {
			break
		}
	}
	return r
}
