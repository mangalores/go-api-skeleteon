package utils

import (
	"reflect"
)

func Contains(elems []string, elem string) bool {
	for _, v := range elems {
		if v == elem {
			return true
		}
	}

	return false
}

func StripPointer(i interface{}) interface{} {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
		return StripPointer(v.Interface())
	}

	return v.Interface()
}
