package query_builder

import (
	"fmt"

	"github.com/mangalores/go-api-skeleton/pkg/db"
)

type InvalidParamValueErr struct {
	name    string
	numeric bool
}

func NewInvalidParamValueErr(name string, numeric bool) InvalidParamValueErr {
	return InvalidParamValueErr{
		name,
		numeric,
	}
}

func (e InvalidParamValueErr) Error() string {
	msg := fmt.Sprintf("param %s has invalid value", e.name)
	if e.numeric {
		msg = msg + " must be numeric"
	}

	return msg
}

type MaxLimitExceededErr struct {
	maxLimit int
}

func NewMaxLimitExceededErr() MaxLimitExceededErr {
	return MaxLimitExceededErr{
		maxLimit,
	}
}

func (e MaxLimitExceededErr) Error() string {
	return fmt.Sprintf("collection limit  exceeded (max limit: %d)", e.maxLimit)
}

type InvalidFilterErr struct {
	filter db.Filter
}

func NewInvalidFilterErr(f db.Filter) InvalidFilterErr {
	return InvalidFilterErr{
		f,
	}
}

func (e InvalidFilterErr) Error() string {
	return fmt.Sprintf("invalid filter value for name: %s and operator %s", e.filter.FieldName, e.filter.Operator)
}

type InvalidMultipleValuesErr struct {
	filter db.Filter
}

func NewInvalidMultipleValuesErr(f db.Filter) InvalidMultipleValuesErr {
	return InvalidMultipleValuesErr{
		f,
	}
}

func (e InvalidMultipleValuesErr) Error() string {
	return fmt.Sprintf("multiple values not allowed for filter withr name: %s and operator %s", e.filter.FieldName, e.filter.Operator)
}

type InvalidEmbedErr struct {
	invalidName string
}

func NewInvalidEmbedErr(invalidName string) InvalidEmbedErr {
	return InvalidEmbedErr{
		invalidName,
	}
}

func (e InvalidEmbedErr) Error() string {
	return fmt.Sprintf("invalid embed name requested: %s", e.invalidName)
}
