package swagger

import (
	"reflect"
)

type InType string

const (
	InBody   InType = "body"
	InQuery  InType = "query"
	InPath   InType = "path"
	InHeader InType = "header"
	InCookie InType = "cookie"
)

type Parameter struct {
	*BaseObject
	// How passed parameter - in body, in query or in path
	IN InType `json:"in,omitempty"`
	// Is a required parameter?
	Req bool `json:"required,omitempty"`
}

func NewParameter(name, description string, t interface{}, required bool, inType InType) *Parameter {
	return &Parameter{
		BaseObject: NewBaseObject(name, description, t),
		Req:        required,
		IN:         inType,
	}
}

// Parse a parameter structure for JSON generation
func (p *Parameter) Parse(sw Doc) {
	ParseRootType(p, sw)
}

func (p *Parameter) GetSchema() *Schema {
	return p.Schema
}

func (p *Parameter) SetTypeName(typeName string) {
	p.TypeName = typeName
}

func (p *Parameter) SetFormat(format string) {
	p.Format = format
}

func (p *Parameter) GetType() reflect.Kind {
	return p.Type
}

type ArrayParameters []*Parameter
