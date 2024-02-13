package api

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/mangalores/go-api-skeleton/pkg/db"
	"github.com/mangalores/go-api-skeleton/pkg/utils"
	"net/url"
	"reflect"
	"regexp"
)

const compositePattern = "^(?P<name>[a-zA-Z0-9]+):(?P<op>[a-zA-Z]+)$"
const plainParamPattern = "^[a-zA-Z0-9]+$"
const (
	maxLimit       = 10000
	defaultLimit   = 500
	defaultOffset  = 0
	sortField      = "_sort"
	limitField     = "_limit"
	offsetField    = "_offset"
	embedField     = "_embed"
	reservedPrefix = "_"
)

var plainRegEx = regexp.MustCompile(plainParamPattern)
var compositeRegEx = regexp.MustCompile(compositePattern)

var operatorMap = map[string]string{
	"eq":  "=",
	"ne":  "<>",
	"lt":  "<",
	"gt":  ">",
	"lte": "<=",
	"gte": ">=",
}

type Operator struct {
	operator    string
	transformFn func(value []string) []string
	matchField  *regexp.Regexp
}

func (o *Operator) Operator() string {
	return o.operator
}

func (o *Operator) FieldMatches(fieldName string) bool {
	return o.matchField.MatchString(fieldName)
}

func (o *Operator) TransformValue(values []string) interface{} {
	return o.transformFn(values)
}

type QueryBuilder struct {
	operators    map[string]*Operator
	allowedEmbed map[string]db.Preload
	isSlice      bool
	defaultSort  []db.Sort
	presetFilter []db.Filter
	preload      []db.Preload
	model        interface{}

	loadQueryParams bool
	loadPathParams  bool
	loadPreloads    bool
}

func NewQueryBuilder(model interface{}) *QueryBuilder {
	builder := &QueryBuilder{}
	builder.Reset(model)

	return builder
}

func (b *QueryBuilder) Reset(model interface{}) *QueryBuilder {
	b.model = model
	b.operators = make(map[string]*Operator, 10)
	b.registerDefaultOperators()
	b.defaultSort = []db.Sort{}
	b.presetFilter = []db.Filter{}
	b.preload = []db.Preload{}

	b.loadQueryParams = true
	b.loadPathParams = true
	b.loadPreloads = false

	return b
}

func (b *QueryBuilder) registerDefaultOperators() {
	for identifier, operator := range operatorMap {
		b.RegisterOperator(identifier, operator, defaultTransformFn, regexp.MustCompile(".*"))
	}
}

func (b *QueryBuilder) RegisterOperator(identifier string, operator string, fn func(value []string) []string, rx *regexp.Regexp) {
	b.operators[identifier] = &Operator{
		operator,
		fn,
		rx,
	}
}

func (b *QueryBuilder) SetSlice(flag bool) {
	b.isSlice = flag
}

func (b *QueryBuilder) SetParsePathParams(flag bool) {
	b.loadPathParams = flag
}

func (b *QueryBuilder) SetParseQueryParams(flag bool) {
	b.loadQueryParams = flag
}

func (b *QueryBuilder) SetParseEmbedding(flag bool) {
	b.loadPreloads = flag
}

func (b *QueryBuilder) AddDefaultSort(fieldName string, direction db.Direction) {
	b.defaultSort = append(b.defaultSort, db.Sort{FieldName: fieldName, Direction: direction})
}

func (b *QueryBuilder) AddPresetFilter(fieldName string, operator string, value interface{}) error {
	var (
		op *Operator
		ok bool
	)

	if op, ok = b.mapOperator(operator); !ok {
		return errors.New("invalid operator id")
	}
	if ok = op.FieldMatches(fieldName); !ok {
		return errors.New("invalid field name")
	}

	b.presetFilter = append(b.presetFilter, db.Filter{FieldName: fieldName, Operator: op.Operator(), Value: value})

	return nil
}

func (b *QueryBuilder) AddPreload(preloads ...db.Preload) {
	for _, preload := range preloads {
		b.preload = append(b.preload, preload)
	}
}

func (b *QueryBuilder) AddAllowedEmbed(name string, associationName string) {
	b.allowedEmbed[name] = db.Preload{Name: associationName}
}

func (b *QueryBuilder) Build(ctx echo.Context) (query db.QueryObject, err error) {
	simpleQuery, err := b.buildQuery(ctx)
	if err != nil {
		return
	}
	query = simpleQuery

	filteredQuery, err := b.buildFilterQuery(ctx, simpleQuery)
	if err != nil || !b.isSlice {
		query = filteredQuery
		return
	}

	query, err = b.buildSlicedQuery(ctx, filteredQuery)
	if err != nil {
		return
	}

	return
}

func (b *QueryBuilder) buildQuery(ctx echo.Context) (*db.Query, error) {
	q := &db.Query{}

	if b.model == nil {
		return q, errors.New("model struct/slice must be set")
	}
	q.SetModel(b.model)

	preloads, err := b.buildPreloads(ctx)
	if len(preloads) > 0 {
		q.SetPreloads(&preloads)
	}

	return q, err
}

func (b *QueryBuilder) buildFilterQuery(ctx echo.Context, query *db.Query) (*db.FilterQuery, error) {
	filters := make([]db.Filter, 0)
	filterQuery := &db.FilterQuery{Query: *query}

	paramFilters, err := b.buildFilters(b.fetchParameters(ctx), acceptedFilterFields(b.model))
	if err != nil {
		return filterQuery, err
	}
	filters = append(filters, paramFilters...)
	filters = append(filters, b.presetFilter...)

	filterQuery.SetFilters(filters)

	return filterQuery, err
}

func (b *QueryBuilder) buildSlicedQuery(ctx echo.Context, query *db.FilterQuery) (*db.CollectionQuery, error) {
	cq := &db.CollectionQuery{FilterQuery: *query}

	if reflect.TypeOf(utils.StripPointer(b.model)).Kind() != reflect.Slice {
		return cq, errors.New("for collections specified destination object must be of type slice")
	}
	slice, err := b.buildSlice(ctx)
	if err != nil {
		return cq, err
	}

	slice.Sort, err = extractSortParamValues(ctx.QueryParams(), acceptedFilterFields(b.model))
	if err != nil {
		return cq, err
	}
	if len(slice.Sort) == 0 && len(b.defaultSort) > 0 {
		slice.Sort = b.defaultSort
	}

	cq.SetSlice(slice)

	return cq, err
}

func (b *QueryBuilder) fetchParameters(ctx echo.Context) url.Values {
	params := make(url.Values)

	if b.loadQueryParams {
		for name, value := range ctx.QueryParams() {
			params[name] = value
		}
	}

	// path parameter take precedence
	if b.loadPathParams {
		pathNames := ctx.ParamNames()
		pathValues := ctx.ParamValues()
		for i, n := range pathNames {
			params[n] = []string{pathValues[i]}
		}
	}

	return params
}

func (b *QueryBuilder) buildFilters(params url.Values, fields map[string]string) ([]db.Filter, error) {
	filters := make([]db.Filter, 0)

	params = stripReservedFields(params)

	for key, value := range params {
		var (
			ok        bool
			fieldName string
			operator  *Operator
		)

		filter := db.Filter{}

		fieldName, operator, ok = b.extractParamAndOperator(key)
		if !ok {
			return filters, NewInvalidFilterErr(filter)
		}

		// todo allow multiple as or list
		if len(value) > 1 {
			return filters, NewInvalidMultipleValuesErr(filter)
		}

		// skip fields not in allowed list
		_, ok = fields[fieldName]
		if !ok || !operator.FieldMatches(fieldName) {
			continue
		}

		filter.FieldName = fields[fieldName]

		// apply transformation function to values
		filter.Value = operator.TransformValue(value)
		filter.Operator = operator.Operator()

		filters = append(filters, filter)
	}

	return filters, nil
}

func (b *QueryBuilder) buildPreloads(ctx echo.Context) ([]db.Preload, error) {
	params := ctx.QueryParams()
	if list := params[embedField]; len(list) > 0 {
		for _, name := range list {
			preload, ok := b.allowedEmbed[name]
			if !ok {
				return b.preload, NewInvalidEmbedErr(name)
			}

			b.preload = append(b.preload, preload)
		}
	}

	return b.preload, nil
}

func (b *QueryBuilder) buildSlice(ctx echo.Context) (slice *db.Slice, err error) {
	offset := defaultOffset
	limit := defaultLimit

	offset, err = extractNumericParamValue(ctx, offsetField)
	if err != nil {
		return &db.Slice{Offset: defaultOffset, Limit: defaultLimit}, NewInvalidParamValueErr(offsetField, true)
	}

	limit, err = extractNumericParamValue(ctx, limitField)
	if err != nil {
		return &db.Slice{Offset: defaultOffset, Limit: defaultLimit}, NewInvalidParamValueErr(limitField, true)
	}

	if limit > maxLimit {
		return &db.Slice{Offset: defaultOffset, Limit: defaultLimit}, NewMaxLimitExceededErr()
	}
	if limit <= 0 {
		limit = defaultLimit
	}

	return &db.Slice{Offset: offset, Limit: limit}, nil
}

func (b *QueryBuilder) extractParamAndOperator(param string) (string, *Operator, bool) {

	if plainRegEx.MatchString(param) {
		param = param + ":eq"
	}

	if compositeRegEx.MatchString(param) {
		matches := compositeRegEx.FindAllStringSubmatch(param, -1)
		name := matches[0][1]
		operator, ok := b.mapOperator(matches[0][2])

		return name, operator, ok
	}

	return "", nil, false
}

func (b *QueryBuilder) mapOperator(v string) (o *Operator, ok bool) {

	if o = b.operators[v]; o != nil {
		return o, true
	}

	return nil, false
}

func defaultTransformFn(value []string) []string {
	return value
}

func SearchTransformFN(value []string) []string {
	transformed := make([]string, 0)
	for _, v := range value {
		transformed = append(transformed, v+"%")
	}
	return transformed
}
