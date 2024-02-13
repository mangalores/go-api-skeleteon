package db

import (
	"errors"
	"fmt"
	"gorm.io/gorm/clause"
	"reflect"
	"sync"

	"gorm.io/gorm"
	gormSchema "gorm.io/gorm/schema"
)

type QueryHandler struct {
	db *gorm.DB
}

type UnsupportedQueryTypeErr struct {
	query interface{}
}

func (e UnsupportedQueryTypeErr) Error() string {
	return fmt.Sprintf("unsupported query object of type %s", reflect.TypeOf(e.query))
}

func NewQueryHandler(db *gorm.DB) *QueryHandler {
	return &QueryHandler{
		db,
	}
}

func (h *QueryHandler) Supports(t interface{}) bool {
	h.db.Model(t)
	return true
}

func (h *QueryHandler) Handle(query QueryObject) QueryObject {
	stmt, schema, err := h.buildStatement(query.Model())
	if err != nil {
		query.SetError(err)
		return query
	}

	var result interface{}

	switch query.(type) {
	case SlicedQueryObject:
		buildPreloads(stmt, query.Preloads())
		buildFilter(stmt, query.(FilteredQueryObject), schema)
		buildSelection(stmt, query.(SlicedQueryObject), schema)
	case FilteredQueryObject:
		buildPreloads(stmt, query.Preloads())
		buildFilter(stmt, query.(FilteredQueryObject), schema)
	case QueryObject:
		buildPreloads(stmt, query.Preloads())
	}

	result = h.buildResult(query.Model())
	if query.Error() != nil {
		return query
	}

	stmt.Find(result)
	query.SetResult(result)

	return query
}

func (h *QueryHandler) buildStatement(result interface{}) (stmt *gorm.DB, schema *gormSchema.Schema, err error) {
	stmt = h.db.Model(result)
	schema, err = gormSchema.Parse(result, &sync.Map{}, gormSchema.NamingStrategy{})

	return
}

func (h *QueryHandler) buildResult(model interface{}) interface{} {
	value := reflect.ValueOf(model)

	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	return reflect.New(value.Type()).Interface()
}

func buildPreloads(stmt *gorm.DB, pre []Preload) {
	for _, p := range pre {
		stmt.Preload(p.Name, p.Conditions...)
	}
}

func buildSelection(stmt *gorm.DB, query SlicedQueryObject, schema *gormSchema.Schema) {
	sel := query.Slice()
	if sel == nil {
		return
	}

	var total int64

	stmt.Session(&gorm.Session{})
	stmt.Count(&total)
	sel.Total = total

	if sel.Offset > 0 {
		stmt.Offset(sel.Offset)
	}
	if sel.Limit > 0 {
		stmt.Limit(sel.Limit)
	}
	if len(sel.Sort) > 0 {
		err := buildSort(stmt, sel.Sort, schema)
		if err != nil {
			query.SetError(err)
		}
	}

	return
}

func buildFilter(stmt *gorm.DB, query FilteredQueryObject, schema *gormSchema.Schema) {
	filters := query.Filters()
	if len(filters) == 0 {
		return
	}

	fields := schema.FieldsByName
	for _, filter := range filters {
		field := fields[filter.FieldName]
		if field == nil {
			query.SetError(errors.New("unknown field name"))
			return
		}

		stmt.Where(fmt.Sprintf("%s %s ?", field.DBName, filter.Operator), filter.Value)
	}

	return
}

func buildSort(stmt *gorm.DB, sorts []Sort, schema *gormSchema.Schema) error {
	fields := schema.FieldsByName
	for _, sort := range sorts {
		field := fields[sort.FieldName]
		if field == nil {
			return errors.New("unknown field name")
		}
		stmt.Order(clause.OrderByColumn{Column: clause.Column{Name: field.DBName}, Desc: sort.Direction == DESC})
	}

	return nil
}
