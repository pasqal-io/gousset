package schema

import (
	"fmt"
	"reflect"

	tags "github.com/pasqal-io/gousset/inner"
	"github.com/pasqal-io/gousset/openapi/doc"
	"github.com/pasqal-io/gousset/openapi/example"
	"github.com/pasqal-io/gousset/openapi/shared"
)

type Schema interface {
	sealed()
}

// Use this interface to customize the schema returned by
// your type.
type HasSchema interface {
	Schema() Schema
}

type None struct{}

func (None) sealed() {}

type Shared struct {
	ExternalDocs *doc.External `json:"externalDocs,omitempty"`
	Example      *shared.Json  `json:"example,omitempty"`
	Type         Type          `json:"type"`
}

type Type string

const (
	TypeString = Type("string")
	TypeObject = Type("object")
	TypeArray  = Type("array")
	TypeNumber = Type("number")
	TypeBool   = Type("bool")
)

type Primitive struct {
	Shared
	// A well-known format, e.g. "email".
	Format *string `json:"format,omitempty"`
}

func (Primitive) sealed() {}

type Object struct {
	Shared
	// Required fields.
	Required   *[]string         `json:"required,omitempty"`
	Properties map[string]Schema `json:"properties"`
}

func (Object) sealed() {}

type Array struct {
	Shared
	Items Schema             `json:"items"`
	Defs  *map[string]Schema `json:"$defs,omitempty"`
}

func (Array) sealed() {}

type OneOf struct {
	OneOf []Schema `json:"oneOf"`
}

func (OneOf) sealed() {}

type AllOf struct {
	AllOf []Schema `json:"allOf"`
}

func (AllOf) sealed() {}

type Discriminator struct {
	// The name of the property in the payload that will hold the discriminating value. This property SHOULD be required in the payload schema, as the behavior when the property is absent is undefined.
	PropertyName string `json:"propertyName"`

	// An object to hold mappings between payload values and schema names or URI references.
	Mapping *map[string]string `json:"mapping"`
}

// Create a schema from a type.
//
// As of this writing, we make no attempt to optimize schemas if e.g. some data structures are repeated.
//
// Arguments:
//
//	typ The type to extract. See HasSchema, HasExternalDocs, HasExample for means to configure it.
func FromType(typ reflect.Type, publicNameKey string) (Schema, error) {
	var externalDocs *doc.External
	var anExample *shared.Json
	phony := reflect.New(typ)
	if phony.CanInterface() {
		asAny := phony.Interface()
		// Give priority to a schema provided by the user.
		if hasSchema, ok := asAny.(HasSchema); ok {
			return hasSchema.Schema(), nil
		}
		if hasExternalDocs, ok := asAny.(doc.HasExternalDocs); ok {
			externalDocs = shared.Ptr(hasExternalDocs.Docs())
		}
		if hasExample, ok := asAny.(example.HasExample); ok {
			anExample = shared.Ptr(hasExample.Example())
		}
	}

	switch typ.Kind() {
	case reflect.Pointer:
		return FromType(typ.Elem(), publicNameKey)
	case reflect.Bool:
		return Primitive{
			Shared: Shared{
				Type:         TypeBool,
				ExternalDocs: externalDocs,
				Example:      anExample,
			},
			Format: nil,
		}, nil
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		fallthrough
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		return Primitive{
			Shared: Shared{
				Type:         TypeNumber,
				ExternalDocs: externalDocs,
				Example:      anExample,
			},
			Format: nil,
		}, nil
	case reflect.String:
		return Primitive{
			Shared: Shared{
				Type:         TypeString,
				ExternalDocs: externalDocs,
				Example:      anExample,
			},
			Format: nil,
		}, nil
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		items, err := FromType(typ.Elem(), publicNameKey)
		if err != nil {
			return None{}, fmt.Errorf("failed to extract type from the elements of array/slice %s", typ.String())
		}
		return Array{
			Shared: Shared{
				Type:         TypeArray,
				ExternalDocs: externalDocs,
				Example:      anExample,
			},
			Items: items,
			Defs:  nil,
		}, nil
	case reflect.Struct:
		var required []string
		properties := make(map[string]Schema)
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			tags, err := tags.Parse(field.Tag)
			if err != nil {
				return None{}, fmt.Errorf("failed to parse tags for field %s of struct %s", field.Name, typ.String())
			}

			name := field.Name
			if publicName := tags.PublicFieldName(publicNameKey); publicName != nil {
				name = *publicName
			}

			fieldSchema, err := FromType(field.Type, publicNameKey)
			if err != nil {
				return None{}, fmt.Errorf("failed to extract scheme from field %s of struct %s", field.Name, typ.String())
			}

			if tags.IsFlattened() {
				panic("not implemented")
			}

			if tags.Default() == nil && !tags.IsPreinitialized() && tags.MethodName() == nil {
				required = append(required, name)
			}
			properties[name] = fieldSchema
		}
		return Object{
			Shared: Shared{
				Type:         TypeArray,
				ExternalDocs: externalDocs,
				Example:      anExample,
			},
			Required:   &required,
			Properties: properties,
		}, nil
	default:
		return None{}, fmt.Errorf("couldn't find any scheme for type %s, you may need to implement HasSchema", typ.String())
	}
}
