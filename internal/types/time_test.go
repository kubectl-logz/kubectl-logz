package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTime(t *testing.T) {

	assert.False(t, ParseTime("2022-12-04T11:34:26,673-0800").IsZero())
}