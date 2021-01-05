package swagger

import "reflect"

// Base object include all common fields of responses and parameters
type BaseObject struct {
	// Type name
	TypeName string `json:"type,omitempty"`
	// List of enumeration values
	Enum []interface{} `json:"enum,omitempty"`
	// Default value
	Default interface{} `json:"default,omitempty"`
	// Example of value
	Example interface{} `json:"example,omitempty"`
	// Could it be null?
	Nullable string `json:"nullable,omitempty"`
	// Detailed object description
	Description string `json:"description,omitempty"`
	// Name of object
	Name string `json:"name,omitempty"`
	// Object format, for example: int64 is integer with format int64
	Format string `json:"format,omitempty"`
	// Reference to a schema definition
	Ref string `json:"$ref,omitempty"`
	// The object schema
	Schema *Schema `json:"schema,omitempty"`
	// Kind of object
	Type reflect.Kind `json:"-"`
	//
	AdditionalProperties *AdditionalProperties `json:"additionalProperties,omitempty"`
}

func NewBaseObject(name, description string, t interface{}) *BaseObject {
	return &BaseObject{
		Name:        name,
		Description: description,
		Schema:      NewSchema(t),
	}
}
