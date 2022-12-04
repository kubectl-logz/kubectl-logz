package internal

import (
	"github.com/kr/logfmt"
	"github.com/kubectl-logz/kubectl-logz/internal/failover"
	"github.com/kubectl-logz/kubectl-logz/internal/fields"
	"github.com/kubectl-logz/kubectl-logz/internal/types"
	"k8s.io/apimachinery/pkg/util/json"
)

func parse(l []byte) (types.Entry, error) {
	r := types.Entry{}
	if _ = logfmt.Unmarshal(l, &r); r.Valid() {
		return r, nil
	}
	if _ = json.Unmarshal(l, &r); r.Valid() {
		return r, nil
	}
	if fields.Unmarshal(l, &r); r.Valid() {
		return r, nil
	}
	failover.Unmarshal(l, &r)
	return r, nil
}
