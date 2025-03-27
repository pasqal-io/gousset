// Definitions of JSON-level data structures.
//
// Any data structure that the API expects or provides, regardless
// of whether it's through body, path, query, headers, is defined
// as a `schema.Schema`.
package schema

import (
	"fmt"
	"reflect"

	"github.com/pasqal-io/gousset/inner/tags"
	"github.com/pasqal-io/gousset/openapi/doc"
	"github.com/pasqal-io/gousset/openapi/example"
	"github.com/pasqal-io/gousset/openapi/shared"
)

// A JSON schema.
type Schema interface {
	sealed()
}

// Use this interface to customize the schema returned by
// your type.
type HasSchema interface {
	Schema() Schema
}

// Data shared between schemas.
type Shared struct {
	// Optional external documentation.
	ExternalDocs *doc.External `json:"externalDocs,omitempty"`

	// An optional example.
	Example *shared.Json `json:"example,omitempty"`

	// The JavaScript type for this schema.
	Type Type `json:"type"`

	// A well-known format, e.g. "email".
	Format *string `json:"format,omitempty"`

	Title            *string  `json:"title,omitempty"`
	MultipleOf       *float64 `json:"multipleOf,omitempty"`
	Maximum          *float64 `json:"maximum,omitempty"`
	ExclusiveMaximum *float64 `json:"exclusiveMaximum,omitempty"`
	Minimum          *float64 `json:"minimum,omitempty"`
	ExclusiveMinimum *float64 `json:"exclusiveMinimum,omitempty"`
	MaxLength        *int64   `json:"maxLength,omitempty"`
	MinLength        *int64   `json:"minLength,omitempty"`
	Pattern          *string  `json:"pattern,omitempty"`
	MaxItems         *int64   `json:"maxItems,omitempty"`
	MinItems         *int64   `json:"minItems,omitempty"`
	MaxProperties    *int64   `json:"maxProperties,omitempty"`
	MinProperties    *int64   `json:"minProperties,omitempty"`
	Enum             *[]any   `json:"enum,omitempty"`
}

type Type string

const (
	TypeString = Type("string")
	TypeObject = Type("object")
	TypeArray  = Type("array")
	TypeNumber = Type("number")
	TypeBool   = Type("boolean")
)

type Primitive struct {
	Shared `flatten:""`
}

func (Primitive) sealed() {}

type Object struct {
	Shared `flatten:""`

	// A list of required fields
	Required []string `json:"required,omitempty"`

	// A list of permitted fields, with their associated type. Use
	// this if your object is used as an object, with a well-known
	// list of properties.
	Properties map[string]Schema `json:"properties,omitempty"`

	// A list of regexps for additional fields, with their associated
	// type. Use this if your object is used as a map, with an unknown
	// list of properties.
	PatternProperties map[string]Schema `json:"patternProperties,omitempty"`
}

func (Object) sealed() {}

// Array combinator.
type Array struct {
	Shared `flatten:""`
	// The type of items.
	Items Schema `json:"items"`
	// Optionally, definitions used within `Items`.
	Defs *map[string]Schema `json:"$defs,omitempty"`
}

func (Array) sealed() {}

// "One of" combinator, e.g. a sum type.
type OneOf struct {
	OneOf []Schema `json:"oneOf"`
}

func (OneOf) sealed() {}

// "All of" combinator, e.g. an intersection type.
type AllOf struct {
	AllOf []Schema `json:"allOf"`
}

func (AllOf) sealed() {}

type Implementation struct {
	Type             reflect.Type
	PublicNameKey    string
	Title            *string
	MultipleOf       *float64
	Maximum          *float64
	ExclusiveMaximum *float64
	Minimum          *float64
	ExclusiveMinimum *float64
	MaxLength        *int64
	MinLength        *int64
	Pattern          *string
	MaxItems         *int64
	MinItems         *int64
	MaxProperties    *int64
	MinProperties    *int64
	Enum             *[]any
	Format           *string
}

// Create a schema from a type.
//
// As of this writing, we make no attempt to optimize schemas if e.g. some data structures are repeated.
//
// Arguments:
//
//   - typ The type to extract. See HasSchema, HasExternalDocs, HasExample for means to configure it.
//   - publicNameKey The tag used to represent the public name of this field, e.g. `json`, `query`, `path`.
func FromImplementation(impl Implementation) (Schema, error) {
	errorReturn := AllOf{}
	share := Shared{
		Title:            impl.Title,
		MultipleOf:       impl.MultipleOf,
		Maximum:          impl.Maximum,
		ExclusiveMaximum: impl.ExclusiveMaximum,
		Minimum:          impl.Minimum,
		ExclusiveMinimum: impl.ExclusiveMinimum,
		MaxLength:        impl.MaxLength,
		MinLength:        impl.MinLength,
		Pattern:          impl.Pattern,
		MaxItems:         impl.MaxItems,
		MinItems:         impl.MinItems,
		MaxProperties:    impl.MaxProperties,
		MinProperties:    impl.MinProperties,
		Enum:             impl.Enum,
		Format:           impl.Format,
	}
	phony := reflect.New(impl.Type)
	if phony.CanInterface() {
		asAny := phony.Interface()
		// Give priority to a schema provided by the user.
		if hasSchema, ok := asAny.(HasSchema); ok {
			return hasSchema.Schema(), nil
		}
		if hasExternalDocs, ok := asAny.(doc.HasExternalDocs); ok {
			share.ExternalDocs = shared.Ptr(hasExternalDocs.Docs())
		}
		if hasExample, ok := asAny.(example.HasExample); ok {
			share.Example = shared.Ptr(hasExample.Example())
		}
		if hasFormat, ok := asAny.(HasFormat); ok {
			share.Format = shared.Ptr(string(hasFormat.Format()))
		}
		if hasEnum, ok := asAny.(HasEnum); ok {
			share.Enum = shared.Ptr(hasEnum.Enum())
		}
	}

	switch impl.Type.Kind() {
	case reflect.Pointer:
		subImpl := impl
		subImpl.Type = impl.Type.Elem()
		return FromImplementation(subImpl)
	case reflect.Bool:
		share.Type = TypeBool
	case reflect.Float32:
		share.Type = TypeNumber
		if share.Format == nil {
			share.Format = shared.Ptr("float")
		}
	case reflect.Float64:
		share.Type = TypeNumber
		if share.Format == nil {
			share.Format = shared.Ptr("double")
		}
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		share.Type = TypeNumber
		if share.Format == nil {
			share.Format = shared.Ptr("int32")
		}
	case reflect.Int64:
		fallthrough
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		share.Type = TypeNumber
		if share.Format == nil {
			share.Format = shared.Ptr("int64")
		}
	case reflect.String:
		share.Type = TypeString
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		subImpl := impl
		subImpl.Type = impl.Type.Elem()
		items, err := FromImplementation(subImpl)
		if err != nil {
			return errorReturn, fmt.Errorf("while compiling schema for %s, failed to extract type from the elements of array/slice", impl.Type.String())
		}
		share.Type = TypeArray
		return Array{
			Shared: share,
			Items:  items,
			Defs:   nil,
		}, nil
	case reflect.Struct:
		var required []string
		properties := make(map[string]Schema)
		patternProperties := make(map[string]Schema)
		var fields []reflect.StructField
		for i := 0; i < impl.Type.NumField(); i++ {
			fields = append(fields, impl.Type.Field(i))
		}
		for _, field := range fields {
			tags, err := tags.Parse(field.Tag)
			if err != nil {
				return errorReturn, fmt.Errorf("while compiling schema for struct %s, failed to parse tags for field %s", impl.Type.String(), field.Name)
			}

			if tags.IsFlattened() {
				switch field.Type.Kind() {
				case reflect.Struct:
					for i := 0; i < field.Type.NumField(); i++ {
						fields = append(fields, field.Type.Field(i))
					}
					continue
				case reflect.Map:
					subImpl := impl
					subImpl.Type = field.Type.Elem()
					switch field.Type.Key().Kind() {
					case reflect.String:
						scheme, err := FromImplementation(subImpl)
						if err != nil {
							return errorReturn, fmt.Errorf("while compiling schema for struct %s, cannot extract scheme for contents of map at field %s: %w", impl.Type.String(), field.Name, err)
						}
						patternProperties["*"] = scheme
						continue
					default:
						return errorReturn, fmt.Errorf("while compiling schema for %s, this type of map is not implemented at field %s", impl.Type.String(), field.Name)
					}
				default:
					return errorReturn, fmt.Errorf("while compiling schema for %s, field %s is marked as flattened but is not a struct", impl.Type.String(), field.Name)
				}
			}

			var name string
			if publicName := tags.PublicFieldName(impl.PublicNameKey); publicName != nil {
				name = *publicName
			} else {
				return errorReturn, fmt.Errorf("while compiling schema for struct %s, field %s doesn't have a public name, expecting a tag `%s`", impl.Type.String(), field.Name, impl.PublicNameKey)
			}

			subImpl, err := ImplementationFromStructField(field, impl.PublicNameKey)
			if err != nil {
				return errorReturn, fmt.Errorf("while compiling a schema for struct %s, error in field %s: %w", impl.Type.String(), field.Name, err)
			}

			fieldSchema, err := FromImplementation(subImpl)
			if err != nil {
				return errorReturn, fmt.Errorf("while compiling schema for %s, failed to extract scheme from field %s: %w", impl.Type.String(), field.Name, err)
			}

			if tags.Default() == nil && !tags.IsPreinitialized() && tags.MethodName() == nil {
				required = append(required, name)
			}
			properties[name] = fieldSchema
		}
		share.Type = TypeObject
		return Object{
			Shared:     share,
			Required:   required,
			Properties: properties,
		}, nil
	case reflect.Map:
		if impl.Type.Key().Kind() != reflect.String {
			return errorReturn, fmt.Errorf("while compiling schema for map %s, this key type isn't supported %s", impl.Type.String(), impl.Type.Key().String())
		}
		subImpl := Implementation{
			Type:          impl.Type.Elem(),
			PublicNameKey: impl.PublicNameKey,
		}
		contentSchema, err := FromImplementation(subImpl)
		if err != nil {
			return errorReturn, fmt.Errorf("while compiling schema for map %s, failed to extract scheme from content: %w", impl.Type.String(), err)
		}
		patternedProperties := make(map[string]Schema)
		patternedProperties["*"] = contentSchema
		share.Type = TypeObject
		return Object{
			Shared:            share,
			Required:          []string{},
			Properties:        make(map[string]Schema),
			PatternProperties: patternedProperties,
		}, nil
	default:
		return errorReturn, fmt.Errorf("while compiling schema for %s, couldn't find any scheme, you may need to implement HasSchema", impl.Type.String())
	}
	// If we have reached this point, we're dealing with a primitive.
	return Primitive{
		Shared: share,
	}, nil
}

func ImplementationFromStructField(field reflect.StructField, publicNameKey string) (Implementation, error) {
	tags, err := tags.Parse(field.Tag)
	if err != nil {
		return Implementation{}, fmt.Errorf("failed to parse tags for field %s: %w", field.Name, err)
	}

	result := Implementation{
		Type:          field.Type,
		PublicNameKey: publicNameKey,
	}
	type Pair[T any, U any] struct {
		First  T
		Second U
	}
	for _, parse := range []Pair[**string, string]{
		{&result.Title, "title"},
		{&result.Format, "format"},
		{&result.Pattern, "pattern"},
	} {
		*parse.First = tags.LookupString(parse.Second)
	}
	for _, parse := range []Pair[**float64, string]{
		{&result.ExclusiveMaximum, "exclusiveMaximum"},
		{&result.ExclusiveMinimum, "exclusiveMinimum"},
		{&result.Minimum, "minimum"},
		{&result.Maximum, "maximum"},
		{&result.MultipleOf, "multipleOf"},
	} {
		parsed, err := tags.LookupFloat(parse.Second)
		if err != nil {
			return Implementation{}, fmt.Errorf("while compiling schema for field %s, error: %w",
				field.Name, err)
		}
		*parse.First = parsed
	}
	for _, parse := range []Pair[**int64, string]{
		{&result.MaxItems, "maxItems"},
		{&result.MinItems, "minItems"},
		{&result.MaxLength, "maxLength"},
		{&result.MinLength, "minLength"},
		{&result.MaxProperties, "maxProperties"},
		{&result.MinProperties, "minProperties"},
	} {
		parsed, err := tags.LookupInt(parse.Second)
		if err != nil {
			return Implementation{}, fmt.Errorf("while compiling schema for field %s, error: %w",
				field.Name, err)
		}
		*parse.First = parsed
	}
	return result, nil
}

// Implement this on a type to specify that it should be marked as an enum.
type HasEnum interface {
	Enum() []shared.Json
}

// Implement this on a type to specify that it should have a given format.
type HasFormat interface {
	Format() Format
}

// A well-known format.
type Format string

const (
	FormatDateTime    = Format("date-time")
	FormatDate        = Format("date")
	FormatTime        = Format("time")
	FormatDuration    = Format("duration")
	FormatEmail       = Format("email")
	FormatIdnEmail    = Format("idn-email")
	FormatHostname    = Format("hostname")
	FormatIdnHostname = Format("idn-hostname")
	FormatUri         = Format("uri")
	FormatRegex       = Format("regex")
	FormatBinary      = Format("binary")
	FormatInt32       = Format("int32")
	FormatInt64       = Format("int64")
	FormatFloat       = Format("float")
	FormatDouble      = Format("double")
	FormatByte        = Format("byte")
	FormatPassword    = Format("password")
)
