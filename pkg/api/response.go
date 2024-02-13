package api

import (
	"fmt"
	"github.com/mangalores/go-api-skeleton/pkg/utils"
	"reflect"

	"github.com/labstack/echo/v4"
	"github.com/mangalores/go-api-skeleton/pkg/db"
)

type MappingFunc func(e interface{}) (interface{}, error)
type PathFunc func(ctx echo.Context, e interface{}) (string, error)

type ResponseHandler struct {
	mappings map[string]MappingFunc
}

func NewResponseHandler() *ResponseHandler {
	return &ResponseHandler{
		make(map[string]MappingFunc),
	}
}

func (r *ResponseHandler) Register(i interface{}, m MappingFunc) {
	r.mappings[reflect.TypeOf(i).String()] = m
}

func (r *ResponseHandler) Handle(ctx echo.Context, query db.QueryObject) interface{} {
	switch query.(type) {
	case db.SlicedQueryObject:
		return r.mapCollection(ctx, query.(db.SlicedQueryObject))
	default:
		return r.mapSimpleQuery(ctx, query.(db.QueryObject))
	}

}

func (r *ResponseHandler) mapSimpleQuery(ctx echo.Context, query db.QueryObject) interface{} {
	println(fmt.Sprintf("%T", query.Result()))
	result := utils.StripPointer(query.Result())
	mapFunc, err := r.getMap(result)
	if err != nil {
		query.SetError(err)
		return nil
	}
	if mapFunc == nil {
		query.SetError(fmt.Errorf("mapping not found for %T", result))
		return nil
	}

	mapped, err := mapFunc(result)
	if err != nil {
		query.SetError(err)
		return nil
	}

	return mapped
}

func (r *ResponseHandler) mapCollection(ctx echo.Context, query db.SlicedQueryObject) interface{} {
	result := utils.StripPointer(query.Result())
	m, err := r.getMap(result)
	if err != nil {
		query.SetError(err)
		return nil
	}

	return NewCollection(ctx, query, m)
}

func (r *ResponseHandler) getMap(i interface{}) (MappingFunc, error) {
	if i == nil {
		return nil, nil
	}

	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	n := t.String()

	m := r.mappings[n]
	if m == nil {
		return nil, fmt.Errorf("mapping not found for %s", n)
	}

	return m, nil
}

func GenerateSelfLink(path string) LinkOpts {
	links := make(Links, 0)

	links["self"] = Link{
		Href: path,
	}

	return LinkOpts{Links: links}
}

func NewCollection(ctx echo.Context, query db.SlicedQueryObject, mapFunc MappingFunc) *Collection {
	result := utils.StripPointer(query.Result())
	items, err := mapFunc(result)
	if err != nil {
		query.SetError(err)
		return nil
	}

	collection := Collection{
		Embedded: Embedded{
			"items": items,
		},
	}

	if sel := query.Slice(); sel != nil {
		collection.LinkOpts = GenerateCollectionLinks(sel, ctx.Request().URL.Path)
		collection.Metadata = CollectionOpts{
			Offset: sel.Offset,
			Limit:  sel.Limit,
			Total:  sel.Total,
		}
	}

	return &collection
}

func GenerateCollectionLinks(s *db.Slice, path string) LinkOpts {
	links := make(Links)

	links["self"] = Link{
		fmt.Sprintf("%s?offset=%d&limit=%d", path, s.Offset, s.Limit),
	}

	firstOffset := 0
	if firstOffset != s.Offset {
		links["first"] = Link{
			fmt.Sprintf("%s?offset=%d&limit=%d", path, firstOffset, s.Limit),
		}
	}

	lastOffset := s.Total - s.Total%int64(s.Limit)
	if lastOffset > int64(s.Offset) {
		links["last"] = Link{
			fmt.Sprintf("%s?offset=%d&limit=%d", path, lastOffset, s.Limit),
		}
	}

	prevOffset := s.Offset - s.Limit
	if prevOffset > 0 {
		links["prev"] = Link{
			fmt.Sprintf("%s?offset=%d&limit=%d", path, prevOffset, s.Limit),
		}
	}

	nextOffset := s.Offset + s.Limit
	if s.Total > int64(nextOffset) {
		links["next"] = Link{
			fmt.Sprintf("%s?offset=%d&limit=%d", path, nextOffset, s.Limit),
		}
	}

	return LinkOpts{links}
}
