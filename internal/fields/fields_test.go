package fields

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBraces(t *testing.T) {
	assert.ElementsMatch(t, Braces("[foo]-[bar] baz qux"), []string{"foo", "bar", "baz qux"})
}
