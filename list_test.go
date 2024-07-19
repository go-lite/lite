package lite

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewList(t *testing.T) {
	l := NewList([]string{"a", "b", "c"})
	assert.Equal(t, 3, l.Length)
}
