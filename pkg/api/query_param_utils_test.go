package api

import (
	"github.com/mangalores/go-api-skeleton/pkg/db"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

type MockEntity struct {
	Foo string `json:"foo"`
	Bar string `json:"bar"`
	Baz string
}

func TestExtractSortParamValues(t *testing.T) {
	acceptedFields := map[string]string{"foo": "Foo", "bar": "Bar"}
	testCases := [][]interface{}{
		{
			// proper multi field sort
			url.Values{"_sort": {"foo:asc", "bar:desc"}},
			[]db.Sort{{FieldName: "Foo", Direction: db.ASC}, {FieldName: "Bar", Direction: db.DESC}},
			nil,
		},
		{
			// attempt to sort with invalid direction value
			url.Values{"_sort": {"foo:asc", "bar:desc", "bar:baz"}},
			[]db.Sort{{FieldName: "Foo", Direction: db.ASC}, {FieldName: "Bar", Direction: db.DESC}},
			InvalidParamValueErr{name: "_sort", numeric: false},
		},
		{
			// attempt to sort by private field
			url.Values{"_sort": {"foo:asc", "bar:desc", "baz:desc"}},
			[]db.Sort{{FieldName: "Foo", Direction: db.ASC}, {FieldName: "Bar", Direction: db.DESC}},
			InvalidParamValueErr{name: "_sort", numeric: false},
		},
		{
			// attempt to sort with unknown field for entity type
			url.Values{"_sort": {"foo:asc", "bar:desc", "boo:desc"}},
			[]db.Sort{{FieldName: "Foo", Direction: db.ASC}, {FieldName: "Bar", Direction: db.DESC}},
			InvalidParamValueErr{name: "_sort", numeric: false},
		},
		{
			// ensure sort order
			url.Values{"_sort": {"bar:desc", "foo:asc"}},
			[]db.Sort{{FieldName: "Bar", Direction: db.DESC}, {FieldName: "Foo", Direction: db.ASC}},
			nil,
		},
	}

	for _, testCase := range testCases {
		test := testCase[0].(url.Values)
		expectedSort := testCase[1]
		expectedErr := testCase[2]

		sort, err := extractSortParamValues(test, acceptedFields)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, expectedSort, sort)
	}
}

func Test_GetType(t *testing.T) {
	getType(struct{}{})
	getType(&struct{}{})
	getType([]struct{}{})
	getType(&[]struct{}{})
	getType(&[]*struct{}{})

}
