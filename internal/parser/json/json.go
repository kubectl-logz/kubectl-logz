package json

import (
	"encoding/json"

	"github.com/kubectl-logz/kubectl-logz/internal/types"
)

func Unmarshal(data []byte, entry *types.Entry) error {
	return json.Unmarshal(data, entry)
}
