package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func subtract(a, b int) int {
	return a - b
}
func testSubtract(t *testing.T) {

	assert.Equal(t, 2, subtract(3, 1))
}
