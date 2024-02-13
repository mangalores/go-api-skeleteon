package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExtractParamAndOperator(t *testing.T) {
	builder := QueryBuilder{operators: map[string]*Operator{}}
	operatorOk := &Operator{}
	builder.operators["eq"] = operatorOk

	name, actual, ok := builder.extractParamAndOperator("foo")
	assert.True(t, ok)
	assert.Equal(t, "foo", name)
	assert.Same(t, operatorOk, actual)

	name, actual, ok = builder.extractParamAndOperator("foo:bar")
	assert.False(t, ok)
	assert.Equal(t, "foo", name)
	assert.Nil(t, actual)
}

func TestMapOperator(t *testing.T) {
	builder := QueryBuilder{operators: map[string]*Operator{}}
	operatorOk := &Operator{}
	builder.operators["foo"] = operatorOk

	actualOk, ok := builder.mapOperator("foo")

	assert.True(t, ok)
	assert.Same(t, operatorOk, actualOk)

	actualNOK, ok := builder.mapOperator("bar")

	assert.False(t, ok)
	assert.Nil(t, actualNOK)
}

func TestDefaultTransformFn(t *testing.T) {
	values := []string{"foo", "bar", "baz"}

	actual := defaultTransformFn(values)

	assert.Equal(t, values, actual)
}
