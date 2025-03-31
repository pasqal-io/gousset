package parameter

import (
	"fmt"
	"reflect"

	tags "github.com/pasqal-io/gousset/inner"
	"github.com/pasqal-io/gousset/openapi/doc"
	"github.com/pasqal-io/gousset/openapi/example"
	"github.com/pasqal-io/gousset/openapi/media"
	"github.com/pasqal-io/gousset/openapi/schema"
	"github.com/pasqal-io/gousset/openapi/shared"
)

type Parameter interface {
	parameter()
}

// Specification for a parameter.
type Spec struct {
	// The name of the parameter. Parameter names are case sensitive.
	//
	// If in is "path", the name field MUST correspond to a template expression occurring within the path field in the Paths Object. See Path Templating for further information.
	// If in is "header" and the name field is "Accept", "Content-Type" or "Authorization", the parameter definition SHALL be ignored.
	// For all other cases, the name corresponds to the parameter name used by the in field.
	Name string `json:"name"`

	// The location of the parameter. Possible values are "query", "header", "path" or "cookie".
	In In `json:"in"`

	// Short summary.
	Summary string `json:"summary"`

	// Longer description. May include Markdown.
	Description *string `json:"description,omitempty"`

	Required bool `json:"required"`

	Deprecated bool `json:"deprecated"`

	// Structure and syntax of the parameter.
	//
	// Mutually exclusive with Schema.
	*ContentSpec `json:"content,omitempty"`

	// Media type and schema for the parameter.
	//
	// Mutually exclusive with Content.
	*SchemaSpec `json:"schema,omitempty"`
}

func FromField(from reflect.StructField, in In) (Spec, error) {
	publicNameKey := string(in)
	tags, err := tags.Parse(from.Tag)
	if err != nil {
		return Spec{}, fmt.Errorf("failed to parse tags for field %s, got %w", from.Name, err)
	}

	// We make the decision of only documenting fields that have a public field name.
	tagFieldName := tags.PublicFieldName(publicNameKey)
	if tagFieldName == nil {
		return Spec{}, fmt.Errorf("field %s is missing a public name", from.Name)
	}

	// Extract summary and description.
	summary := doc.GetSummary(from.Type)
	description := doc.GetDescription(from.Type)

	if summary == nil {
		if tagSummary, ok := tags.Lookup("summary"); ok && len(tagSummary) >= 1 {
			summary = shared.Ptr(tagSummary[0])
		} else {
			return Spec{}, fmt.Errorf("field %s doesn't have a summary, please provide a method Summary() or a tag `summary`", from.Name)
		}
	}

	deprecated := false
	if _, ok := tags.Lookup("deprecated"); ok {
		deprecated = true
	}

	required := true
	if (tags.Default() != nil) || tags.IsPreinitialized() || (tags.MethodName() != nil) {
		required = false
	}

	schema, err := schema.FromType(from.Type, publicNameKey)
	if err != nil {
		return Spec{}, fmt.Errorf("failed to find schema for field %s: %w", from.Name, err)
	}
	schemaSpec := SchemaSpec{
		Style:   nil,
		Explode: nil,
		Schema:  schema,
	}

	return Spec{
		Name:        *tagFieldName, // Non-nil, checked above.
		In:          in,
		Summary:     *summary, // Non-nil, checked above.
		Description: description,
		Deprecated:  deprecated,
		Required:    required,
		SchemaSpec:  &schemaSpec,
	}, nil
}

func FromStruct(Struct reflect.Type, in In) ([]Parameter, error) {
	if Struct.Kind() != reflect.Struct {
		return []Parameter{}, fmt.Errorf("invalid type %s, expected a struct, got %v.", Struct.String(), Struct.Kind())
	}

	var parameters []Parameter
	for i := 0; i < Struct.NumField(); i++ { // We have checked above that it's a struct.
		field := Struct.Field(i)
		// FIXME: We'll need to know if there are any default values.
		param, err := FromField(field, in)
		if err != nil {
			return []Parameter{}, fmt.Errorf("failed to generate spec for parameter %s of type %s: %w", field.Name, Struct.String(), err)
		}
		parameters = append(parameters, param)
	}

	return parameters, nil
}

type In string

const (
	InQuery  = In("query")
	InHeader = In("header")
	InCookie = In("cookie")
	InPath   = In("path")
)

func (Spec) parameter() {

}

var _ Parameter = Spec{}

// Reference to a component.
type Reference struct {
}

func (Reference) parameter() {

}

var _ Parameter = Reference{}

type SchemaSpec struct {
	// Describes how the parameter value will be serialized depending on the type of the parameter value. Default values (based on value of in): for "query" - "form"; for "path" - "simple"; for "header" - "simple"; for "cookie" - "form".
	Style *string `json:"style,omitempty"`

	// When this is true, parameter values of type array or object generate separate parameters for each value of the array or key-value pair of the map. For other types of parameters this field has no effect. When style is "form", the default value is true. For all other styles, the default value is false. Note that despite false being the default for deepObject, the combination of false with deepObject is undefined.
	Explode *bool `json:"explode,omitempty"`

	Schema schema.Schema `json:"schema"`

	Example  *shared.Json       `json:"example,omitempty"`
	Examples *[]example.Example `json:"examples,omitempty"`
}

type ContentSpec struct {
	Content map[string]media.Type `json:"content"`
}
