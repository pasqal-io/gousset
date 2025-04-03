// Headers, both within requests and within responses.
package header

import (
	"fmt"
	"reflect"

	"github.com/pasqal-io/gousset/openapi/example"
	"github.com/pasqal-io/gousset/openapi/media"
	"github.com/pasqal-io/gousset/openapi/schema"
	"github.com/pasqal-io/gousset/shared"
)

type Header interface {
	sealed()
}

// A reference.
type Reference shared.Reference

func Ref(to string) Reference {
	return Reference(shared.Ref(to))
}
func (Reference) sealed() {}

var _ Header = Ref("")

// https://spec.openapis.org/oas/v3.0.1.html#header-object
type Spec struct {
	Description  *string `json:"description"`
	Required     bool    `json:"required"`
	Deprecated   bool    `json:"deprecated"`
	*SchemaSpec  `json:"-1,omitempty" flatten:""`
	*ContentSpec `json:"-2,omitempty" flatten:""`
}

func (Spec) sealed() {}

var _ Header = Spec{}

// https://spec.openapis.org/oas/v3.0.1.html#schema-object
type SchemaSpec struct {
	// Describes how the parameter value will be serialized depending on the type of the parameter value. Default values (based on value of in): for "query" - "form"; for "path" - "simple"; for "header" - "simple"; for "cookie" - "form".
	Style string `json:"style"`

	// When this is true, parameter values of type array or object generate separate parameters for each value of the array or key-value pair of the map. For other types of parameters this field has no effect. When style is "form", the default value is true. For all other styles, the default value is false. Note that despite false being the default for deepObject, the combination of false with deepObject is undefined.
	Explode bool `json:"explode"`

	Schema schema.Schema `json:"schema"`

	Example  *shared.Json       `json:"example,omitempty"`
	Examples *[]example.Example `json:"examples,omitempty"`
}

type ContentSpec struct {
	Content map[string]media.Type `json:"content"`
}

// User-provided metadata containing information on the implementation
// to be converted to OpenAPI spec.
type Implementation struct {
	Description *string
	Required    bool
	Deprecated  bool
	SchemaSpec  *SchemaImplementation
	ContentSpec *ContentImplementation
}

type SchemaImplementation struct {
	Type     reflect.Type
	Example  *shared.Json
	Examples *[]example.Example
}

type ContentImplementation struct {
	Content map[string]media.Implementation
}

func FromImplementation(impl Implementation) (Header, error) {
	result := Spec{
		Description: impl.Description,
		Required:    impl.Required,
		Deprecated:  impl.Deprecated,
	}
	if impl.SchemaSpec != nil {
		schema, err := FromSchemaImplementation(*impl.SchemaSpec)
		if err != nil {
			return result, fmt.Errorf("while compiling header, error in schema: %w", err)
		}
		result.SchemaSpec = &schema
	}
	if impl.ContentSpec != nil {
		content, err := FromContentImplementation(*impl.ContentSpec)
		if err != nil {
			return result, fmt.Errorf("while compiling header, error in content: %w", err)
		}
		result.ContentSpec = &content
	}
	return result, nil
}

func FromSchemaImplementation(impl SchemaImplementation) (SchemaSpec, error) {
	result := SchemaSpec{
		Style:    "simple",
		Example:  impl.Example,
		Examples: impl.Examples,
	}
	schema, err := schema.FromImplementation(schema.Implementation{Type: impl.Type, PublicNameKey: "header"})
	if err != nil {
		return result, fmt.Errorf("while compiling schema, error: %w", err)
	}
	result.Schema = schema
	return result, nil
}

func FromContentImplementation(impl ContentImplementation) (ContentSpec, error) {
	result := ContentSpec{
		Content: map[string]media.Type{},
	}
	for k, v := range impl.Content {
		content, err := media.FromImplementation(v)
		if err != nil {
			return result, fmt.Errorf("while compiling content type %s, error: %w", k, err)
		}
		result.Content[k] = content
	}
	return result, nil
}
