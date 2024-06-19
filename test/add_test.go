package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func add(a, b int) int {
	return a + b
}
func testAdd(t *testing.T) {

	assert.Equal(t, 2, add(1, 1))
}
