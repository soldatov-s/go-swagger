package swagger

import (
	"reflect"
	"strings"
)

type TypeDictElement struct {
	TypeName string
	Format   string
}

var (
	// Dictinary with self-types
	typeDict = map[string]TypeDictElement{
		"Time": {
			TypeName: "string",
			Format:   "date-time",
		},
		"NullTime": {
			TypeName: "string",
			Format:   "date-time",
		},
		"UUID": {
			TypeName: "string",
			Format:   "uuid",
		},
		"NullString": {
			TypeName: "string",
			Format:   "",
		},
		"NullMeta": {
			TypeName: "object",
			Format:   "",
		},
	}
)

type Definition struct {
	// Type name
	TypeName string `json:"type,omitempty"`
	// List of properties of definition object
	Properties map[string]*Property `json:"properties,omitempty"`
}

// AddNewDefinition is a helper for add new definition in map
func AddNewDefinition(objName string, s interface{}, sw *Doc) {
	if _, ok := sw.Definitions[objName]; !ok {
		sw.Definitions[objName] = &Definition{}
		sw.Definitions[objName].Parse(s, sw)
	}
}

func parseInterfaceOrStruct(obj interface{}, sw *Doc) (ref string) {
	o := valueFromPtr(obj)
	name := reflect.TypeOf(o).Name()
	ref = "#/definitions/" + name
	AddNewDefinition(name, o, sw)
	return ref
}

func parseArrayOrSlice(obj interface{}, sw *Doc) *BaseObject {
	var itemRef string

	o := reflect.New(reflect.TypeOf(obj).Elem()).Interface()

	o = valueFromPtr(o)
	if o == nil {
		return buildBaseObject(reflect.Interface, itemRef)
	}

	switch reflect.TypeOf(o).Kind() {
	case reflect.Interface, reflect.Struct:
		itemRef = parseInterfaceOrStruct(o, sw)
	}

	return buildBaseObject(reflect.TypeOf(o).Kind(), itemRef)
}

func parseMap(obj interface{}, sw *Doc) *AdditionalProperties {
	var addPropRef string

	o := reflect.New(reflect.TypeOf(obj).Elem()).Interface()

	o = valueFromPtr(o)
	if o == nil {
		return buildAdditionalProperties(reflect.Interface, addPropRef)
	}

	switch reflect.TypeOf(o).Kind() {
	case reflect.Interface, reflect.Struct:
		addPropRef = parseInterfaceOrStruct(o, sw)
	}
	return buildAdditionalProperties(reflect.TypeOf(o).Kind(), addPropRef)
}

func buildBaseObject(kind reflect.Kind, ref string) *BaseObject {
	tp, fmt := ParseKind(kind)
	return &BaseObject{
		TypeName: tp,
		Ref:      ref,
		Format:   fmt,
	}
}

func buildAdditionalProperties(kind reflect.Kind, ref string) *AdditionalProperties {
	tp, fmt := ParseKind(kind)
	return &AdditionalProperties{
		Ref:      ref,
		TypeName: tp,
		Format:   fmt,
	}
}

func valueFromPtr(s interface{}) interface{} {
	if s == nil {
		return nil
	}

	if reflect.TypeOf(s).Kind() == reflect.Ptr {
		if !reflect.ValueOf(s).IsNil() {
			return valueFromPtr(reflect.ValueOf(s).Elem().Interface())
		}
		return valueFromPtr(reflect.New(reflect.TypeOf(s).Elem()).Interface())
	}

	return s
}

// Parse definition of object
func (d *Definition) Parse(s interface{}, sw *Doc) {
	if d.TypeName == "" {
		d.TypeName = constObject
	}

	if d.Properties == nil {
		d.Properties = make(MapProperty)
	}

	Types := reflect.TypeOf(s)
	Values := reflect.ValueOf(s)

	// Walk through all the fields of the structure
	for j := 0; j < Types.NumField(); j++ {
		var name string

		tag, ok := Types.Field(j).Tag.Lookup("json")
		if !ok {
			name = Types.Field(j).Name
		} else {
			name = strings.Split(tag, ",")[0]
		}

		// skip hided fields
		if name == "-" {
			continue
		}

		// Parse swagtype for overriding
		swagType, _ := Types.Field(j).Tag.Lookup("swagtype")

		// Parse enumeration
		var swagEnumArray []interface{}
		if swagEnum, ok := Types.Field(j).Tag.Lookup("swagenum"); ok {
			enumArray := strings.Split(swagEnum, ",")
			if len(enumArray) > 0 {
				swagEnumArray = make([]interface{}, 0, len(enumArray))
				for _, item := range enumArray {
					swagEnumArray = append(swagEnumArray, item)
				}
			}
		}

		if Types.Field(j).Anonymous {
			d.Parse(Values.Field(j).Interface(), sw)
		} else {
			d.Properties[name] = parseStructField(Types.Field(j).Type, Values.Field(j), sw, swagType, swagEnumArray)
		}
	}
}

func parseStructField(tp reflect.Type, val reflect.Value, sw *Doc, swagType string, swagEnum []interface{}) *Property {
	var (
		addProp  *AdditionalProperties
		item     *BaseObject
		typeName string
		ref      string
		format   string
	)

	obj := typeDict[tp.Name()]
	if obj.TypeName != "" {
		typeName = obj.TypeName
		format = obj.Format
	} else {
		switch tp.Kind() {
		case reflect.Interface:
			typeName = constObject
			if !val.IsNil() {
				return parseStructField(reflect.TypeOf(val.Interface()), val, sw, swagType, swagEnum)
			}
		case reflect.Ptr:
			return parseStructField(tp.Elem(), reflect.New(tp.Elem()), sw, swagType, swagEnum)
		case reflect.Struct:
			typeName = constObject
			ref = parseInterfaceOrStruct(val.Interface(), sw)
		case reflect.Array, reflect.Slice:
			typeName = constArray
			item = parseArrayOrSlice(val.Interface(), sw)
		case reflect.Map:
			typeName = constObject
			addProp = parseMap(val.Interface(), sw)
		default:
			typeName, format = ParseKind(tp.Kind())
			if swagType != "" {
				typeName = swagType
			}
		}
	}

	if ref != "" {
		typeName = ""
	}

	return &Property{
		BaseObject: BaseObject{
			TypeName:             typeName,
			Ref:                  ref,
			Format:               format,
			Enum:                 swagEnum,
			AdditionalProperties: addProp,
		},
		Item: item,
	}
}
