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

// MyResult represents a sum between types Ok and Err.
type MyResult interface {
	// Prevent other types from implementing MySum.
	sealed()
	// Make sure that we can generate proper OpenAPI spec
	// for this type.
	schema.IsOneOf
}

type Ok struct {
	Result any `json:"result"`
}

type Err struct {
	Error error `json:"error"`
}

// Implementing sealed, to make sure that only Ok and Err are part of the sum type.
func (Ok) sealed() {
}
func (Err) sealed() {
}

func (Ok) Type() reflect.Type {
	return reflect.TypeFor[MyResult]()
}
func (Ok) Variants() []reflect.Type {
	return []reflect.Type{reflect.TypeOf(Ok{}), reflect.TypeOf(Err{})}
}
func (Err) Type() reflect.Type {
	return Ok{}.Type()
}
func (Err) Variants() []reflect.Type {
	return Ok{}.Variants()
}

var _ MyResult = Ok{}

var _ MyResult = Err{}

// Do NOT forget to call schema.RegisterOneOf!
var _ = schema.RegisterOneOf(&Ok{})

func ExampleIsOneOf() {
	spec, err := schema.FromImplementation(schema.Implementation{
		Type:          reflect.TypeFor[MyResult](),
		PublicNameKey: "json",
	})
	if err != nil {
		panic(err)
	}
	marshaled, err := json.MarshalIndent(spec, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Print(string(marshaled))
	// Output: {
	// 	"oneOf": [
	// 		{
	// 			"type": "object",
	// 			"required": [
	// 				"result"
	// 			],
	// 			"properties": {
	// 				"result": {
	// 					"type": ""
	// 				}
	// 			}
	// 		},
	// 		{
	// 			"type": "object",
	// 			"required": [
	// 				"error"
	// 			],
	// 			"properties": {
	// 				"error": {
	// 					"type": ""
	// 				}
	// 			}
	// 		}
	// 	]
	// }
}

func TestIsOneOf(t *testing.T) {
	spec, err := schema.FromImplementation(schema.Implementation{
		Type:          reflect.TypeFor[MyResult](),
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
				"result"
			],
			"properties": {
				"result": {
				"type": ""
				}
			}
			},
			{
			"type": "object",
			"required": [
				"error"
			],
			"properties": {
				"error": {
				"type": ""
				}
			}
			}
		]
		}`)
}
