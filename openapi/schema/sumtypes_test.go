package schema_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/pasqal-io/gousset/openapi/schema"
	"github.com/pasqal-io/gousset/testutils"
)

// -----   Example: Using sum types.

type Unflattened[T any, U any] struct {
	Comments string `json:"comments"`
	Ok       *T     `json:"ok" variant:"ok"`
	Error    *U     `json:"error" variant:"error"`
}

func ExampleSchema() {
	spec, err := schema.FromImplementation(schema.Implementation{
		Type:          reflect.TypeFor[Unflattened[int, bool]](),
		PublicNameKey: "json",
	})
	if err != nil {
		panic(err)
	}
	marshaled, err := json.Marshal(spec)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(marshaled))
	// Output: {"oneOf":[{"type":"object","required":["comments","ok"],"properties":{"comments":{"type":"string"},"ok":{"type":"number","format":"int32"}}},{"type":"object","required":["comments","error"],"properties":{"comments":{"type":"string"},"error":{"type":"boolean"}}}]}
}

type Flattened[T any, U any] struct {
	Comments string `json:"comments"`
	Ok       *T     `flatten:"" variant:"ok"`
	Error    *U     `flatten:"" variant:"error"`
}

func TestFlattened(t *testing.T) {
	type Success struct {
		Success string `json:"success"`
	}
	type Failure struct {
		Failure bool `json:"failure"`
	}
	spec, err := schema.FromImplementation(schema.Implementation{
		Type:          reflect.TypeFor[Flattened[Success, Failure]](),
		PublicNameKey: "json",
	})
	if err != nil {
		panic(err)
	}
	testutils.EqualJSON(t, spec, `{
		"oneOf": [
			{
				"type": "object",
				"required": [
					"comments",
					"success"
				],
				"properties": {
					"comments": {
						"type": "string"
					},
					"success": {
						"type": "string"
					}
				}
			},
			{
				"type": "object",
				"required": [
					"comments",
					"failure"
				],
				"properties": {
					"comments": {
						"type": "string"
					},
					"failure": {
						"type": "boolean"
					}
				}
			}
		]
		}`)
}

func TestMoreFlattened(t *testing.T) {
	type Success struct {
		Success string `json:"result"`
	}
	type Failure struct {
		Failure bool `json:"result"`
	}
	spec, err := schema.FromImplementation(schema.Implementation{
		Type:          reflect.TypeFor[Flattened[Success, Failure]](),
		PublicNameKey: "json",
	})
	if err != nil {
		panic(err)
	}
	testutils.EqualJSON(t, spec, `{
		"oneOf": [
			{
			"type": "object",
			"required": [
				"comments",
				"result"
			],
			"properties": {
				"comments": {
					"type": "string"
				},
				"result": {
					"type": "string"
				}
			}
			},
			{
			"type": "object",
			"required": [
				"comments",
				"result"
			],
			"properties": {
				"comments": {
					"type": "string"
				},
				"result": {
					"type": "boolean"
				}
			}
			}
		]
		}`)
}
