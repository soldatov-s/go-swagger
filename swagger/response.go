package swagger

import (
	"reflect"
	"strconv"
)

type Response struct {
	*BaseObject
}

func NewResponse(description string) *Response {
	return &Response{
		BaseObject: NewBaseObject("", description, nil),
	}
}

// Parse a response structure for JSON generation
func (r *Response) Parse(sw Doc) {
	ParseRootType(r, sw)
}

func (r *Response) GetSchema() *Schema {
	return r.Schema
}

func (r *Response) SetTypeName(typeName string) {
	if r.TypeName == "" {
		r.TypeName = typeName
	}
}

func (r *Response) GetType() reflect.Kind {
	return r.Type
}

func (r *Response) SetFormat(format string) {
	if r.Format == "" {
		r.Format = format
	}
}

type MapResponse map[string]*Response

type MapResponseCode map[int]int

func (arc MapResponseCode) Append(respCode int) string {
	if arc == nil {
		arc = make(MapResponseCode)
	}

	if _, ok := arc[respCode]; !ok {
		arc[respCode] = 1
		return strconv.Itoa(respCode)
	}

	arc[respCode]++
	return strconv.Itoa(respCode) + "(" + strconv.Itoa(arc[respCode]) + ")"
}
