package swagger

import (
	"fmt"
	"reflect"
)

const (
	constInteger = "integer"
	constNumber  = "number"
	constBoolean = "boolean"
	constObject  = "object"
	constArray   = "array"
)

type AdditionalProperties struct {
	// Object format, for example: int64 is integer with format int64
	Format string `json:"format,omitempty"`
	// Reference to a schema definition
	Ref string `json:"$ref,omitempty"`
	// Type name
	TypeName string `json:"type,omitempty"`
}

type Schema struct {
	// Object format, for example: int64 is integer with format int64
	Format string `json:"format,omitempty"`
	// Reference to a schema definition
	Ref string `json:"$ref,omitempty"`
	// Type name
	TypeName string `json:"type,omitempty"`
	//
	AdditionalProperties *AdditionalProperties `json:"additionalProperties,omitempty"`
	//
	Item *BaseObject `json:"items,omitempty"`
	//
	Type interface{} `json:"-"`
}

func NewSchema(schema interface{}) *Schema {
	return &Schema{
		Type: schema,
	}
}

type Schemater interface {
	GetSchema() *Schema
	SetTypeName(string)
	GetType() reflect.Kind
	SetFormat(string)
}

// ParseRootType is a method for analyzing Types and Schemas of Parameters and
// Response
func ParseRootType(obj Schemater, sw Doc) {
	// Parse Type when Schema unused
	if obj.GetSchema() == nil {
		TypeName, Format := ParseKind(obj.GetType())
		obj.SetTypeName(TypeName)
		obj.SetFormat(Format)
		return
	}

	value := valueFromPtr(obj.GetSchema().Type)

	// Parse Schema, when it is Structure or Pointer to structure
	switch reflect.ValueOf(value).Kind() {
	case reflect.Struct:
		obj.GetSchema().TypeName = ""
		obj.GetSchema().Ref = parseInterfaceOrStruct(value, &sw)
		return
	case reflect.Slice, reflect.Array:
		obj.GetSchema().TypeName = constArray
		obj.GetSchema().Item = parseArrayOrSlice(value, &sw)
		return
	case reflect.Map:
		obj.GetSchema().TypeName = constObject
		obj.GetSchema().AdditionalProperties = parseMap(value, &sw)
		return
	}

	// Parse Schema when it is reflect.Kind type
	TypeName, Format := ParseKind(value)
	obj.SetTypeName(TypeName)
	obj.SetFormat(Format)
}

// ParseKind parse simple types and return it TypeName and Format
func ParseKind(kind interface{}) (typeName, format string) {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Uint, reflect.Uint8, reflect.Uint16:
		typeName = constInteger
	case reflect.Int32, reflect.Int64, reflect.Uint32, reflect.Uint64:
		typeName = constInteger
		format = fmt.Sprint(kind)
	case reflect.Float32, reflect.Float64:
		typeName = constNumber
	case reflect.Bool:
		typeName = constBoolean
	case reflect.Invalid, reflect.Struct:
		// No changes if type not setted
	case reflect.Interface, reflect.Map:
		typeName = constObject
	default:
		typeName = fmt.Sprint(kind)
	}
	return
}
