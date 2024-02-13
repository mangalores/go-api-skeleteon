package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContains_OK(t *testing.T) {
	input := []string{"foo", "bar", "baz"}

	assert.True(t, Contains(input, "foo"))
	assert.True(t, Contains(input, "bar"))
	assert.True(t, Contains(input, "baz"))
	assert.False(t, Contains(input, "not_contains"))
}

type TestStruct struct {
}

func TestStripPointer_OK(t *testing.T) {
	var actual interface{}

	actual = StripPointer(TestStruct{})
	assert.IsType(t, TestStruct{}, actual)

	actual = StripPointer(&TestStruct{})
	assert.IsType(t, TestStruct{}, actual)

	actual = StripPointer(&TestStruct{})
	assert.IsType(t, TestStruct{}, actual)

	actual = StripPointer(&[]TestStruct{})
	assert.IsType(t, []TestStruct{}, actual)

	pointer1 := &TestStruct{}
	pointer2 := &pointer1
	actual = StripPointer(pointer2)
	assert.IsType(t, TestStruct{}, actual)
}
