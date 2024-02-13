package api

import (
	"fmt"
)

type TypeErr struct {
	expected interface{}
	actual   interface{}
}

func NewTypeErr(expected interface{}, actual interface{}) TypeErr {
	return TypeErr{
		expected: expected,
		actual:   actual,
	}
}

func (e TypeErr) Error() string {
	return fmt.Sprintf("wrong type, expected '%T', got '%T'", e.expected, e.actual)
}

type UnsupportedResultTypeErr struct {
	result interface{}
}

func NewUnsupportedQueryType(result interface{}) UnsupportedResultTypeErr {
	return UnsupportedResultTypeErr{
		result: result,
	}
}

func (e UnsupportedResultTypeErr) Error() string {
	return fmt.Sprintf("unsupported query object of type %T", e.result)
}
