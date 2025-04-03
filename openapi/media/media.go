package media

import (
	"fmt"
	"reflect"

	"github.com/pasqal-io/gousset/openapi/example"
	"github.com/pasqal-io/gousset/openapi/schema"
	"github.com/pasqal-io/gousset/shared"
)

type Type struct {
	Schema   *schema.Schema              `json:"schema,omitempty"`
	Example  *shared.Json                `json:"example,omitempty"`
	Examples *map[string]example.Example `json:"examples,omitempty"`
}

func FromBody(body reflect.Type, publicNameKey string) (Type, error) {
	schema, err := schema.FromImplementation(schema.Implementation{Type: body, PublicNameKey: publicNameKey})
	if err != nil {
		return Type{}, fmt.Errorf("failed to extract schema for body %s: %w", body.String(), err)
	}
	anExample := example.GetExample(body)
	examples := example.GetExamples(body)
	return Type{
		Schema:   &schema,
		Example:  anExample,
		Examples: examples,
	}, nil
}

// User-provided metadata containing information on the implementation
// to be converted to OpenAPI spec.
type Implementation struct {
	Type          reflect.Type
	Example       *shared.Json
	Examples      *map[string]example.Example
	PublicNameKey string
}

func FromImplementation(impl Implementation) (Type, error) {
	result := Type{
		Example:  impl.Example,
		Examples: impl.Examples,
	}
	// Default to json.
	if impl.PublicNameKey == "" {
		impl.PublicNameKey = "json"
	}
	if impl.Type.Kind() != reflect.Invalid {
		typ, err := schema.FromImplementation(schema.Implementation{Type: impl.Type, PublicNameKey: impl.PublicNameKey})
		if err != nil {
			return result, fmt.Errorf("while collecting media type, error: %w", err)
		}
		result.Schema = &typ
	}
	return result, nil
}
