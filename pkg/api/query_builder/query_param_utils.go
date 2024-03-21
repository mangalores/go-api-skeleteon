package query_builder

import (
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"

	"github.com/mangalores/go-api-skeleton/pkg/db"
)

const sortParamPattern = "^(?P<name>[a-zA-Z0-9]+):(?P<direction>([a-zA-Z0-9]+))$"

var sortRegEx = regexp.MustCompile(sortParamPattern)

// extractSortParamValues
// provided query param _sort=fieldA:asc&_sort=fieldB:desc => []Sort{{fieldName:"fieldA",direction:"ASC"},{fieldName:"fieldB",direction:"DESC"}}
func extractSortParamValues(params url.Values, fields map[string]string) ([]db.Sort, error) {
	sortFields := make([]db.Sort, 0)

	if list := params[sortField]; len(list) > 0 {
		for _, v := range list {
			if !sortRegEx.MatchString(v) {
				return sortFields, NewInvalidParamValueErr(sortField, false)
			}

			matches := sortRegEx.FindAllStringSubmatch(v, -1)
			name := matches[0][1]
			fieldName, ok := fields[name]
			if !ok {
				return sortFields, NewInvalidParamValueErr(sortField, false)
			}

			direction, ok := mapDirection(matches[0][2])
			if !ok {
				return sortFields, NewInvalidParamValueErr(sortField, false)
			}

			sortFields = append(sortFields, db.Sort{FieldName: fieldName, Direction: direction})
		}
	}

	return sortFields, nil
}

func mapDirection(v string) (direction db.Direction, ok bool) {
	switch v {
	case string(db.ASC):
		return db.ASC, true
	case string(db.DESC):
		return db.DESC, true
	default:
		return "", false
	}
}
func extractNumericParamValue(params url.Values, name string) (value int, err error) {
	if l, ok := params[name]; ok {
		value, err = strconv.Atoi(l[0])
		if err != nil {
			return value, NewInvalidParamValueErr(name, true)
		}
	}

	return value, err
}

func stripReservedFields(values url.Values) url.Values {
	reserved := regexp.MustCompile(fmt.Sprintf("^%s[a-zA-Z0-9]+$", reservedPrefix))
	params := make(url.Values)
	for k, v := range values {
		if reserved.MatchString(k) {
			continue
		}

		params[k] = v
	}

	return params
}

func acceptedFilterFields(i interface{}) map[string]string {
	fields := make(map[string]string)
	t := getType(i)

	if t.Kind() != reflect.Struct {
		return fields
	}
	for _, field := range reflect.VisibleFields(t) {
		tag := field.Tag.Get("json")

		if tag != "" {
			fieldName := field.Name
			fields[tag] = fieldName
		}
	}

	return fields
}

func getType(i interface{}) reflect.Type {
	return traverseType(reflect.TypeOf(i))
}

func traverseType(typeOf reflect.Type) reflect.Type {
	switch typeOf.Kind() {
	case reflect.Array, reflect.Slice, reflect.Pointer, reflect.UnsafePointer:
		return traverseType(typeOf.Elem())
	default:
		return typeOf
	}
}
