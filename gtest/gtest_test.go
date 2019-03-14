package gtest

import (
	"testing"
)

func TestAuxiliary(t *testing.T) {
	Equal(t, true, true)
	Unequal(t, true, false)
	PrintMyName(t)
}
