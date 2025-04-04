// Definitions of JSON-level data structures.
//
// Any data structure that the API expects or provides, regardless
// of whether it's through body, path, query, headers, is defined
// as a `schema.Schema`.
package schema

import (
	"fmt"
	"log/slog"
	"maps"
	"reflect"
	"time"

	"github.com/pasqal-io/gousset/inner/tags"
	"github.com/pasqal-io/gousset/openapi/doc"
	"github.com/pasqal-io/gousset/openapi/example"
	"github.com/pasqal-io/gousset/shared"
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

	// The types of values used in this object. Use this if your object
	// is used as a map, with an unknown list of properties, all of them
	// with the same type.
	AdditionalProperties *Schema `json:"additionalProperties,omitempty"`
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
	Example          *string
}

var stringType = reflect.TypeOf("")

func fill[T any, I any](field **T, value any, cb func(I) T) {
	if *field != nil {
		return
	}
	asInterface, ok := value.(I)
	if !ok {
		return
	}
	*field = shared.Ptr(cb(asInterface))
}

// Extract the scheme for a single variant.
//
// If `restriction` is not nil, only consider the fields that are keys of `restriction`.
func fromImplementationSingleVariant(impl Implementation, restriction map[string]bool) (Schema, error) {
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
		fill(&share.ExternalDocs, asAny, func(value doc.HasExternalDocs) doc.External { return value.Docs() })
		fill(&share.Example, asAny, func(value example.HasExample) shared.Json { return value.Example() })
		fill(&share.Format, asAny, func(value HasFormat) string { return string(value.Format()) })
		fill(&share.Enum, asAny, func(value IsEnum) []shared.Json { return value.Enum() })
		fill(&share.MinItems, asAny, func(value HasMinArrayLength) int64 { return value.MinArrayLength() })
		fill(&share.MaxItems, asAny, func(value HasMaxArrayLength) int64 { return value.MaxArrayLength() })
		fill(&share.MinLength, asAny, func(value HasMinStringLength) int64 { return value.MinStringLength() })
		fill(&share.MaxLength, asAny, func(value HasMaxStringLength) int64 { return value.MaxStringLength() })
		fill(&share.MinProperties, asAny, func(value HasMinMapLength) int64 { return value.MinMapLength() })
		fill(&share.MaxProperties, asAny, func(value HasMaxMapLength) int64 { return value.MaxMapLength() })
		fill(&share.Minimum, asAny, func(value HasMin) float64 { return value.Min() })
		fill(&share.Maximum, asAny, func(value HasMax) float64 { return value.Max() })

		// If this is a time, let's not look further.
		if isTime(asAny) {
			share.Type = TypeString
			if share.Format == nil {
				share.Format = shared.Ptr(string(FormatDateTime))
			}
			return Primitive{
				Shared: share,
			}, nil
		}
	}
	if impl.Example != nil {
		var example any = *impl.Example
		share.Example = &example
	}

	switch impl.Type.Kind() {
	case reflect.Interface:
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
			field := impl.Type.Field(i)
			if restriction != nil {
				if _, ok := restriction[field.Name]; !ok {
					continue
				}
			}
			fields = append(fields, field)
		}
		for len(fields) != 0 {
			currentFields := fields
			fields = []reflect.StructField{}
			for _, field := range currentFields {
				if !field.IsExported() {
					continue
				}
				tags, err := tags.Parse(field.Tag)
				if err != nil {
					return errorReturn, fmt.Errorf("while compiling schema for struct %s, failed to parse tags for field %s", impl.Type.String(), field.Name)
				}

				if tags.IsFlattened() {
					// Copy fields/entries from this type into the parent.
					typ := field.Type
					for typ.Kind() == reflect.Pointer {
						typ = typ.Elem()
					}
					switch typ.Kind() {
					case reflect.Struct:
						// Let's do it again, but with children fields instead of this field.
						for i := 0; i < typ.NumField(); i++ {
							fields = append(fields, typ.Field(i))
						}
						continue
					case reflect.Map:
						subImpl := impl
						subImpl.Type = typ.Elem()
						switch typ.Key().Kind() {
						case reflect.String:
							scheme, err := FromImplementation(subImpl)
							if err != nil {
								return errorReturn, fmt.Errorf("while compiling schema for struct %s, cannot extract scheme for contents of map at field %s: %w", impl.Type.String(), field.Name, err)
							}
							patternProperties["*"] = scheme
						default:
							return errorReturn, fmt.Errorf("while compiling schema for %s, this type of map is not implemented at field %s", impl.Type.String(), field.Name)
						}
					default:
						return errorReturn, fmt.Errorf("while compiling schema for %s, field %s is marked as flattened but is not a struct, got %s", impl.Type.String(), field.Name, field.Type.String())
					}
				} else {
					// Treat field as an object.
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
			}
		}
		share.Type = TypeObject
		return Object{
			Shared:     share,
			Required:   required,
			Properties: properties,
		}, nil
	case reflect.Map:
		if !isStringifiable(impl.Type.Key()) {
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
		share.Type = TypeObject
		return Object{
			Shared:               share,
			Required:             []string{},
			Properties:           make(map[string]Schema),
			AdditionalProperties: &contentSchema,
		}, nil
	default:
		return errorReturn, fmt.Errorf("while compiling schema for %s, couldn't find any scheme, you may need to implement IsOneOf or HasSchema or to call RegisterOneOf", impl.Type.String())
	}
	// If we have reached this point, we're dealing with a primitive.
	return Primitive{
		Shared: share,
	}, nil
}

// Create a schema from a type.
//
// As of this writing, we make no attempt to optimize schemas if e.g. some data structures are repeated.
func FromImplementation(impl Implementation) (Schema, error) {
	if impl.Type.Kind() != reflect.Struct {
		return fromImplementationSingleVariant(impl, nil)
	}
	// Make a first pass to determine whether this is a sum type.
	shared := []string{}
	variants := map[string][]string{}
	for i := 0; i < impl.Type.NumField(); i++ {
		field := impl.Type.Field(i)
		tags, err := tags.Parse(field.Tag)
		if err != nil {
			return OneOf{}, fmt.Errorf("while compiling schema for struct type %s, during the first pass, error while parsing tags for field %s: %w",
				impl.Type.String(),
				field.Name,
				err)
		}
		if fieldVariants, ok := tags.Variants(); ok {
			// This field belongs to at least one variant (which means that there are variants).
			for _, variant := range fieldVariants {
				var addMeTo []string
				if v, found := variants[variant]; found {
					addMeTo = v
				} else {
					addMeTo = []string{}
				}
				variants[variant] = append(addMeTo, field.Name)
			}
		} else {
			// This field belongs to all variants (which might mean that there are no variants).
			shared = append(shared, field.Name)
		}
	}
	switch len(variants) {
	case 0:
		// No variants, just a regular struct.
		return fromImplementationSingleVariant(impl, nil)
	case 1:
		var key string
		for k, _ := range variants {
			key = k
			break
		}
		slog.Warn("gousset.openapi.schema.FromImplementation: encountering a struct with a single variant, tag `variants` seems misapplied",
			"type", impl.Type.String(),
			"variant", key)
		return fromImplementationSingleVariant(impl, nil)
	default:
		// Alright, this is a true sum type, let's build it.
		result := OneOf{}

		// Order the variant names by a guaranteed stable order.
		// This is mostly to make testing easier.
		variantNames := []string{}
		for variantName := range maps.Keys(variants) {
			variantNames = append(variantNames, variantName)
		}

		for _, variantName := range variantNames {
			restriction := map[string]bool{}
			for _, name := range variants[variantName] {
				restriction[name] = true
			}
			for _, name := range shared {
				restriction[name] = true
			}
			variant, err := fromImplementationSingleVariant(impl, restriction)
			if err != nil {
				return AllOf{}, fmt.Errorf("while compiling schema for sum type %s, error in variant %s: %w", impl.Type.String(), variantName, err)
			}
			result.OneOf = append(result.OneOf, variant)
		}
		return result, nil
	}
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
		{&result.Example, "example"},
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
			return Implementation{}, fmt.Errorf("while compiling schema for field %s, error within %s key %s: %w",
				"float", parse.Second,
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
			return Implementation{}, fmt.Errorf("while compiling schema for field %s, error within %s key %s: %w",
				"int", parse.Second,
				field.Name, err)
		}
		*parse.First = parsed
	}
	return result, nil
}

// Implement this on a type to specify that it should be marked as an enum.
//
// For more sophisticated cases, see `IsOneOf`.
type IsEnum interface {
	// The list of possibilities for this enum.
	Enum() []shared.Json
}

// Implement this on a string or number to specify that it should match a given format.
type HasFormat interface {
	// The format restriction.
	Format() Format
}

// A well-known format.
type Format string

// Well-known formats.
//
// Note that this list doesn't attempt to be exhaustive.
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

func isTime(value any) bool {
	if _, ok := value.(time.Time); ok {
		return true
	}
	if _, ok := value.(*time.Time); ok {
		return true
	}
	return false
}

type hasString interface {
	String() string
}

func isStringifiable(typ reflect.Type) bool {
	if stringType.ConvertibleTo(typ) {
		return true
	}
	phony := reflect.New(typ)
	if !phony.CanInterface() {
		return false
	}
	if _, ok := phony.Interface().(hasString); ok {
		return true
	}
	return false
}

// Implement this to mark a minimal length for a string.
type HasMinStringLength interface {
	MinStringLength() int64
}

// Implement this to mark a maximal length for a string.
type HasMaxStringLength interface {
	MaxStringLength() int64
}

// Implement this to mark a minimal length for an array.
type HasMinArrayLength interface {
	MinArrayLength() int64
}

// Implement this to mark a maximal length for an array.
type HasMaxArrayLength interface {
	MaxArrayLength() int64
}

// Implement this to mark a minimal length for an array.
type HasMinMapLength interface {
	MinMapLength() int64
}

// Implement this to mark a maximal length for an array.
type HasMaxMapLength interface {
	MaxMapLength() int64
}

// Implement this to mark a minimal value for a number.
type HasMin interface {
	Min() float64
}

// Implement this to mark a maximal value for a number.
type HasMax interface {
	Max() float64
}
