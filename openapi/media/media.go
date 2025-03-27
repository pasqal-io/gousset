package media

import (
	"fmt"
	"reflect"

	"github.com/pasqal-io/gousset/openapi/example"
	"github.com/pasqal-io/gousset/openapi/schema"
	"github.com/pasqal-io/gousset/openapi/shared"
)

type Type struct {
	Schema   *schema.Schema              `json:"schema"`
	Example  *shared.Json                `json:"example"`
	Examples *map[string]example.Example `json:"examples"`
}

func FromBody(body reflect.Type, publicNameKey string) (Type, error) {
	schema, err := schema.FromType(body, publicNameKey)
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
