package db

import (
	"fmt"
	"reflect"
)

type Direction string

const (
	ASC  Direction = "asc"
	DESC Direction = "desc"
)

type QueryObject interface {
	Model() interface{}
	Error() error
	SetError(err error)
	Result() interface{}
	SetResult(result interface{})
	Preloads() []Preload
}

type FilteredQueryObject interface {
	QueryObject
	Filters() []Filter
}

type SlicedQueryObject interface {
	FilteredQueryObject
	Slice() *Slice
}

type Preload struct {
	Name       string
	Conditions []interface{}
}

type Query struct {
	modelType  string
	resultType string
	model      interface{}
	result     interface{}
	error      error
	preloads   *[]Preload
}

func NewQuery(model interface{}) *Query {
	return &Query{
		model: model,
	}
}

func (q *Query) Model() interface{} {
	return q.model
}

func (q *Query) SetModel(model interface{}) {
	q.modelType = fmt.Sprintf("%T", model)
	q.model = model
}

func (q *Query) Error() error {
	return q.error
}

func (q *Query) SetError(err error) {
	q.error = err
}

func (q *Query) Result() interface{} {
	return q.result
}

func (q *Query) SetResult(result interface{}) {
	q.resultType = fmt.Sprintf("%T", result)
	q.result = result
}

func (q *Query) Type() string {
	return reflect.TypeOf(q.Model).String()
}

func (q *Query) AddPreload(preload Preload) {
	if q.preloads == nil {
		q.preloads = &[]Preload{}
	}

	*q.preloads = append(*q.preloads, preload)
}

func (q *Query) Preloads() []Preload {
	if q.preloads == nil {
		return []Preload{}
	}

	return *q.preloads
}

func (q *Query) SetPreloads(preloads *[]Preload) {
	if q.preloads == nil {
		q.preloads = &[]Preload{}
	}

	q.preloads = preloads
}

type Filter struct {
	FieldName string
	Operator  string
	Value     interface{}
}

type FilterQuery struct {
	Query
	filters []Filter
}

func NewFilterQuery(model interface{}) *FilterQuery {
	return &FilterQuery{
		Query: Query{
			model: model,
		},
		filters: make([]Filter, 0),
	}
}

func (q *FilterQuery) Filters() []Filter {
	return q.filters
}

func (q *FilterQuery) SetFilters(filters []Filter) {
	q.filters = filters
}

type Slice struct {
	Offset int
	Limit  int
	Total  int64
	Sort   []Sort
}

type Sort struct {
	FieldName string
	Direction Direction
}

type CollectionQuery struct {
	FilterQuery
	slice *Slice
}

func NewCollectionQuery(model interface{}) *CollectionQuery {
	return &CollectionQuery{
		FilterQuery: FilterQuery{
			Query:   Query{model: model},
			filters: make([]Filter, 0),
		},
	}
}

func (q *CollectionQuery) Slice() *Slice {
	return q.slice
}

func (q *CollectionQuery) SetSlice(slice *Slice) {
	q.slice = slice
}
